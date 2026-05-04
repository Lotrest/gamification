import { Link } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";

function ArticlesPage() {
  const { portalState } = useAppData();
  const publishedArticles = portalState.articles.published;
  const drafts = portalState.articles.drafts;

  return (
    <section className="page">
      <div className="page__intro page__intro--split">
        <div>
          <h1>Статьи и руководства</h1>
          <p>Здесь пока пусто. Первые материалы появятся после публикации.</p>
        </div>
        <Link className="action-button action-button--small" to="/articles/new">
          Написать статью
        </Link>
      </div>

      {drafts.length > 0 ? (
        <section className="surface-card">
          <div className="section-heading">
            <div>
              <h2>Черновики</h2>
              <p>Временное хранилище материалов перед публикацией</p>
            </div>
          </div>

          <div className="stack-list">
            {drafts.map((article) => (
              <article key={article.id} className="draft-card">
                <div>
                  <strong>{article.title}</strong>
                  <p>{article.summary || "Краткое описание пока не заполнено."}</p>
                </div>
                <span>{article.publishedAt}</span>
              </article>
            ))}
          </div>
        </section>
      ) : null}

      <section className="article-list">
        {publishedArticles.length === 0 ? (
          <article className="surface-card">
            <div className="empty-state">
              <h2>На платформе пока нет статей</h2>
              <p>Опубликуй первый материал, и здесь появится живая лента знаний.</p>
              <Link className="action-button action-button--small" to="/articles/new">
                Создать первую статью
              </Link>
            </div>
          </article>
        ) : (
          publishedArticles.map((article) => (
            <Link key={article.id} className="knowledge-card" to={`/articles/${article.id}`}>
              <div className="knowledge-card__icon">Статья</div>
              <div className="knowledge-card__body">
                <h2>{article.title}</h2>
                <p>{article.summary}</p>
                <div className="knowledge-card__meta">
                  <span>{article.author.name}</span>
                  <span>{article.metrics.views} просмотров</span>
                  <span>{article.comments.length} комментариев</span>
                </div>
              </div>
            </Link>
          ))
        )}
      </section>
    </section>
  );
}

export default ArticlesPage;
