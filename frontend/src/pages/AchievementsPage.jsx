import { useMemo, useState } from "react";

import { useAppData } from "../context/AppDataContext";

const filters = [
  { id: "all", label: "Все" },
  { id: "received", label: "Получены" },
  { id: "locked", label: "Заблокированы" },
];

function AchievementsPage() {
  const { portalState } = useAppData();
  const [filter, setFilter] = useState("all");

  const achievements = useMemo(() => {
    const items = portalState.achievements.items;

    if (filter === "received") {
      return items.filter((item) => item.status === "unlocked");
    }

    if (filter === "locked") {
      return items.filter((item) => item.status === "locked");
    }

    return items;
  }, [filter, portalState.achievements.items]);

  return (
    <section className="page">
      <div className="page__intro">
        <div>
          <h1>Достижения</h1>
          <p>Собери все бейджи и покажи свое мастерство.</p>
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

      <section className="bucket-grid">
        {portalState.achievements.buckets.map((bucket) => (
          <article key={bucket.id} className="bucket-card">
            <strong>
              {bucket.collected} / {bucket.total}
            </strong>
            <span>{bucket.label}</span>
            <div className="bucket-card__track">
              <div
                className={`bucket-card__fill bucket-card__fill--${bucket.accent}`}
                style={{ width: `${bucket.total === 0 ? 0 : (bucket.collected / bucket.total) * 100}%` }}
              />
            </div>
          </article>
        ))}
      </section>

      <section className="badge-grid">
        {achievements.map((achievement) => (
          <article
            key={achievement.id}
            className={`achievement-card${achievement.status === "locked" ? " achievement-card--locked" : ""}`}
          >
            <strong>{achievement.title}</strong>
            <span>{achievement.rarity}</span>
            <p>{achievement.description}</p>
            <em>+{achievement.rewardXp} XP</em>
          </article>
        ))}
      </section>
    </section>
  );
}

export default AchievementsPage;
