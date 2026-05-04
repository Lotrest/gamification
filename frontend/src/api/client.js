const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://127.0.0.1:18080";
const REQUEST_TIMEOUT_MS = 4000;
const TOKEN_STORAGE_KEY = "cdek-platform-token";

function getAuthHeaders(extraHeaders = {}) {
  const token = getAuthToken();
  return {
    ...extraHeaders,
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
}

async function fetchWithTimeout(url, options = {}) {
  const controller = new AbortController();
  const timeoutId = window.setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);

  try {
    return await fetch(url, {
      ...options,
      signal: controller.signal,
    });
  } finally {
    window.clearTimeout(timeoutId);
  }
}

async function parseResponse(response) {
  if (!response.ok) {
    const payload = await response.json().catch(() => ({}));
    throw new Error(payload.message || "Request failed");
  }

  return response.json();
}

async function postJSON(path, body) {
  const response = await fetchWithTimeout(`${API_BASE_URL}${path}`, {
    method: "POST",
    headers: getAuthHeaders({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify(body),
  });

  return parseResponse(response);
}

export function getAuthToken() {
  return window.localStorage.getItem(TOKEN_STORAGE_KEY) || "";
}

export function setAuthToken(token) {
  window.localStorage.setItem(TOKEN_STORAGE_KEY, token);
}

export function clearAuthToken() {
  window.localStorage.removeItem(TOKEN_STORAGE_KEY);
}

export async function registerUser(payload) {
  return postJSON("/api/v1/auth/register", payload);
}

export async function loginUser(payload) {
  return postJSON("/api/v1/auth/login", payload);
}

export async function fetchBootstrap() {
  const response = await fetchWithTimeout(`${API_BASE_URL}/api/v1/bootstrap`, {
    headers: getAuthHeaders(),
  });

  return parseResponse(response);
}

export async function acceptTask(taskId) {
  const response = await fetchWithTimeout(`${API_BASE_URL}/api/v1/tasks/${taskId}/accept`, {
    method: "POST",
    headers: getAuthHeaders(),
  });

  return parseResponse(response);
}

export async function redeemReward(rewardId) {
  const response = await fetchWithTimeout(`${API_BASE_URL}/api/v1/rewards/${rewardId}/redeem`, {
    method: "POST",
    headers: getAuthHeaders(),
  });

  return parseResponse(response);
}

export async function fetchArticles() {
  const response = await fetchWithTimeout(`${API_BASE_URL}/api/v1/articles`, {
    headers: getAuthHeaders(),
  });

  return parseResponse(response);
}

export async function fetchArticle(articleId) {
  const response = await fetchWithTimeout(`${API_BASE_URL}/api/v1/articles/${articleId}`, {
    headers: getAuthHeaders(),
  });

  return parseResponse(response);
}

export async function createArticle(payload) {
  return postJSON("/api/v1/articles", payload);
}

export async function toggleReaction(articleId, reaction) {
  return postJSON(`/api/v1/articles/${articleId}/reactions`, { reaction });
}

export async function createComment(articleId, body) {
  return postJSON(`/api/v1/articles/${articleId}/comments`, { body });
}

export async function deleteComment(articleId, commentId) {
  const response = await fetchWithTimeout(`${API_BASE_URL}/api/v1/articles/${articleId}/comments/${commentId}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });

  return parseResponse(response);
}
