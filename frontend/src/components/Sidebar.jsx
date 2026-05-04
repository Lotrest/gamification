import { NavLink } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";
import logo from "../img/logo.png";

const navigationItems = [
  { to: "/dashboard", label: "Дашборд" },
  { to: "/profile", label: "Профиль" },
  { to: "/tasks", label: "Задания" },
  { to: "/achievements", label: "Достижения" },
  { to: "/rating", label: "Рейтинг" },
  { to: "/shop", label: "Магазин наград" },
  { to: "/purchases", label: "Покупки" },
];

function Sidebar() {
  const { portalState } = useAppData();
  const currentUser = portalState?.currentUser;

  return (
    <aside className="sidebar">
      <div className="sidebar__brand">
        <img src={logo} alt="CDEK" />
      </div>

      <div className="sidebar__caption">Навигация</div>

      <nav className="sidebar__nav">
        {navigationItems.map((item) => (
          <NavLink
            key={item.to}
            className={({ isActive }) => `sidebar__link${isActive ? " sidebar__link--active" : ""}`}
            to={item.to}
          >
            {item.label}
          </NavLink>
        ))}
      </nav>

      <div className="sidebar__streak">
        <div className="sidebar__streak-row">
          <span>Сегодня заработано</span>
          <strong>+{currentUser?.todayEarned ?? 0} xp</strong>
        </div>
        <div className="sidebar__streak-track">
          {Array.from({ length: 7 }).map((_, index) => (
            <span key={index} className="sidebar__streak-step" />
          ))}
        </div>
        <div className="sidebar__streak-label">{currentUser?.streakDays ?? 0}-дней серия</div>
      </div>
    </aside>
  );
}

export default Sidebar;
