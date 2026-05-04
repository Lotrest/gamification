import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";

function ArticleDetailPage() {
  const { articleId } = useParams();
  const { portalState, viewArticle, reactToArticle, commentArticle, deleteComment } = useAppData();
  const [comment, setComment] = useState("");

  const article = portalState.articles.published.find((item) => item.id === articleId);
  const recommended = useMemo(
    () => portalState.articles.published.filter((item) => item.id !== articleId).slice(0, 2),
    [articleId, portalState.articles.published],
  );

  useEffect(() => {
    let active = true;

    if (articleId) {
      Promise.resolve(viewArticle(articleId)).catch(() => {
        if (!active) {
          return;
        }
      });
    }

    return () => {
      active = false;
    };
  }, [articleId]);

  if (!article) {
    return (
      <section className="page">
        <article className="surface-card">
          <div className="empty-state">
            <h1>Статья не найдена</h1>
            <p>Материал еще не опубликован или отсутствует в текущем состоянии платформы.</p>
            <Link className="action-button action-button--small" to="/articles">
              Вернуться к статьям
            </Link>
          </div>
        </article>
      </section>
    );
  }

  async function submitComment(event) {
    event.preventDefault();
    if (!comment.trim()) {
      return;
    }

    await commentArticle(article.id, comment);
    setComment("");
  }

  return (
    <section className="page page--article">
      <article className="article-page">
        <div className="article-page__author">
          <div className="article-page__avatar">{article.author.name.slice(0, 1)}</div>
          <div>
            <strong>{article.author.name}</strong>
            <span>{article.author.level}</span>
          </div>
        </div>

        <h1>{article.title}</h1>

        <div className="article-page__body">
          {article.body.map((paragraph, index) => (
            <p key={index}>{paragraph}</p>
          ))}
        </div>

        <div className="article-actions">
          <button
            className={`article-actions__button${article.viewerActions.liked ? " article-actions__button--active" : ""}`}
            type="button"
            onClick={() => void reactToArticle(article.id, "like")}
          >
            Лайк {article.metrics.likes}
          </button>
          <button
            className={`article-actions__button${article.viewerActions.disliked ? " article-actions__button--active" : ""}`}
            type="button"
            onClick={() => void reactToArticle(article.id, "dislike")}
          >
            Дизлайк {article.metrics.dislikes}
          </button>
          <button
            className={`article-actions__button${article.viewerActions.reposted ? " article-actions__button--active" : ""}`}
            type="button"
            onClick={() => void reactToArticle(article.id, "repost")}
          >
            Репост {article.metrics.reposts}
          </button>
          <span className="article-actions__counter">Просмотры {article.metrics.views}</span>
          <span className="article-actions__counter">Комментарии {article.comments.length}</span>
        </div>

        <section className="comments-block">
          <h2>Комментарии</h2>

          <form className="comment-form" onSubmit={submitComment}>
            <textarea
              rows="4"
              value={comment}
              placeholder="Напиши содержательный комментарий к статье"
              onChange={(event) => setComment(event.target.value)}
            />
            <div className="comment-form__actions">
              <button className="action-button action-button--small" type="submit">
                Отправить
              </button>
            </div>
          </form>

          <div className="stack-list">
            {article.comments.length === 0 ? (
              <div className="empty-state empty-state--small">
                <h2>Пока комментариев нет</h2>
                <p>Ты можешь начать обсуждение первым.</p>
              </div>
            ) : (
              article.comments.map((item) => {
                const isOwnComment = item.authorId === portalState.currentUser.id;

                return (
                  <article key={item.id} className="comment-card">
                    <div className="comment-card__header">
                      <div className="comment-card__author">
                        <div className="comment-card__avatar">{item.author.slice(0, 1)}</div>
                        <div>
                          <strong>{item.author}</strong>
                          <span>{item.level}</span>
                        </div>
                      </div>
                      <span>{item.timestamp}</span>
                    </div>
                    <p>{item.body}</p>
                    {isOwnComment ? (
                      <div className="comment-card__actions">
                        <button
                          className="ghost-button"
                          type="button"
                          onClick={() => void deleteComment(article.id, item.id)}
                        >
                          Удалить комментарий
                        </button>
                      </div>
                    ) : null}
                  </article>
                );
              })
            )}
          </div>
        </section>
      </article>

      <section className="surface-card">
        <div className="section-heading">
          <div>
            <h2>Рекомендуемые статьи</h2>
            <p>Что еще можно почитать по теме</p>
          </div>
        </div>

        <div className="article-list article-list--compact">
          {recommended.length === 0 ? (
            <div className="empty-state empty-state--small">
              <h2>Пока рекомендаций нет</h2>
              <p>Они появятся, когда на платформе станет больше материалов.</p>
            </div>
          ) : (
            recommended.map((item) => (
              <Link key={item.id} className="knowledge-card" to={`/articles/${item.id}`}>
                <div className="knowledge-card__icon">Статья</div>
                <div className="knowledge-card__body">
                  <h2>{item.title}</h2>
                  <p>{item.summary}</p>
                </div>
              </Link>
            ))
          )}
        </div>
      </section>
    </section>
  );
}

export default ArticleDetailPage;
