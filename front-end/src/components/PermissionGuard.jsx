import { hasPermission } from "../services/api";

const PermissionGuard = ({
  action,
  resource,
  children,
  fallback = null,
  requireRole = null,
}) => {
  if (requireRole) {
    const payload = JSON.parse(
      atob(localStorage.getItem("token")?.split(".")[1] || "{}")
    );
    if (payload.role !== requireRole && payload.role !== "admin") {
      return fallback;
    }
  }

  if (!hasPermission(action, resource)) {
    return fallback;
  }

  return children;
};

export default PermissionGuard;
