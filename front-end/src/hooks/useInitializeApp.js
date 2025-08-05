import { useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import useAppStore from "../store/useAppStore";
import api, { tokenUtils } from "../services/api";

export const useInitializeApp = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const location = useLocation();

  const { initializeStore, isAuthenticated } = useAppStore();

  useEffect(() => {
    const initialize = async () => {
      try {
        const publicRoutes = ["/login", "/initial-setup", "/set-new-password"];
        if (publicRoutes.includes(location.pathname)) {
          setLoading(false);
          return;
        }

        const initResponse = await api.init.check();

        if (!initResponse.initialized || initResponse.userCount === 0) {
          navigate("/initial-setup", { replace: true });
          setLoading(false);
          return;
        }

        if (!tokenUtils.isValid()) {
          navigate("/login", { replace: true });
          setLoading(false);
          return;
        }

        await initializeStore();
        
        // Check if user needs to change password after loading profile
        const userProfile = useAppStore.getState().user;
        if (userProfile?.needsPasswordChange) {
          navigate('/set-new-password', { replace: true });
          setLoading(false);
          return;
        }
      } catch (error) {
        console.error("Initialization error:", error);
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    initialize();
  }, [initializeStore, navigate, location.pathname]);

  return { loading, error, isAuthenticated };
};

export default useInitializeApp;
