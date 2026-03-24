import { Navigate, Route, Routes } from "react-router-dom";
import { AppLayout } from "./components/AppLayout";
import { ProtectedRoute } from "./routes/ProtectedRoute";
import { CategoriesPage } from "./pages/CategoriesPage";
import { ContactsPage } from "./pages/ContactsPage";
import { LedgersPage } from "./pages/LedgersPage";
import { LoginPage } from "./pages/LoginPage";
import { RecordsPage } from "./pages/RecordsPage";
import { RegisterPage } from "./pages/RegisterPage";
import { RemindersPage } from "./pages/RemindersPage";
import { StatsPage } from "./pages/StatsPage";

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route element={<ProtectedRoute />}>
        <Route element={<AppLayout />}>
          <Route path="/" element={<Navigate to="/records" replace />} />
          <Route path="/records" element={<RecordsPage />} />
          <Route path="/contacts" element={<ContactsPage />} />
          <Route path="/ledgers" element={<LedgersPage />} />
          <Route path="/categories" element={<CategoriesPage />} />
          <Route path="/stats" element={<StatsPage />} />
          <Route path="/reminders" element={<RemindersPage />} />
        </Route>
      </Route>
      <Route path="*" element={<Navigate to="/records" replace />} />
    </Routes>
  );
}
