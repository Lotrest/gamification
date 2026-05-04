import { Link } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";

function ProfilePage() {
  const { portalState } = useAppData();
  const { currentUser, dashboard } = portalState;

  return (
    <section className="page">
      <div className="page__intro">
        <div>
          <h1>Профиль</h1>
          <p>Ты только зашел на платформу, поэтому весь прогресс пока нулевой.</p>
        </div>
      </div>

      <section className="profile-hero">
        <div className="profile-hero__identity">
          <div className="profile-hero__avatar">{currentUser.name.slice(0, 1)}</div>
          <div>
            <h2>{currentUser.name}</h2>
            <p>{currentUser.title}</p>
            <span>
              На портале с {currentUser.joinedAt} · Ранг #{currentUser.rank}
            </span>
          </div>
        </div>

        <div className="profile-hero__level">
          <div className="profile-hero__level-ring">
            <strong>{currentUser.levelText}</strong>
            <span>{currentUser.progressPercent}%</span>
          </div>
          <em>{currentUser.currentXp} XP</em>
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
              <h2>Топ статей</h2>
              <p>После первой публикации здесь появятся твои материалы</p>
            </div>
            <Link className="ghost-link" to="/articles">
              Все статьи
            </Link>
          </div>

          <div className="stack-list">
            {dashboard.articles.length === 0 ? (
              <div className="empty-state empty-state--small">
                <h2>Пока статей нет</h2>
                <p>Создай первую статью, чтобы начать собирать просмотры и реакции.</p>
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
              <p>История активности появится после первых действий на платформе</p>
            </div>
          </div>

          <div className="stack-list">
            {dashboard.recentActivity.length === 0 ? (
              <div className="empty-state empty-state--small">
                <h2>Пока активность пустая</h2>
                <p>Прими задание или опубликуй материал, чтобы увидеть первые события.</p>
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

export default ProfilePage;
