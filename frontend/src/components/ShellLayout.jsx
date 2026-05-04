import { Navigate, Outlet } from "react-router-dom";

import { useAppData } from "../context/AppDataContext";
import Sidebar from "./Sidebar";
import ToastViewport from "./ToastViewport";
import Topbar from "./Topbar";
import WelcomeModal from "./WelcomeModal";

function ShellLayout() {
  const { portalState, loading, isAuthenticated } = useAppData();

  if (!loading && !isAuthenticated) {
    return <Navigate to="/auth" replace />;
  }

  if (loading || !portalState) {
    return <div className="screen-center">Загружаем платформу...</div>;
  }

  return (
    <div className="app-shell">
      <Sidebar />

      <main className="app-main">
        <Topbar />
        <Outlet />
      </main>

      <ToastViewport />
      <WelcomeModal />
    </div>
  );
}

export default ShellLayout;
