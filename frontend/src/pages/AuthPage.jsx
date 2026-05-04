import { useState } from "react";
import { Navigate } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";

function AuthPage() {
  const { isAuthenticated, register, login, busyKey } = useAppData();
  const [mode, setMode] = useState("login");
  const [error, setError] = useState("");
  const [form, setForm] = useState({
    name: "",
    title: "",
    email: "",
    password: "",
  });

  if (isAuthenticated) {
    return <Navigate to="/profile" replace />;
  }

  function updateField(field, value) {
    setForm((current) => ({
      ...current,
      [field]: value,
    }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");

    try {
      if (mode === "register") {
        await register(form);
      } else {
        await login({
          email: form.email,
          password: form.password,
        });
      }
    } catch (submitError) {
      setError(submitError.message || "Не удалось выполнить вход");
    }
  }

  return (
    <section className="auth-shell">
      <div className="auth-card">
        <div className="auth-card__intro">
          <p className="auth-card__eyebrow">CDEK Gamification</p>
          <h1>Вход и регистрация</h1>
          <p>
            Создай несколько пользователей, публикуй статьи, комментируй материалы коллег и смотри,
            как всё сохраняется в общей базе данных.
          </p>
        </div>

        <div className="auth-toggle">
          <button
            className={mode === "login" ? "auth-toggle__button auth-toggle__button--active" : "auth-toggle__button"}
            type="button"
            onClick={() => setMode("login")}
          >
            Вход
          </button>
          <button
            className={mode === "register" ? "auth-toggle__button auth-toggle__button--active" : "auth-toggle__button"}
            type="button"
            onClick={() => setMode("register")}
          >
            Регистрация
          </button>
        </div>

        <form className="auth-form" onSubmit={handleSubmit}>
          {mode === "register" ? (
            <>
              <label className="field">
                <span>Имя</span>
                <input
                  type="text"
                  value={form.name}
                  placeholder="Например, Анна Иванова"
                  onChange={(event) => updateField("name", event.target.value)}
                />
              </label>

              <label className="field">
                <span>Роль</span>
                <input
                  type="text"
                  value={form.title}
                  placeholder="QA-инженер, frontend-разработчик, аналитик"
                  onChange={(event) => updateField("title", event.target.value)}
                />
              </label>
            </>
          ) : null}

          <label className="field">
            <span>Email</span>
            <input
              type="email"
              value={form.email}
              placeholder="user@cdek.ru"
              onChange={(event) => updateField("email", event.target.value)}
            />
          </label>

          <label className="field">
            <span>Пароль</span>
            <input
              type="password"
              value={form.password}
              placeholder="Минимум один пароль для демо"
              onChange={(event) => updateField("password", event.target.value)}
            />
          </label>

          {error ? <div className="auth-form__error">{error}</div> : null}

          <button className="action-button" type="submit">
            {busyKey === "auth:register" || busyKey === "auth:login"
              ? "Сохраняем..."
              : mode === "register"
                ? "Создать пользователя"
                : "Войти"}
          </button>
        </form>
      </div>
    </section>
  );
}

export default AuthPage;
