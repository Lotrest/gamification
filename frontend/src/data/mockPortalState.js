const justNow = "только что";

function createNotification(title, body, variant = "info") {
  return {
    id: `notification-${crypto.randomUUID()}`,
    title,
    body,
    variant,
    createdAt: new Date().toISOString(),
    unread: true,
  };
}

function createArticle({
  id,
  title,
  summary,
  body,
  authorId,
  authorName,
  authorLevel,
  likes,
  dislikes,
  reposts,
  views,
  comments,
  tags,
  isPublished = true,
  publishedAt = "Сегодня",
}) {
  return {
    id,
    title,
    summary,
    body,
    author: {
      id: authorId,
      name: authorName,
      level: authorLevel,
    },
    metrics: {
      likes,
      dislikes,
      reposts,
      views,
    },
    comments,
    tags,
    isPublished,
    publishedAt,
    viewerActions: {
      liked: false,
      disliked: false,
      reposted: false,
    },
  };
}

const initialTasks = [
  {
    id: "task-first-article",
    title: "Задание на статью",
    description: "Напиши и опубликуй первую статью в базе знаний платформы.",
    status: "available",
    progress: 0,
    target: 1,
    rewardXp: 150,
    kind: "publish_article",
    actionLabel: "Написать статью",
    actionRoute: "/articles/new",
  },
  {
    id: "task-scouting",
    title: "Скаутинг",
    description: "Сделай 3 полезных действия под статьями: лайк, репост или комментарий.",
    status: "available",
    progress: 0,
    target: 3,
    rewardXp: 80,
    kind: "article_engagement",
    actionLabel: "К статьям",
    actionRoute: "/articles",
  },
  {
    id: "task-feedback",
    title: "Корректор",
    description: "Оставь один содержательный комментарий к статье.",
    status: "available",
    progress: 0,
    target: 1,
    rewardXp: 40,
    kind: "article_comment",
    actionLabel: "Открыть статьи",
    actionRoute: "/articles",
  },
];

const starterNotifications = [
  createNotification(
    "Добро пожаловать!",
    "Пока у тебя нулевой прогресс. Начни с заданий и первой статьи.",
    "success",
  ),
];

function clone(value) {
  return JSON.parse(JSON.stringify(value));
}

function computeLevel(currentXp) {
  const thresholds = [0, 500, 1200, 2200, 3600, 5400, 7600, 10200];
  let level = 1;

  for (let index = 1; index < thresholds.length; index += 1) {
    if (currentXp >= thresholds[index]) {
      level = index + 1;
    }
  }

  const currentBase = thresholds[level - 1];
  const nextBase = thresholds[Math.min(level, thresholds.length - 1)];
  const progressPercent =
    nextBase === currentBase
      ? 0
      : Math.round(((currentXp - currentBase) / (nextBase - currentBase)) * 100);

  return {
    level,
    levelText: `Ур.${level}`,
    progressPercent: Math.max(Math.min(progressPercent, 100), 0),
    xpToNextLevel: Math.max(nextBase - currentXp, 0),
  };
}

function sortLeaderboardRows(rows) {
  return [...rows].sort((left, right) => right.xp - left.xp).map((row, index) => ({
    ...row,
    rank: index + 1,
  }));
}

function pushNotification(state, title, body, variant = "info") {
  const notification = createNotification(title, body, variant);
  state.notifications.unshift(notification);
  state.ui.toastQueue.unshift(notification.id);
}

function appendActivity(state, title, xp = 0, timestamp = justNow) {
  state.dashboard.recentActivity.unshift({
    id: `activity-${crypto.randomUUID()}`,
    title,
    timestamp,
    xp,
  });
}

function calculateEngagementProgress(state) {
  return state.articles.published.reduce((total, article) => {
    let articleTotal = 0;

    if (article.viewerActions.liked) {
      articleTotal += 1;
    }

    if (article.viewerActions.reposted) {
      articleTotal += 1;
    }

    articleTotal += article.comments.filter((comment) => comment.authorId === state.currentUser.id).length;

    return total + articleTotal;
  }, 0);
}

function calculateCommentProgress(state) {
  return state.articles.published.reduce(
    (total, article) => total + article.comments.filter((comment) => comment.authorId === state.currentUser.id).length,
    0,
  );
}

function calculatePublishedArticleProgress(state) {
  return state.articles.published.filter((article) => article.author.id === state.currentUser.id).length;
}

