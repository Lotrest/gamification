function NotificationCenter({ items, open, onClose }) {
  return (
    <div className={`notification-center${open ? " notification-center--open" : ""}`}>
      <div className="notification-center__header">
        <strong>Уведомления</strong>
        <button type="button" onClick={onClose}>
          Закрыть
        </button>
      </div>

      <div className="notification-center__list">
        {items.length === 0 ? (
          <div className="notification-center__empty">Пока уведомлений нет</div>
        ) : (
          items.map((item) => (
            <article
              key={item.id}
              className={`notification-center__item notification-center__item--${item.variant}`}
            >
              <strong>{item.title}</strong>
              <p>{item.body}</p>
            </article>
          ))
        )}
      </div>
    </div>
  );
}

export default NotificationCenter;
