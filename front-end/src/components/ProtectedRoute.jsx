import { Navigate } from "react-router-dom";
import { tokenUtils } from "../services/api";

const ProtectedRoute = ({ children, requiredRole }) => {
  const token = tokenUtils.get();

  if (!token || !tokenUtils.isValid()) {
    return <Navigate to="/login" replace />;
  }

  try {
    const payload = tokenUtils.getPayload();

    if (!payload) {
      tokenUtils.remove();
      return <Navigate to="/login" replace />;
    }

    const { role } = payload;

    if (requiredRole) {
      if (role !== requiredRole && role !== "admin") {
        return <Navigate to="/unauthorized" replace />;
      }
    }

    return children;
  } catch (error) {
    console.error("Error validating token:", error);
    tokenUtils.remove();
    return <Navigate to="/login" replace />;
  }
};

export default ProtectedRoute;