function syncCurrentUser(state) {
  const levelState = computeLevel(state.currentUser.currentXp);
  state.currentUser.level = levelState.level;
  state.currentUser.levelText = levelState.levelText;
  state.currentUser.progressPercent = levelState.progressPercent;
  state.currentUser.xpToNextLevel = levelState.xpToNextLevel;
  state.currentUser.totalBadges = state.achievements.items.filter((item) => item.status === "unlocked").length;
  state.currentUser.completedTasks = state.tasks.filter((task) => task.status === "completed").length;
}

function syncDashboard(state) {
  const publishedArticles = state.articles.published;
  const totalComments = publishedArticles.reduce((sum, article) => sum + article.comments.length, 0);

  state.dashboard.metrics = [
    { id: "api", label: "API-запросов", value: "0", caption: "пока нет активности" },
    { id: "articles", label: "Статей", value: String(publishedArticles.length), caption: "добавятся после публикации" },
    { id: "comments", label: "Комментариев", value: String(totalComments), caption: "появятся после обсуждений" },
    { id: "xp", label: "Всего XP", value: String(state.currentUser.currentXp), caption: "пока прогресс нулевой" },
  ];

  state.dashboard.articles = publishedArticles.slice(0, 3).map((article) => ({
    id: article.id,
    title: article.title,
    views: article.metrics.views,
    comments: article.comments.length,
    xp: article.metrics.likes * 5 + article.metrics.reposts * 10,
    rating: `${Math.max(4.6, 4.9 - article.metrics.dislikes * 0.1).toFixed(1)}`,
  }));
}

function syncAchievements(state) {
  const publishTask = state.tasks.find((task) => task.id === "task-first-article");
  const engageTask = state.tasks.find((task) => task.id === "task-scouting");
  const commentTask = state.tasks.find((task) => task.id === "task-feedback");

  state.achievements.items = state.achievements.items.map((achievement) => {
    if (achievement.id === "achievement-author" && publishTask?.status === "completed") {
      return { ...achievement, status: "unlocked" };
    }

    if (achievement.id === "achievement-guardian" && engageTask?.status === "completed") {
      return { ...achievement, status: "unlocked" };
    }

    if (achievement.id === "achievement-empathy" && commentTask?.status === "completed") {
      return { ...achievement, status: "unlocked" };
    }

    return achievement;
  });

  const buckets = [
    { id: "common", label: "Обычные", accent: "green", rarity: "Обычные" },
    { id: "rare", label: "Редкие", accent: "gray", rarity: "Редкие" },
    { id: "epic", label: "Эпические", accent: "green", rarity: "Эпические" },
    { id: "legendary", label: "Легендарные", accent: "gray", rarity: "Легендарные" },
  ];

  state.achievements.buckets = buckets.map((bucket) => {
    const items = state.achievements.items.filter((item) => item.rarity === bucket.rarity);

    return {
      id: bucket.id,
      label: bucket.label,
      accent: bucket.accent,
      total: items.length,
      collected: items.filter((item) => item.status === "unlocked").length,
    };
  });
}

function syncLeaderboard(state) {
  const rows = state.leaderboard.rows.map((row) =>
    row.userId === state.currentUser.id
      ? {
          ...row,
          name: state.currentUser.name,
          title: state.currentUser.title,
          company: state.currentUser.company,
          xp: state.currentUser.currentXp,
          levelText: `Lv.${state.currentUser.level}`,
          isCurrent: true,
        }
      : row,
  );

  state.leaderboard.rows = sortLeaderboardRows(rows);
  state.currentUser.rank = state.leaderboard.rows.find((row) => row.userId === state.currentUser.id)?.rank ?? 1;
  state.leaderboard.podium = state.leaderboard.rows.slice(0, 3);
}

function rewardTask(state, taskId) {
  const task = state.tasks.find((item) => item.id === taskId);
  if (!task || task.status === "completed") {
    return;
  }

  task.progress = task.target;
  task.status = "completed";
  state.currentUser.currentXp += task.rewardXp;
  state.currentUser.coins += task.rewardXp * 4;
  state.currentUser.todayEarned += task.rewardXp;

  appendActivity(state, `Задание «${task.title}» выполнено`, task.rewardXp);
  pushNotification(
    state,
    "Задание завершено",
    `Ты выполнил задание «${task.title}» и получил ${task.rewardXp} XP.`,
    "success",
  );

  syncCurrentUser(state);
  syncAchievements(state);
  syncLeaderboard(state);
  syncDashboard(state);
}

