import { useAppData } from "../context/AppDataContext";

function PurchasesPage() {
  const { portalState } = useAppData();

  return (
    <section className="page">
      <div className="page__intro">
        <div>
          <h1>Покупки</h1>
          <p>История обмена баллов на награды, бонусы и внутренние призы.</p>
        </div>
      </div>

      <section className="surface-card">
        <div className="stack-list">
          {portalState.purchases.length === 0 ? (
            <div className="empty-state empty-state--small">
              <h2>Пока покупок нет</h2>
              <p>Когда ты обменяешь coins на награду, она появится здесь.</p>
            </div>
          ) : (
            portalState.purchases.map((purchase) => (
              <div key={purchase.id} className="purchase-row">
                <div>
                  <strong>{purchase.title}</strong>
                  <p>{purchase.redeemedAt}</p>
                </div>
                <div className="purchase-row__meta">
                  <span>{purchase.status}</span>
                  <strong>{new Intl.NumberFormat("ru-RU").format(purchase.cost)} coins</strong>
                </div>
              </div>
            ))
          )}
        </div>
      </section>
    </section>
  );
}

export default PurchasesPage;
