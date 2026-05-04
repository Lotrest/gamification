import { useAppData } from "../context/AppDataContext";

function ShopPage() {
  const { portalState, redeemReward, busyKey } = useAppData();
  const { rewards, currentUser } = portalState;

  return (
    <section className="page">
      <div className="page__intro">
        <div>
          <h1>Магазин наград</h1>
          <p>Обменивай накопленные баллы на бонусы, мерч и внутренние привилегии.</p>
        </div>
      </div>

      <section className="shop-balance">
        <span>Доступный баланс</span>
        <strong>{new Intl.NumberFormat("ru-RU").format(currentUser.coins)} coins</strong>
      </section>

      <section className="reward-grid">
        {rewards.map((reward) => {
          const disabled = reward.status === "redeemed" || reward.cost > currentUser.coins;
          const loading = busyKey === `reward:${reward.id}:redeem`;

          return (
            <article key={reward.id} className="reward-card">
              <div className="reward-card__category">{reward.category}</div>
              <h2>{reward.title}</h2>
              <p>{reward.description}</p>
              <div className="reward-card__footer">
                <strong>{new Intl.NumberFormat("ru-RU").format(reward.cost)} coins</strong>
                <button
                  className="action-button action-button--small"
                  type="button"
                  disabled={disabled || loading}
                  onClick={() => redeemReward(reward.id)}
                >
                  {loading ? "Сохраняем..." : reward.status === "redeemed" ? "Получено" : "Обменять"}
                </button>
              </div>
            </article>
          );
        })}
      </section>
    </section>
  );
}

export default ShopPage;