function syncTaskProgress(state) {
  const progressByKind = {
    publish_article: calculatePublishedArticleProgress(state),
    article_engagement: calculateEngagementProgress(state),
    article_comment: calculateCommentProgress(state),
  };

  state.tasks.forEach((task) => {
    if (task.status !== "in_progress") {
      return;
    }

    task.progress = Math.min(progressByKind[task.kind] ?? 0, task.target);
  });

  state.tasks
    .filter((task) => task.status === "in_progress" && task.progress >= task.target)
    .forEach((task) => rewardTask(state, task.id));
}

export function createInitialPortalState(remoteState) {
  const currentUser = {
    id: remoteState?.currentUser?.id ?? "me",
    name: remoteState?.currentUser?.name ?? "Пользователь",
    title: remoteState?.currentUser?.title ?? "QA-инженер",
    company: remoteState?.currentUser?.company ?? "CDEK Digital",
    joinedAt: remoteState?.currentUser?.joinedAt ?? "сегодня",
    location: remoteState?.currentUser?.location ?? "Новосибирск",
    team: remoteState?.currentUser?.team ?? "Платформа комьюнити",
    currentXp: 0,
    coins: 0,
    todayEarned: 0,
    streakDays: 0,
    rank: 1,
    level: 1,
    levelText: "Ур.1",
    progressPercent: 0,
    xpToNextLevel: 0,
    totalBadges: 0,
    completedTasks: 0,
  };

  const state = {
    currentUser,
    dashboard: {
      metrics: [],
      weeklyActivity: [
        { day: "Пн", xp: 0 },
        { day: "Вт", xp: 0 },
        { day: "Ср", xp: 0 },
        { day: "Чт", xp: 0 },
        { day: "Пт", xp: 0 },
        { day: "Сб", xp: 0 },
        { day: "Вс", xp: 0 },
      ],
      recentActivity: [],
      articles: [],
    },
    tasks: clone(initialTasks),
    achievements: {
      buckets: [],
      items: [
        {
          id: "achievement-eye",
          title: "Зоркий глаз",
          description: "Сделай первое полезное действие на платформе",
          rarity: "Обычные",
          status: "locked",
          rewardXp: 20,
        },
        {
          id: "achievement-empathy",
          title: "Эмпат",
          description: "Оставь первый комментарий к статье",
          rarity: "Обычные",
          status: "locked",
          rewardXp: 10,
        },
        {
          id: "achievement-author",
          title: "Автор",
          description: "Опубликуй первую статью",
          rarity: "Редкие",
          status: "locked",
          rewardXp: 70,
        },
        {
          id: "achievement-guardian",
          title: "Страж",
          description: "Сделай 3 полезных действия под статьями",
          rarity: "Эпические",
          status: "locked",
          rewardXp: 150,
        },
      ],
    },
    leaderboard: {
      podium: [],
      rows: [
        {
          userId: "me",
          rank: 1,
          xp: currentUser.currentXp,
          levelText: "Lv.1",
          name: currentUser.name,
          title: currentUser.title,
          company: currentUser.company,
          isCurrent: true,
        },
      ],
    },
    rewards: remoteState?.rewards ?? [
      {
        id: "reward-hoodie",
        title: "Фирменный худи",
        description: "Мерч команды платформы",
        cost: 2500,
        status: "available",
        category: "мерч",
      },
      {
        id: "reward-coffee",
        title: "Кофе с архитектором",
        description: "Часовая встреча и разбор архитектурных решений",
        cost: 1800,
        status: "available",
        category: "бонус",
      },
    ],
    purchases: remoteState?.purchases ?? [],
    notifications: starterNotifications,
    articles: {
      published: [],
      drafts: [],
    },
    ui: {
      toastQueue: [],
      notificationPanelOpen: false,
      welcomeModalOpen: true,
    },
  };

  syncCurrentUser(state);
  syncAchievements(state);
  syncLeaderboard(state);
  syncDashboard(state);

  return state;
}

export function dismissToastInState(state, notificationId) {
  const draft = clone(state);
  draft.ui.toastQueue = draft.ui.toastQueue.filter((id) => id !== notificationId);
  return draft;
}

export function toggleNotificationPanelInState(state) {
  const draft = clone(state);
  draft.ui.notificationPanelOpen = !draft.ui.notificationPanelOpen;
  if (draft.ui.notificationPanelOpen) {
    draft.notifications = draft.notifications.map((item) => ({ ...item, unread: false }));
  }
  return draft;
}

