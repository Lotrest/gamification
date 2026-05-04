import { useEffect } from "react";

import { useAppData } from "../context/AppDataContext";

function ToastViewport() {
  const { portalState, dismissToast } = useAppData();
  const activeIds = portalState.ui.toastQueue;
  const activeToasts = activeIds
    .map((id) => portalState.notifications.find((item) => item.id === id))
    .filter(Boolean)
    .slice(0, 3);

  useEffect(() => {
    if (activeToasts.length === 0) {
      return undefined;
    }

    const timer = window.setTimeout(() => {
      dismissToast(activeToasts[0].id);
    }, 3500);

    return () => window.clearTimeout(timer);
  }, [activeToasts, dismissToast]);

  return (
    <div className="toast-viewport">
      {activeToasts.map((toast) => (
        <div key={toast.id} className={`toast toast--${toast.variant}`}>
          <strong>{toast.title}</strong>
          <p>{toast.body}</p>
        </div>
      ))}
    </div>
  );
}

export default ToastViewport;
