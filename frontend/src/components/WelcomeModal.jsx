import { useNavigate } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";

function WelcomeModal() {
  const navigate = useNavigate();
  const { portalState, dismissWelcome, openTasksFromWelcome } = useAppData();

  if (!portalState.ui.welcomeModalOpen) {
    return null;
  }

  return (
    <div className="welcome-modal__backdrop">
      <div className="welcome-modal">
        <div className="welcome-modal__icon">OK</div>
        <h2>Добро пожаловать!</h2>
        <p>
          Вы успешно зарегистрировались. Для вас есть маленький подарок:
          посмотрите раздел «Задания».
        </p>
        <div className="welcome-modal__actions">
          <button className="ghost-button" type="button" onClick={() => dismissWelcome()}>
            Отмена
          </button>
          <button
            className="action-button"
            type="button"
            onClick={() => {
              openTasksFromWelcome();
              navigate("/tasks");
            }}
          >
            Посмотреть
          </button>
        </div>
      </div>
    </div>
  );
}

export default WelcomeModal;