export function dismissWelcomeInState(state) {
  const draft = clone(state);
  draft.ui.welcomeModalOpen = false;
  return draft;
}

export function openTasksFromWelcomeInState(state) {
  const draft = clone(state);
  draft.ui.welcomeModalOpen = false;
  pushNotification(
    draft,
    "Вам назначены задания",
    "Переходи в раздел «Задания» и начинай собирать первый прогресс.",
    "success",
  );
  return draft;
}

export function acceptTaskInState(state, taskId) {
  const draft = clone(state);
  const task = draft.tasks.find((item) => item.id === taskId);
  if (!task || task.status !== "available") {
    return draft;
  }

  task.status = "in_progress";
  syncTaskProgress(draft);
  pushNotification(draft, "Вы приняли задание!", `Задание «${task.title}» перемещено в активные.`, "success");
  appendActivity(draft, `Принято задание «${task.title}»`);
  return draft;
}

export function saveDraftInState(state, payload) {
  const draft = clone(state);
  draft.articles.drafts.unshift(
    createArticle({
      id: `draft-${crypto.randomUUID()}`,
      title: payload.title,
      summary: payload.summary,
      body: payload.body
        .split("\n")
        .map((paragraph) => paragraph.trim())
        .filter(Boolean),
      authorId: draft.currentUser.id,
      authorName: draft.currentUser.name,
      authorLevel: draft.currentUser.levelText,
      likes: 0,
      dislikes: 0,
      reposts: 0,
      views: 0,
      comments: [],
      tags: ["черновик"],
      isPublished: false,
      publishedAt: "Черновик",
    }),
  );

  pushNotification(draft, "Черновик сохранен", "Материал сохранен и доступен в черновиках.", "info");
  return draft;
}

export function publishArticleInState(state, payload) {
  const draft = clone(state);
  const body = payload.body
    .split("\n")
    .map((paragraph) => paragraph.trim())
    .filter(Boolean);

  const article = createArticle({
    id: `article-${crypto.randomUUID()}`,
    title: payload.title,
    summary: payload.summary,
    body,
    authorId: draft.currentUser.id,
    authorName: draft.currentUser.name,
    authorLevel: draft.currentUser.levelText,
    likes: 0,
    dislikes: 0,
    reposts: 0,
    views: 1,
    comments: [],
    tags: ["новая статья", "комьюнити"],
  });

  draft.articles.published.unshift(article);
  draft.currentUser.currentXp += 120;
  draft.currentUser.coins += 140;
  draft.currentUser.todayEarned += 120;
  appendActivity(draft, `Опубликована статья «${article.title}»`, 120);
  pushNotification(
    draft,
    "Статья опубликована",
    "Материал появился в базе знаний и теперь доступен для реакций и комментариев.",
    "success",
  );
  syncTaskProgress(draft);
  syncCurrentUser(draft);
  syncAchievements(draft);
  syncLeaderboard(draft);
  syncDashboard(draft);

  return {
    state: draft,
    articleId: article.id,
  };
}

export function viewArticleInState(state, articleId) {
  const draft = clone(state);
  draft.articles.published = draft.articles.published.map((article) =>
    article.id === articleId
      ? {
          ...article,
          metrics: {
            ...article.metrics,
            views: article.metrics.views + 1,
          },
        }
      : article,
  );
  syncDashboard(draft);
  return draft;
}

