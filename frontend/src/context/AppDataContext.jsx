import { createContext, useContext, useEffect, useMemo, useState } from "react";

import {
  acceptTask as acceptTaskRequest,
  clearAuthToken,
  createArticle as createArticleRequest,
  createComment as createCommentRequest,
  deleteComment as deleteCommentRequest,
  fetchArticle,
  fetchArticles,
  fetchBootstrap,
  getAuthToken,
  loginUser,
  redeemReward as redeemRewardRequest,
  registerUser,
  setAuthToken,
  toggleReaction as toggleReactionRequest,
} from "../api/client";
import {
  createInitialPortalState,
  deleteCommentInState,
  dismissToastInState,
  dismissWelcomeInState,
  openTasksFromWelcomeInState,
  saveDraftInState,
  toggleNotificationPanelInState,
} from "../data/mockPortalState";

const AppDataContext = createContext(null);

function enrichTasks(tasks) {
  const taskMeta = {
    "Задание на статью": {
      kind: "publish_article",
      actionLabel: "Написать статью",
      actionRoute: "/articles/new",
    },
    Скаутинг: {
      kind: "article_engagement",
      actionLabel: "К статьям",
      actionRoute: "/articles",
    },
    Корректор: {
      kind: "article_comment",
      actionLabel: "Открыть статьи",
      actionRoute: "/articles",
    },
  };

  return (tasks ?? []).map((task) => ({
    ...task,
    ...(taskMeta[task.title] ?? {}),
  }));
}

function mergeRemoteState(remoteState, previousState, articleItems, options = {}) {
  const baseState = createInitialPortalState(remoteState);
  const previousUnread = new Map(
    (previousState?.notifications ?? []).map((item) => [item.id, Boolean(item.unread)]),
  );

  const notifications = (remoteState?.notifications ?? []).map((item, index) => ({
    ...item,
    unread: previousUnread.has(item.id) ? previousUnread.get(item.id) : index < 3,
  }));

  return {
    ...baseState,
    currentUser: {
      ...baseState.currentUser,
      ...(remoteState?.currentUser ?? {}),
    },
    dashboard: {
      metrics: remoteState?.dashboard?.metrics ?? baseState.dashboard.metrics,
      weeklyActivity: remoteState?.dashboard?.weeklyActivity ?? baseState.dashboard.weeklyActivity,
      recentActivity: remoteState?.dashboard?.recentActivity ?? baseState.dashboard.recentActivity,
      articles: remoteState?.dashboard?.articles ?? baseState.dashboard.articles,
    },
    tasks: enrichTasks(remoteState?.tasks ?? baseState.tasks),
    achievements: remoteState?.achievements ?? baseState.achievements,
    leaderboard: remoteState?.leaderboard ?? baseState.leaderboard,
    rewards: remoteState?.rewards ?? baseState.rewards,
    purchases: remoteState?.purchases ?? baseState.purchases,
    notifications,
    articles: {
      published: articleItems ?? previousState?.articles?.published ?? [],
      drafts: previousState?.articles?.drafts ?? [],
    },
    ui: {
      toastQueue: notifications.filter((item) => item.unread).slice(0, 3).map((item) => item.id),
      notificationPanelOpen: previousState?.ui?.notificationPanelOpen ?? false,
      welcomeModalOpen: Boolean(options.showWelcome),
    },
  };
}

function replaceArticleInState(state, article) {
  if (!state || !article) {
    return state;
  }

  const exists = state.articles.published.some((item) => item.id === article.id);

  return {
    ...state,
    articles: {
      ...state.articles,
      published: exists
        ? state.articles.published.map((item) => (item.id === article.id ? article : item))
        : [article, ...state.articles.published],
    },
  };
}

