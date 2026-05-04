import { useState } from "react";
import { useNavigate } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";

function ArticleEditorPage() {
  const navigate = useNavigate();
  const { saveDraft, publishArticle, busyKey } = useAppData();
  const [form, setForm] = useState({
    title: "",
    summary: "",
    body: "",
  });

  function updateField(field, value) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  function validate() {
    return form.title.trim() && form.summary.trim() && form.body.trim();
  }

  function handleSaveDraft() {
    if (!form.title.trim() && !form.body.trim()) {
      return;
    }

    saveDraft(form);
    navigate("/articles");
  }

  async function handlePublish() {
    if (!validate()) {
      return;
    }

    const articleId = await publishArticle(form);
    navigate(`/articles/${articleId}`);
  }

  return (
    <section className="page">
      <div className="page__intro">
        <div>
          <h1>Написать статью</h1>
          <p>Материал сначала можно сохранить как черновик, а затем опубликовать в базе знаний.</p>
        </div>
      </div>

      <section className="surface-card editor-card">
        <label className="field">
          <span>Заголовок</span>
          <input
            type="text"
            value={form.title}
            placeholder="Например: Как мы строим BFF для комьюнити-платформы"
            onChange={(event) => updateField("title", event.target.value)}
          />
        </label>

        <label className="field">
          <span>Краткое описание</span>
          <textarea
            rows="3"
            value={form.summary}
            placeholder="Коротко опиши, о чем статья и какую пользу она дает."
            onChange={(event) => updateField("summary", event.target.value)}
          />
        </label>

        <label className="field">
          <span>Текст статьи</span>
          <textarea
            rows="12"
            value={form.body}
            placeholder="Разделяй абзацы пустой строкой или переносами строк."
            onChange={(event) => updateField("body", event.target.value)}
          />
        </label>

        <div className="editor-card__actions">
          <button className="ghost-button" type="button" onClick={handleSaveDraft}>
            {busyKey === "article:draft" ? "Сохраняем..." : "Сохранить черновик"}
          </button>
          <button className="action-button" type="button" disabled={!validate()} onClick={handlePublish}>
            {busyKey === "article:publish" ? "Публикуем..." : "Опубликовать статью"}
          </button>
        </div>
      </section>
    </section>
  );
}

export default ArticleEditorPage;
