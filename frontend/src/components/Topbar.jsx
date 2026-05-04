import { Link } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";
import userAvatar from "../img/user.png";
import NotificationCenter from "./NotificationCenter";

function BellIcon() {
  return (
    <svg aria-hidden="true" viewBox="0 0 24 24">
      <path
        d="M12 4a4 4 0 0 0-4 4v2.7c0 .8-.3 1.5-.8 2.1L5.7 14a1 1 0 0 0 .8 1.7h11a1 1 0 0 0 .8-1.7l-1.5-1.2a3.2 3.2 0 0 1-.8-2.1V8a4 4 0 0 0-4-4Z"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M10 18a2.1 2.1 0 0 0 4 0"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.8"
        strokeLinecap="round"
      />
    </svg>
  );
}

function Topbar() {
  const { portalState, toggleNotificationPanel, logout } = useAppData();
  const currentUser = portalState.currentUser;
  const unreadCount = portalState.notifications.filter((item) => item.unread).length;

  return (
    <header className="topbar">
      <div className="topbar__left">
        <Link className="outline-button" to="/articles">
          Статьи и руководства
        </Link>
        <Link className="action-button action-button--small" to="/articles/new">
          Написать статью
        </Link>
      </div>

      <div className="topbar__right">
        <div className="topbar__coins">
          <span className="topbar__coin-dot" />
          {new Intl.NumberFormat("ru-RU").format(currentUser.coins)}
        </div>

        <button
          className="topbar__icon-button topbar__icon-button--bell"
          type="button"
          aria-label="Уведомления"
          onClick={() => toggleNotificationPanel()}
        >
          <BellIcon />
          {unreadCount > 0 ? <span className="topbar__badge">{unreadCount}</span> : null}
        </button>

        <button className="ghost-button topbar__logout" type="button" onClick={() => logout()}>
          Выйти
        </button>

        <div className="topbar__profile">
          <img src={userAvatar} alt={currentUser.name} />
          <div>
            <p>{currentUser.name}</p>
            <span>Уровень {currentUser.level}</span>
          </div>
        </div>
      </div>

      <NotificationCenter
        items={portalState.notifications}
        open={portalState.ui.notificationPanelOpen}
        onClose={() => toggleNotificationPanel()}
      />
    </header>
  );
}

export default Topbar;