export function AppDataProvider({ children }) {
  const [portalState, setPortalState] = useState(null);
  const [loading, setLoading] = useState(true);
  const [busyKey, setBusyKey] = useState("");
  const [authToken, setAuthTokenState] = useState(() => getAuthToken());

  async function hydratePortalState(options = {}) {
    const token = options.token ?? getAuthToken();
    if (!token) {
      setPortalState(null);
      setLoading(false);
      return;
    }

    const [remoteState, articlesResponse] = await Promise.all([fetchBootstrap(), fetchArticles()]);
    setPortalState((current) =>
      mergeRemoteState(remoteState, current, articlesResponse.articles ?? [], options),
    );
    setLoading(false);
  }

  useEffect(() => {
    let active = true;

    async function load() {
      const token = getAuthToken();
      if (!token) {
        if (active) {
          setLoading(false);
        }
        return;
      }

      try {
        await hydratePortalState({ token, showWelcome: false });
      } catch {
        if (active) {
          clearAuthToken();
          setAuthTokenState("");
          setPortalState(null);
          setLoading(false);
        }
      }
    }

    load();

    return () => {
      active = false;
    };
  }, []);

  function updateState(mutator) {
    setPortalState((current) => mutator(current));
  }

  const value = useMemo(
    () => ({
      portalState,
      loading,
      busyKey,
      isAuthenticated: Boolean(authToken),
      async register(payload) {
        setBusyKey("auth:register");
        try {
          const response = await registerUser(payload);
          setAuthToken(response.token);
          setAuthTokenState(response.token);
          await hydratePortalState({ token: response.token, showWelcome: true });
        } catch (error) {
          clearAuthToken();
          setAuthTokenState("");
          setPortalState(null);
          setLoading(false);
          throw error;
        } finally {
          setBusyKey("");
        }
      },
      async login(payload) {
        setBusyKey("auth:login");
        try {
          const response = await loginUser(payload);
          setAuthToken(response.token);
          setAuthTokenState(response.token);
          await hydratePortalState({ token: response.token, showWelcome: false });
        } catch (error) {
          clearAuthToken();
          setAuthTokenState("");
          setPortalState(null);
          setLoading(false);
          throw error;
        } finally {
          setBusyKey("");
        }
      },
      logout() {
        clearAuthToken();
        setAuthTokenState("");
        setPortalState(null);
      },
      async acceptTask(taskId) {
        setBusyKey(`task:${taskId}:accept`);
        try {
          await acceptTaskRequest(taskId);
          await hydratePortalState({ showWelcome: false });
        } finally {
          setBusyKey("");
        }
      },
      saveDraft(payload) {
        setBusyKey("article:draft");
        updateState((current) => saveDraftInState(current, payload));
        setBusyKey("");
      },
      async publishArticle(payload) {
        setBusyKey("article:publish");
        try {
          const response = await createArticleRequest(payload);
          await hydratePortalState({ showWelcome: false });
          return response.articleId;
        } finally {
          setBusyKey("");
        }
      },
      async viewArticle(articleId) {
        const response = await fetchArticle(articleId);
        updateState((current) => replaceArticleInState(current, response.article));
      },
      async reactToArticle(articleId, reaction) {
        setBusyKey(`article:${articleId}:reaction`);
        try {
          await toggleReactionRequest(articleId, reaction);
          await hydratePortalState({ showWelcome: false });
        } finally {
          setBusyKey("");
        }
      },
      async commentArticle(articleId, body) {
        setBusyKey(`article:${articleId}:comment`);
        try {
          await createCommentRequest(articleId, body);
          await hydratePortalState({ showWelcome: false });
        } finally {
          setBusyKey("");
        }
      },
      async deleteComment(articleId, commentId) {
        setBusyKey(`article:${articleId}:delete-comment`);
        try {
          await deleteCommentRequest(articleId, commentId);
          await hydratePortalState({ showWelcome: false });
        } finally {
          setBusyKey("");
        }
      },
      dismissToast(notificationId) {
        updateState((current) => dismissToastInState(current, notificationId));
      },
      toggleNotificationPanel() {
        updateState((current) => toggleNotificationPanelInState(current));
      },
      dismissWelcome() {
        updateState((current) => dismissWelcomeInState(current));
      },
      openTasksFromWelcome() {
        updateState((current) => openTasksFromWelcomeInState(current));
      },
      async redeemReward(rewardId) {
        setBusyKey(`reward:${rewardId}:redeem`);
        try {
          await redeemRewardRequest(rewardId);
          await hydratePortalState({ showWelcome: false });
        } finally {
          setBusyKey("");
        }
      },
      removeCommentLocally(articleId, commentId) {
        updateState((current) => deleteCommentInState(current, articleId, commentId));
      },
    }),
    [portalState, loading, busyKey, authToken],
  );

  return <AppDataContext.Provider value={value}>{children}</AppDataContext.Provider>;
}

export function useAppData() {
  const context = useContext(AppDataContext);
  if (!context) {
    throw new Error("useAppData must be used within AppDataProvider");
  }

  return context;
}