export function reactToArticleInState(state, articleId, reaction) {
  const draft = clone(state);

  draft.articles.published = draft.articles.published.map((article) => {
    if (article.id !== articleId) {
      return article;
    }

    const nextArticle = clone(article);

    if (reaction === "like") {
      if (nextArticle.viewerActions.liked) {
        nextArticle.viewerActions.liked = false;
        nextArticle.metrics.likes = Math.max(0, nextArticle.metrics.likes - 1);
        pushNotification(draft, "Лайк отменен", `Ты убрал лайк со статьи «${article.title}».`, "info");
      } else {
        nextArticle.viewerActions.liked = true;
        nextArticle.metrics.likes += 1;

        if (nextArticle.viewerActions.disliked) {
          nextArticle.viewerActions.disliked = false;
          nextArticle.metrics.dislikes = Math.max(0, nextArticle.metrics.dislikes - 1);
        }

        pushNotification(draft, "Лайк поставлен", `Ты поддержал статью «${article.title}».`, "info");
      }
    }

    if (reaction === "dislike") {
      if (nextArticle.viewerActions.disliked) {
        nextArticle.viewerActions.disliked = false;
        nextArticle.metrics.dislikes = Math.max(0, nextArticle.metrics.dislikes - 1);
        pushNotification(draft, "Дизлайк отменен", `Ты убрал дизлайк со статьи «${article.title}».`, "info");
      } else {
        nextArticle.viewerActions.disliked = true;
        nextArticle.metrics.dislikes += 1;

        if (nextArticle.viewerActions.liked) {
          nextArticle.viewerActions.liked = false;
          nextArticle.metrics.likes = Math.max(0, nextArticle.metrics.likes - 1);
        }

        pushNotification(draft, "Дизлайк поставлен", `Ты оставил сигнал по статье «${article.title}».`, "info");
      }
    }

    if (reaction === "repost") {
      if (nextArticle.viewerActions.reposted) {
        nextArticle.viewerActions.reposted = false;
        nextArticle.metrics.reposts = Math.max(0, nextArticle.metrics.reposts - 1);
        pushNotification(draft, "Репост отменен", `Ты отменил репост статьи «${article.title}».`, "info");
      } else {
        nextArticle.viewerActions.reposted = true;
        nextArticle.metrics.reposts += 1;
        pushNotification(draft, "Репост выполнен", `Статья «${article.title}» отправлена коллегам.`, "success");
      }
    }

    return nextArticle;
  });

  syncTaskProgress(draft);
  syncDashboard(draft);
  return draft;
}

export function commentArticleInState(state, articleId, body) {
  const draft = clone(state);
  const trimmedBody = body.trim();
  if (!trimmedBody) {
    return draft;
  }

  draft.articles.published = draft.articles.published.map((article) => {
    if (article.id !== articleId) {
      return article;
    }

    return {
      ...article,
      comments: [
        {
          id: `comment-${crypto.randomUUID()}`,
          authorId: draft.currentUser.id,
          author: draft.currentUser.name,
          level: draft.currentUser.levelText,
          timestamp: justNow,
          body: trimmedBody,
        },
        ...article.comments,
      ],
    };
  });

  draft.currentUser.currentXp += 20;
  draft.currentUser.coins += 20;
  draft.currentUser.todayEarned += 20;
  appendActivity(draft, "Оставлен комментарий к статье", 20);
  pushNotification(
    draft,
    "Комментарий опубликован",
    "Твой комментарий добавлен и влияет на прогресс заданий.",
    "success",
  );
  syncTaskProgress(draft);
  syncCurrentUser(draft);
  syncAchievements(draft);
  syncLeaderboard(draft);
  syncDashboard(draft);
  return draft;
}

export function deleteCommentInState(state, articleId, commentId) {
  const draft = clone(state);
  let removedOwnComment = false;

  draft.articles.published = draft.articles.published.map((article) => {
    if (article.id !== articleId) {
      return article;
    }

    const nextComments = article.comments.filter((comment) => {
      const remove = comment.id === commentId && comment.authorId === draft.currentUser.id;
      if (remove) {
        removedOwnComment = true;
      }
      return !remove;
    });

    return {
      ...article,
      comments: nextComments,
    };
  });

  if (!removedOwnComment) {
    return draft;
  }

  draft.currentUser.currentXp = Math.max(0, draft.currentUser.currentXp - 20);
  draft.currentUser.coins = Math.max(0, draft.currentUser.coins - 20);
  draft.currentUser.todayEarned = Math.max(0, draft.currentUser.todayEarned - 20);
  pushNotification(draft, "Комментарий удален", "Твой комментарий убран из статьи.", "info");
  syncTaskProgress(draft);
  syncCurrentUser(draft);
  syncAchievements(draft);
  syncLeaderboard(draft);
  syncDashboard(draft);
  return draft;
}

export function redeemRewardInState(state, rewardId) {
  const draft = clone(state);
  const reward = draft.rewards.find((item) => item.id === rewardId);

  if (!reward || reward.status === "redeemed" || reward.cost > draft.currentUser.coins) {
    return draft;
  }

  reward.status = "redeemed";
  draft.currentUser.coins -= reward.cost;
  draft.purchases.unshift({
    id: `purchase-${crypto.randomUUID()}`,
    title: reward.title,
    redeemedAt: "сегодня",
    status: "Ожидает выдачи",
    cost: reward.cost,
  });
  pushNotification(draft, "Награда оформлена", `Награда «${reward.title}» добавлена в покупки.`, "success");
  appendActivity(draft, `Оформлена награда «${reward.title}»`);
  return draft;
}
