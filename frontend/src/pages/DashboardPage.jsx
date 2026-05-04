import { Link } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";

function DashboardPage() {
  const { portalState } = useAppData();
  const { currentUser, dashboard, tasks } = portalState;
  const maxXp = Math.max(...dashboard.weeklyActivity.map((item) => item.xp), 1);

  return (
    <section className="page">
      <div className="page__intro">
        <div>
          <h1>Дашборд</h1>
          <p>Новый пользователь стартует с нуля. Вся активность будет накапливаться постепенно.</p>
        </div>
      </div>

      <section className="hero-card hero-card--compact">
        <div className="hero-card__left">
          <div className="hero-card__avatar-ring">
            <strong>{currentUser.level}</strong>
            <span>ур.</span>
          </div>
          <div>
            <h2>{currentUser.name}</h2>
            <p>{currentUser.title}</p>
            <span>
              {currentUser.company} · Ранг #{currentUser.rank}
            </span>
          </div>
        </div>

        <div className="hero-card__progress">
          <div className="progress-label">
            <span>До следующего уровня</span>
            <strong>{currentUser.xpToNextLevel} XP</strong>
          </div>
          <div className="progress-bar">
            <div style={{ width: `${currentUser.progressPercent}%` }} />
          </div>
        </div>
      </section>

      <section className="metrics-grid">
        {dashboard.metrics.map((metric) => (
          <article key={metric.id} className="metric-card metric-card--compact">
            <span>{metric.label}</span>
            <strong>{metric.value}</strong>
            <p>{metric.caption}</p>
          </article>
        ))}
      </section>

      <section className="content-grid">
        <article className="surface-card">
          <div className="section-heading">
            <div>
              <h2>Активность за неделю</h2>
              <p>Пока все значения нулевые, потому что действий еще не было</p>
            </div>
          </div>

          <div className="chart chart--compact">
            {dashboard.weeklyActivity.map((item) => (
              <div key={item.day} className="chart__item">
                <span>{item.xp} XP</span>
                <div className="chart__bar" style={{ height: `${Math.max((item.xp / maxXp) * 140, 16)}px` }} />
                <p>{item.day}</p>
              </div>
            ))}
          </div>
        </article>

        <article className="surface-card">
          <div className="section-heading">
            <div>
              <h2>Задания</h2>
              <p>Прими задания и начни заполнять прогресс реальными действиями</p>
            </div>
            <Link className="ghost-link" to="/tasks">
              Открыть
            </Link>
          </div>

          <div className="stack-list">
            {tasks.map((task) => (
              <div key={task.id} className="compact-task">
                <div>
                  <div className="compact-task__head">
                    <strong>{task.title}</strong>
                    <span>{task.rewardXp} XP</span>
                  </div>
                  <p>{task.description}</p>
                </div>
                <div className="progress-bar progress-bar--soft">
                  <div style={{ width: `${(task.progress / task.target) * 100}%` }} />
                </div>
              </div>
            ))}
          </div>
        </article>
      </section>

      <section className="content-grid">
        <article className="surface-card">
          <div className="section-heading">
            <div>
              <h2>Статьи и руководства</h2>
              <p>Лента знаний начнет заполняться после публикаций</p>
            </div>
            <Link className="ghost-link" to="/articles">
              К статьям
            </Link>
          </div>

          <div className="stack-list">
            {dashboard.articles.length === 0 ? (
              <div className="empty-state empty-state--small">
                <h2>Статей пока нет</h2>
                <p>Первая статья создаст ленту и откроет новые сценарии для заданий.</p>
              </div>
            ) : (
              dashboard.articles.map((article) => (
                <Link key={article.id} className="article-inline" to={`/articles/${article.id}`}>
                  <div>
                    <strong>{article.title}</strong>
                    <p>
                      {new Intl.NumberFormat("ru-RU").format(article.views)} просмотров · {article.comments} комментариев
                    </p>
                  </div>
                  <span>+{article.xp} XP</span>
                </Link>
              ))
            )}
          </div>
        </article>

        <article className="surface-card">
          <div className="section-heading">
            <div>
              <h2>Последние действия</h2>
              <p>Здесь появится история после принятия заданий и первых публикаций</p>
            </div>
          </div>

          <div className="stack-list">
            {dashboard.recentActivity.length === 0 ? (
              <div className="empty-state empty-state--small">
                <h2>История пока пустая</h2>
                <p>Начни с любого действия на платформе, и блок сразу оживет.</p>
              </div>
            ) : (
              dashboard.recentActivity.map((item) => (
                <div key={item.id} className="activity-row">
                  <div>
                    <strong>{item.title}</strong>
                    <p>{item.timestamp}</p>
                  </div>
                  <span>+{item.xp} XP</span>
                </div>
              ))
            )}
          </div>
        </article>
      </section>
    </section>
  );
}

export default DashboardPage;
