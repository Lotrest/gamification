import { useMemo, useState } from "react";
import { Link } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";

const filters = [
  { id: "all", label: "Все" },
  { id: "active", label: "Активные" },
  { id: "completed", label: "Выполненные" },
];

function TasksPage() {
  const { portalState, acceptTask, busyKey } = useAppData();
  const [filter, setFilter] = useState("all");

  const tasks = useMemo(() => {
    if (filter === "active") {
      return portalState.tasks.filter((task) => task.status === "in_progress");
    }

    if (filter === "completed") {
      return portalState.tasks.filter((task) => task.status === "completed");
    }

    return portalState.tasks;
  }, [filter, portalState.tasks]);

  return (
    <section className="page">
      <div className="page__intro">
        <div>
          <h1>Задания</h1>
          <p>Выполняй задания и зарабатывай XP. Прогресс обновляется автоматически по твоим действиям.</p>
        </div>
      </div>

      <div className="tabs">
        {filters.map((item) => (
          <button
            key={item.id}
            className={`pill-button${filter === item.id ? " pill-button--active" : ""}`}
            type="button"
            onClick={() => setFilter(item.id)}
          >
            {item.label}
          </button>
        ))}
      </div>

      <div className="stack-list stack-list--wide">
        {tasks.map((task) => {
          const progress = Math.round((task.progress / task.target) * 100);
          const isAcceptBusy = busyKey === `task:${task.id}:accept`;

          return (
            <article
              key={task.id}
              className={`task-card${task.status === "completed" ? " task-card--completed" : ""}`}
            >
              <div className="task-card__head">
                <div>
                  <strong>{task.title}</strong>
                  <p>{task.description}</p>
                </div>
                <span className="task-card__reward">{task.rewardXp} XP</span>
              </div>

              <div className="task-card__meta">
                <span>
                  {task.progress}/{task.target}
                </span>
                <strong>{progress}%</strong>
              </div>

              <div className="progress-bar progress-bar--soft progress-bar--large">
                <div style={{ width: `${progress}%` }} />
              </div>

              <div className="task-card__controls">
                <span className={`status-tag status-tag--${task.status}`}>
                  {task.status === "completed"
                    ? "Выполнено"
                    : task.status === "in_progress"
                      ? "В процессе"
                      : "Новое"}
                </span>

                {task.status === "available" ? (
                  <button
                    className="action-button action-button--small"
                    type="button"
                    disabled={busyKey !== "" && !isAcceptBusy}
                    onClick={() => acceptTask(task.id)}
                  >
                    {isAcceptBusy ? "Сохраняем..." : "Принять"}
                  </button>
                ) : null}

                {task.status === "in_progress" ? (
                  <Link className="action-button action-button--small" to={task.actionRoute}>
                    {task.actionLabel}
                  </Link>
                ) : null}

                {task.status === "completed" ? (
                  <span className="task-card__done">Система завершила задание автоматически</span>
                ) : null}
              </div>
            </article>
          );
        })}
      </div>
    </section>
  );
}

export default TasksPage;
