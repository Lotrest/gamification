package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	gamificationv1 "cdek/platform/shared/contracts/gamification/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

const (
	defaultCompany  = "CDEK Digital"
	defaultLocation = "Новосибирск"
	defaultTeam     = "Платформа комьюнити"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Title    string `json:"title"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createArticleRequest struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Body    string `json:"body"`
}

type reactionRequest struct {
	Reaction string `json:"reaction"`
}

type createCommentRequest struct {
	Body string `json:"body"`
}

func (s *Server) Register(ctx *fiber.Ctx) error {
	var request registerRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Password = strings.TrimSpace(request.Password)
	request.Title = strings.TrimSpace(request.Title)

	if request.Name == "" || request.Email == "" || request.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "name, email and password are required"})
	}

	if !strings.Contains(request.Email, "@") {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "email must be valid"})
	}

	if request.Title == "" {
		request.Title = "Участник платформы"
	}

	requestContext, span := s.tracer.Start(ctx.UserContext(), "register")
	defer span.End()

	tx, err := s.db.BeginTx(requestContext, pgx.TxOptions{})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to start transaction"})
	}
	defer tx.Rollback(requestContext)

	var existingUserID string
	err = tx.QueryRow(
		requestContext,
		`select user_id from user_service.credentials where email = $1`,
		request.Email,
	).Scan(&existingUserID)
	if err == nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "user with this email already exists"})
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.logger.Error("check credential failed", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create user"})
	}

	userID := newID("user")
	joinedAt := time.Now().Format("02.01.2006")
	passwordHash := hashPassword(request.Password)

	if _, err = tx.Exec(
		requestContext,
		`insert into user_service.users
		 (id, name, title, company, level, level_text, joined_at, location, team)
		 values ($1, $2, $3, $4, 1, 'Lv.1', $5, $6, $7)`,
		userID,
		request.Name,
		request.Title,
		defaultCompany,
		joinedAt,
		defaultLocation,
		defaultTeam,
	); err != nil {
		s.logger.Error("insert user failed", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create user"})
	}

	if _, err = tx.Exec(
		requestContext,
		`insert into user_service.credentials (user_id, email, password_hash)
		 values ($1, $2, $3)`,
		userID,
		request.Email,
		passwordHash,
	); err != nil {
		s.logger.Error("insert credential failed", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create user"})
	}

	if err = seedUserState(requestContext, tx, userID); err != nil {
		s.logger.Error("seed user state failed", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to initialize user"})
	}

	if err = tx.Commit(requestContext); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to save user"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token": createTokenForUser(userID),
		"user": fiber.Map{
			"id":      userID,
			"name":    request.Name,
			"title":   request.Title,
			"company": defaultCompany,
			"email":   request.Email,
		},
	})
}

func (s *Server) Login(ctx *fiber.Ctx) error {
	var request loginRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Password = strings.TrimSpace(request.Password)

	if request.Email == "" || request.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "email and password are required"})
	}

	requestContext, span := s.tracer.Start(ctx.UserContext(), "login")
	defer span.End()

	var (
		userID       string
		name         string
		title        string
		company      string
		passwordHash string
	)

	err := s.db.QueryRow(
		requestContext,
		`select u.id, u.name, u.title, u.company, c.password_hash
		 from user_service.credentials c
		 join user_service.users u on u.id = c.user_id
		 where c.email = $1`,
		request.Email,
	).Scan(&userID, &name, &title, &company, &passwordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid email or password"})
		}

		s.logger.Error("login query failed", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to login"})
	}

	if !comparePassword(passwordHash, request.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid email or password"})
	}

	return ctx.JSON(fiber.Map{
		"token": createTokenForUser(userID),
		"user": fiber.Map{
			"id":      userID,
			"name":    name,
			"title":   title,
			"company": company,
			"email":   request.Email,
		},
	})
}

func (s *Server) ListArticles(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	requestContext, span := s.tracer.Start(ctx.UserContext(), "list-articles")
	defer span.End()

	articles, err := s.loadArticles(requestContext, userID)
	if err != nil {
		s.logger.Error("list articles failed", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to load articles"})
	}

	return ctx.JSON(fiber.Map{"articles": articles})
}

func (s *Server) GetArticle(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	articleID := ctx.Params("articleId")
	requestContext, span := s.tracer.Start(ctx.UserContext(), "get-article")
	defer span.End()

	if err := s.incrementArticleViews(requestContext, articleID); err != nil {
		s.logger.Warn("increment article views failed", "articleId", articleID, "error", err)
	}

	article, err := s.loadArticle(requestContext, userID, articleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "article not found"})
		}

		s.logger.Error("get article failed", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to load article"})
	}

	return ctx.JSON(fiber.Map{"article": article})
}

func (s *Server) CreateArticle(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	var request createArticleRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	request.Title = strings.TrimSpace(request.Title)
	request.Summary = strings.TrimSpace(request.Summary)
	request.Body = strings.TrimSpace(request.Body)
	if request.Title == "" || request.Summary == "" || request.Body == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "title, summary and body are required"})
	}

	requestContext, span := s.tracer.Start(ctx.UserContext(), "create-article")
	defer span.End()

	tx, err := s.db.BeginTx(requestContext, pgx.TxOptions{})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to start transaction"})
	}
	defer tx.Rollback(requestContext)

	articleID := newID("article")
	if _, err = tx.Exec(
		requestContext,
		`insert into gamification.articles (id, author_id, title, summary, body)
		 values ($1, $2, $3, $4, $5)`,
		articleID,
		userID,
		request.Title,
		request.Summary,
		request.Body,
	); err != nil {
		s.logger.Error("create article failed", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create article"})
	}

	if _, err = tx.Exec(
		requestContext,
		`insert into gamification.article_cards (id, user_id, title, views, comments, xp, rating)
		 values ($1, $2, $3, 0, 0, 0, '5.0')`,
		articleID,
		userID,
		request.Title,
	); err != nil {
		s.logger.Error("create article card failed", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create article"})
	}

	if _, err = tx.Exec(
		requestContext,
		`update gamification.user_state
		 set articles_count = articles_count + 1,
		     updated_at = now()
		 where user_id = $1`,
		userID,
	); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update article statistics"})
	}

	if _, err = tx.Exec(
		requestContext,
		`insert into gamification.recent_activity (id, user_id, title, timestamp_label, xp, sort_order)
		 values ($1, $2, $3, $4, 0, $5)`,
		newID("activity"),
		userID,
		fmt.Sprintf("Опубликована статья «%s»", request.Title),
		"Только что",
		time.Now().Unix(),
	); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update activity feed"})
	}

	if err = insertPortalNotification(
		requestContext,
		tx,
		userID,
		"Статья опубликована",
		fmt.Sprintf("Материал «%s» теперь доступен другим пользователям.", request.Title),
		"success",
	); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create notification"})
	}

	if err = tx.Commit(requestContext); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to save article"})
	}

	s.advanceAcceptedTask(requestContext, userID, "Задание на статью")

	article, err := s.loadArticle(requestContext, userID, articleID)
	if err != nil {
		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"articleId": articleID})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"articleId": articleID,
		"article":   article,
	})
}

func (s *Server) ToggleReaction(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	articleID := ctx.Params("articleId")
	var request reactionRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	request.Reaction = strings.TrimSpace(strings.ToLower(request.Reaction))
	if request.Reaction != "like" && request.Reaction != "dislike" && request.Reaction != "repost" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "unsupported reaction"})
	}

	requestContext, span := s.tracer.Start(ctx.UserContext(), "toggle-reaction")
	defer span.End()

	tx, err := s.db.BeginTx(requestContext, pgx.TxOptions{})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to start transaction"})
	}
	defer tx.Rollback(requestContext)

	var (
		articleTitle string
		authorID     string
	)
	if err = tx.QueryRow(
		requestContext,
		`select title, author_id from gamification.articles where id = $1`,
		articleID,
	).Scan(&articleTitle, &authorID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "article not found"})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update reaction"})
	}

	toggledOn := false
	switch request.Reaction {
	case "like", "dislike":
		var existing bool
		if err = tx.QueryRow(
			requestContext,
			`select exists(
				select 1 from gamification.article_reactions
				where article_id = $1 and user_id = $2 and reaction_type = $3
			)`,
			articleID,
			userID,
			request.Reaction,
		).Scan(&existing); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update reaction"})
		}

		if existing {
			if _, err = tx.Exec(
				requestContext,
				`delete from gamification.article_reactions
				 where article_id = $1 and user_id = $2 and reaction_type = $3`,
				articleID,
				userID,
				request.Reaction,
			); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update reaction"})
			}
		} else {
			opposite := "like"
			if request.Reaction == "like" {
				opposite = "dislike"
			}

			if _, err = tx.Exec(
				requestContext,
				`delete from gamification.article_reactions
				 where article_id = $1 and user_id = $2 and reaction_type = $3`,
				articleID,
				userID,
				opposite,
			); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update reaction"})
			}

			if _, err = tx.Exec(
				requestContext,
				`insert into gamification.article_reactions (article_id, user_id, reaction_type)
				 values ($1, $2, $3)`,
				articleID,
				userID,
				request.Reaction,
			); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update reaction"})
			}
			toggledOn = true
		}
	case "repost":
		commandTag, execErr := tx.Exec(
			requestContext,
			`delete from gamification.article_reactions
			 where article_id = $1 and user_id = $2 and reaction_type = 'repost'`,
			articleID,
			userID,
		)
		if execErr != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update reaction"})
		}
		if commandTag.RowsAffected() == 0 {
			if _, err = tx.Exec(
				requestContext,
				`insert into gamification.article_reactions (article_id, user_id, reaction_type)
				 values ($1, $2, 'repost')`,
				articleID,
				userID,
			); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update reaction"})
			}
			toggledOn = true
		}
	}

	if err = refreshArticleCardStats(requestContext, tx, articleID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update article statistics"})
	}

	if toggledOn && authorID != userID {
		if err = insertPortalNotification(
			requestContext,
			tx,
			authorID,
			"Новая реакция на статью",
			fmt.Sprintf("Пользователь отреагировал на статью «%s».", articleTitle),
			"info",
		); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create notification"})
		}
	}

	if err = tx.Commit(requestContext); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to save reaction"})
	}

	if toggledOn {
		s.advanceAcceptedTask(requestContext, userID, "Скаутинг")
	}

	article, err := s.loadArticle(requestContext, userID, articleID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "reaction saved but article refresh failed"})
	}

	return ctx.JSON(fiber.Map{"article": article})
}

func (s *Server) CreateComment(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	articleID := ctx.Params("articleId")
	var request createCommentRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	request.Body = strings.TrimSpace(request.Body)
	if request.Body == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "comment body is required"})
	}

	requestContext, span := s.tracer.Start(ctx.UserContext(), "create-comment")
	defer span.End()

	tx, err := s.db.BeginTx(requestContext, pgx.TxOptions{})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to start transaction"})
	}
	defer tx.Rollback(requestContext)

	var (
		articleTitle string
		authorID     string
	)
	if err = tx.QueryRow(
		requestContext,
		`select title, author_id from gamification.articles where id = $1`,
		articleID,
	).Scan(&articleTitle, &authorID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "article not found"})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create comment"})
	}

	if _, err = tx.Exec(
		requestContext,
		`insert into gamification.article_comments (id, article_id, author_id, body)
		 values ($1, $2, $3, $4)`,
		newID("comment"),
		articleID,
		userID,
		request.Body,
	); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create comment"})
	}

	if _, err = tx.Exec(
		requestContext,
		`update gamification.user_state
		 set comments_count = comments_count + 1,
		     updated_at = now()
		 where user_id = $1`,
		userID,
	); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update comment statistics"})
	}

	if _, err = tx.Exec(
		requestContext,
		`insert into gamification.recent_activity (id, user_id, title, timestamp_label, xp, sort_order)
		 values ($1, $2, $3, $4, 0, $5)`,
		newID("activity"),
		userID,
		fmt.Sprintf("Добавлен комментарий к статье «%s»", articleTitle),
		"Только что",
		time.Now().Unix(),
	); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update activity feed"})
	}

	if err = refreshArticleCardStats(requestContext, tx, articleID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update article statistics"})
	}

	if authorID != userID {
		if err = insertPortalNotification(
			requestContext,
			tx,
			authorID,
			"Новый комментарий",
			fmt.Sprintf("К статье «%s» добавлен новый комментарий.", articleTitle),
			"success",
		); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create notification"})
		}
	}

	if err = tx.Commit(requestContext); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to save comment"})
	}

	s.advanceAcceptedTask(requestContext, userID, "Корректор")
	s.advanceAcceptedTask(requestContext, userID, "Скаутинг")

	article, err := s.loadArticle(requestContext, userID, articleID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "comment saved but article refresh failed"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"article": article})
}

func (s *Server) DeleteComment(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	articleID := ctx.Params("articleId")
	commentID := ctx.Params("commentId")
	requestContext, span := s.tracer.Start(ctx.UserContext(), "delete-comment")
	defer span.End()

	tx, err := s.db.BeginTx(requestContext, pgx.TxOptions{})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to start transaction"})
	}
	defer tx.Rollback(requestContext)

	var commentAuthorID string
	err = tx.QueryRow(
		requestContext,
		`select author_id from gamification.article_comments
		 where id = $1 and article_id = $2`,
		commentID,
		articleID,
	).Scan(&commentAuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "comment not found"})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to delete comment"})
	}

	if commentAuthorID != userID {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "you can delete only your own comment"})
	}

	if _, err = tx.Exec(
		requestContext,
		`delete from gamification.article_comments
		 where id = $1 and article_id = $2 and author_id = $3`,
		commentID,
		articleID,
		userID,
	); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to delete comment"})
	}

	if _, err = tx.Exec(
		requestContext,
		`update gamification.user_state
		 set comments_count = greatest(comments_count - 1, 0),
		     updated_at = now()
		 where user_id = $1`,
		userID,
	); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update comment statistics"})
	}

	if err = refreshArticleCardStats(requestContext, tx, articleID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update article statistics"})
	}

	if err = tx.Commit(requestContext); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to delete comment"})
	}

	article, err := s.loadArticle(requestContext, userID, articleID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "comment deleted but article refresh failed"})
	}

	return ctx.JSON(fiber.Map{"article": article})
}

func (s *Server) loadArticles(ctx context.Context, viewerID string) ([]fiber.Map, error) {
	rows, err := s.db.Query(
		ctx,
		`select id
		 from gamification.articles
		 order by created_at desc, id desc`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articleIDs []string
	for rows.Next() {
		var articleID string
		if err := rows.Scan(&articleID); err != nil {
			return nil, err
		}
		articleIDs = append(articleIDs, articleID)
	}

	articles := make([]fiber.Map, 0, len(articleIDs))
	for _, articleID := range articleIDs {
		article, loadErr := s.loadArticle(ctx, viewerID, articleID)
		if loadErr != nil {
			return nil, loadErr
		}
		articles = append(articles, article)
	}

	return articles, rows.Err()
}

func (s *Server) loadArticle(ctx context.Context, viewerID, articleID string) (fiber.Map, error) {
	var (
		id           string
		title        string
		summary      string
		body         string
		authorID     string
		authorName   string
		authorXP     int32
		views        int32
		commentCount int32
		likes        int32
		dislikes     int32
		reposts      int32
		liked        bool
		disliked     bool
		reposted     bool
		createdAt    time.Time
	)

	err := s.db.QueryRow(
		ctx,
		`select
			a.id,
			a.title,
			a.summary,
			a.body,
			a.author_id,
			u.name,
			coalesce(us.current_xp, 0),
			coalesce(ac.views, 0),
			coalesce(ac.comments, 0),
			coalesce((select count(*) from gamification.article_reactions r where r.article_id = a.id and r.reaction_type = 'like'), 0),
			coalesce((select count(*) from gamification.article_reactions r where r.article_id = a.id and r.reaction_type = 'dislike'), 0),
			coalesce((select count(*) from gamification.article_reactions r where r.article_id = a.id and r.reaction_type = 'repost'), 0),
			exists(select 1 from gamification.article_reactions r where r.article_id = a.id and r.user_id = $1 and r.reaction_type = 'like'),
			exists(select 1 from gamification.article_reactions r where r.article_id = a.id and r.user_id = $1 and r.reaction_type = 'dislike'),
			exists(select 1 from gamification.article_reactions r where r.article_id = a.id and r.user_id = $1 and r.reaction_type = 'repost'),
			a.created_at
		 from gamification.articles a
		 join user_service.users u on u.id = a.author_id
		 left join gamification.user_state us on us.user_id = a.author_id
		 left join gamification.article_cards ac on ac.id = a.id
		 where a.id = $2`,
		viewerID,
		articleID,
	).Scan(
		&id,
		&title,
		&summary,
		&body,
		&authorID,
		&authorName,
		&authorXP,
		&views,
		&commentCount,
		&likes,
		&dislikes,
		&reposts,
		&liked,
		&disliked,
		&reposted,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	commentRows, err := s.db.Query(
		ctx,
		`select
			c.id,
			c.author_id,
			u.name,
			coalesce(us.current_xp, 0),
			c.body,
			c.created_at
		 from gamification.article_comments c
		 join user_service.users u on u.id = c.author_id
		 left join gamification.user_state us on us.user_id = c.author_id
		 where c.article_id = $1
		 order by c.created_at desc, c.id desc`,
		articleID,
	)
	if err != nil {
		return nil, err
	}
	defer commentRows.Close()

	comments := make([]fiber.Map, 0)
	for commentRows.Next() {
		var (
			commentID       string
			commentAuthorID string
			commentAuthor   string
			commentAuthorXP int32
			commentBody     string
			commentCreated  time.Time
		)

		if err := commentRows.Scan(
			&commentID,
			&commentAuthorID,
			&commentAuthor,
			&commentAuthorXP,
			&commentBody,
			&commentCreated,
		); err != nil {
			return nil, err
		}

		comments = append(comments, fiber.Map{
			"id":        commentID,
			"authorId":  commentAuthorID,
			"author":    commentAuthor,
			"level":     levelTextFromXP(commentAuthorXP),
			"timestamp": formatTimestampLabel(commentCreated),
			"body":      commentBody,
		})
	}

	if err := commentRows.Err(); err != nil {
		return nil, err
	}

	return fiber.Map{
		"id":      id,
		"title":   title,
		"summary": summary,
		"body":    splitArticleBody(body),
		"author": fiber.Map{
			"id":    authorID,
			"name":  authorName,
			"level": levelTextFromXP(authorXP),
		},
		"metrics": fiber.Map{
			"likes":    likes,
			"dislikes": dislikes,
			"reposts":  reposts,
			"views":    views,
		},
		"comments": comments,
		"publishedAt": formatDateLabel(createdAt),
		"viewerActions": fiber.Map{
			"liked":    liked,
			"disliked": disliked,
			"reposted": reposted,
		},
		"commentCount": commentCount,
	}, nil
}

func (s *Server) incrementArticleViews(ctx context.Context, articleID string) error {
	_, err := s.db.Exec(
		ctx,
		`update gamification.article_cards
		 set views = views + 1
		 where id = $1`,
		articleID,
	)
	return err
}

func (s *Server) advanceAcceptedTask(ctx context.Context, userID, taskTitle string) {
	var (
		taskID string
		status string
	)

	err := s.db.QueryRow(
		ctx,
		`select id, status
		 from gamification.tasks
		 where user_id = $1 and title = $2
		 order by id
		 limit 1`,
		userID,
		taskTitle,
	).Scan(&taskID, &status)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warn("load task for auto advance failed", "userId", userID, "taskTitle", taskTitle, "error", err)
		}
		return
	}

	if status != "in_progress" {
		return
	}

	if _, err = s.gameClient.AdvanceTask(ctx, &gamificationv1.AdvanceTaskRequest{
		UserId: userID,
		TaskId: taskID,
	}); err != nil {
		s.logger.Warn("auto advance task failed", "userId", userID, "taskId", taskID, "error", err)
	}
}

func seedUserState(ctx context.Context, tx pgx.Tx, userID string) error {
	if _, err := tx.Exec(
		ctx,
		`insert into gamification.user_state
		 (user_id, current_xp, coins, today_earned, streak_days, rank, completed_tasks, api_requests, articles_count, comments_count)
		 values ($1, 0, 0, 0, 0, 1, 0, 0, 0, 0)`,
		userID,
	); err != nil {
		return err
	}

	for index, dayCode := range []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"} {
		if _, err := tx.Exec(
			ctx,
			`insert into gamification.weekly_activity (user_id, day_code, xp, sort_order)
			 values ($1, $2, 0, $3)`,
			userID,
			dayCode,
			index+1,
		); err != nil {
			return err
		}
	}

	for _, task := range []struct {
		title       string
		description string
		target      int
		rewardXP    int
	}{
		{"Задание на статью", "Напиши и опубликуй первую статью в базе знаний платформы.", 1, 150},
		{"Скаутинг", "Сделай 3 полезных действия под статьями: лайк, дизлайк, репост или комментарий.", 3, 80},
		{"Корректор", "Оставь один содержательный комментарий к статье.", 1, 40},
	} {
		if _, err := tx.Exec(
			ctx,
			`insert into gamification.tasks (id, user_id, title, description, status, progress, target, reward_xp)
			 values ($1, $2, $3, $4, 'available', 0, $5, $6)`,
			newID("task"),
			userID,
			task.title,
			task.description,
			task.target,
			task.rewardXP,
		); err != nil {
			return err
		}
	}

	for _, achievement := range []struct {
		title       string
		description string
		rarity      string
		rewardXP    int
	}{
		{"Зоркий глаз", "Сделай первое полезное действие на платформе.", "Обычные", 20},
		{"Эмпат", "Оставь первый комментарий к статье.", "Обычные", 10},
		{"Автор", "Опубликуй первую статью.", "Редкие", 70},
		{"Страж", "Сделай 3 полезных действия под статьями.", "Эпические", 150},
	} {
		if _, err := tx.Exec(
			ctx,
			`insert into gamification.achievements (id, user_id, title, description, rarity, status, reward_xp)
			 values ($1, $2, $3, $4, $5, 'locked', $6)`,
			newID("achievement"),
			userID,
			achievement.title,
			achievement.description,
			achievement.rarity,
			achievement.rewardXP,
		); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(
		ctx,
		`insert into gamification.leaderboard (user_id, rank, xp)
		 values ($1, (select coalesce(max(rank), 0) + 1 from gamification.leaderboard), 0)
		 on conflict (user_id) do nothing`,
		userID,
	); err != nil {
		return err
	}

	for _, reward := range []struct {
		title       string
		description string
		cost        int
		category    string
	}{
		{"Фирменный худи", "Мерч команды платформы.", 2500, "мерч"},
		{"Кофе с архитектором", "Часовая встреча и разбор архитектурных решений.", 1800, "бонус"},
	} {
		if _, err := tx.Exec(
			ctx,
			`insert into gamification.rewards (id, user_id, title, description, cost, status, category)
			 values ($1, $2, $3, $4, $5, 'available', $6)`,
			newID("reward"),
			userID,
			reward.title,
			reward.description,
			reward.cost,
			reward.category,
		); err != nil {
			return err
		}
	}

	return insertPortalNotification(
		ctx,
		tx,
		userID,
		"Добро пожаловать!",
		"Профиль создан. Теперь можно брать задания, публиковать статьи и взаимодействовать с другими пользователями.",
		"success",
	)
}

func refreshArticleCardStats(ctx context.Context, tx pgx.Tx, articleID string) error {
	var likeCount, dislikeCount, repostCount, commentCount int32
	if err := tx.QueryRow(
		ctx,
		`select
			coalesce((select count(*) from gamification.article_reactions where article_id = $1 and reaction_type = 'like'), 0),
			coalesce((select count(*) from gamification.article_reactions where article_id = $1 and reaction_type = 'dislike'), 0),
			coalesce((select count(*) from gamification.article_reactions where article_id = $1 and reaction_type = 'repost'), 0),
			coalesce((select count(*) from gamification.article_comments where article_id = $1), 0)`,
		articleID,
	).Scan(&likeCount, &dislikeCount, &repostCount, &commentCount); err != nil {
		return err
	}

	xp := likeCount*5 + repostCount*10 + commentCount*3
	ratingValue := 4.6 + float64(likeCount)*0.06 + float64(repostCount)*0.08 - float64(dislikeCount)*0.07
	ratingValue = math.Max(3.8, math.Min(5.0, ratingValue))

	_, err := tx.Exec(
		ctx,
		`update gamification.article_cards
		 set comments = $2,
		     xp = $3,
		     rating = $4
		 where id = $1`,
		articleID,
		commentCount,
		xp,
		fmt.Sprintf("%.1f", ratingValue),
	)
	return err
}

func insertPortalNotification(ctx context.Context, tx pgx.Tx, userID, title, body, variant string) error {
	_, err := tx.Exec(
		ctx,
		`insert into gamification.notifications (id, user_id, title, body, variant)
		 values ($1, $2, $3, $4, $5)`,
		newID("notification"),
		userID,
		title,
		body,
		variant,
	)
	return err
}

func createTokenForUser(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func comparePassword(expectedHash, password string) bool {
	actualHash := hashPassword(password)
	return subtle.ConstantTimeCompare([]byte(expectedHash), []byte(actualHash)) == 1
}

func newID(prefix string) string {
	buffer := make([]byte, 8)
	if _, err := rand.Read(buffer); err != nil {
		return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	}

	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(buffer))
}

func splitArticleBody(body string) []string {
	parts := strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n")
	paragraphs := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			paragraphs = append(paragraphs, trimmed)
		}
	}
	if len(paragraphs) == 0 {
		return []string{body}
	}

	return paragraphs
}

func levelTextFromXP(currentXP int32) string {
	thresholds := []int32{0, 500, 1200, 2200, 3600, 5400, 7600, 10200, 13200, 16800}
	level := 1
	for index := 1; index < len(thresholds); index++ {
		if currentXP >= thresholds[index] {
			level = index + 1
		}
	}

	return fmt.Sprintf("Уровень %d", level)
}

func formatTimestampLabel(timestamp time.Time) string {
	age := time.Since(timestamp)
	if age < time.Minute {
		return "только что"
	}
	if age < time.Hour {
		return fmt.Sprintf("%d мин назад", int(age.Minutes()))
	}
	if age < 24*time.Hour {
		return fmt.Sprintf("%d ч назад", int(age.Hours()))
	}

	return timestamp.Format("02.01.2006 15:04")
}

func formatDateLabel(timestamp time.Time) string {
	now := time.Now()
	year, month, day := now.Date()
	articleYear, articleMonth, articleDay := timestamp.Date()
	if year == articleYear && month == articleMonth && day == articleDay {
		return "Сегодня"
	}

	return timestamp.Format("02.01.2006")
}
