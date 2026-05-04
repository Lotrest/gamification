import { Navigate, Route, Routes } from "react-router-dom";

import ShellLayout from "./components/ShellLayout";
import AchievementsPage from "./pages/AchievementsPage";
import AuthPage from "./pages/AuthPage";
import ArticleDetailPage from "./pages/ArticleDetailPage";
import ArticleEditorPage from "./pages/ArticleEditorPage";
import ArticlesPage from "./pages/ArticlesPage";
import DashboardPage from "./pages/DashboardPage";
import LeaderboardPage from "./pages/LeaderboardPage";
import ProfilePage from "./pages/ProfilePage";
import PurchasesPage from "./pages/PurchasesPage";
import ShopPage from "./pages/ShopPage";
import TasksPage from "./pages/TasksPage";

function App() {
  return (
    <Routes>
      <Route path="/auth" element={<AuthPage />} />
      <Route element={<ShellLayout />}>
        <Route path="/" element={<Navigate to="/profile" replace />} />
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/tasks" element={<TasksPage />} />
        <Route path="/achievements" element={<AchievementsPage />} />
        <Route path="/rating" element={<LeaderboardPage />} />
        <Route path="/shop" element={<ShopPage />} />
        <Route path="/purchases" element={<PurchasesPage />} />
        <Route path="/articles" element={<ArticlesPage />} />
        <Route path="/articles/new" element={<ArticleEditorPage />} />
        <Route path="/articles/:articleId" element={<ArticleDetailPage />} />
      </Route>
    </Routes>
  );
}

export default App;
