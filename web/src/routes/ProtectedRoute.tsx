import { Navigate, Outlet, useLocation } from "react-router-dom";
import { getAccessToken } from "../api/client";

export function ProtectedRoute() {
  const loc = useLocation();
  if (!getAccessToken()) {
    return <Navigate to="/login" replace state={{ from: loc.pathname }} />;
  }
  return <Outlet />;
}
