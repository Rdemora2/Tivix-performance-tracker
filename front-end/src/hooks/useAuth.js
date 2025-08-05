import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import useAppStore from '../store/useAppStore';
import { tokenUtils } from '../services/api';

export const useAuth = () => {
  const { user, isAuthenticated, logout, loadUserProfile } = useAppStore();
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuth = async () => {
      if (!tokenUtils.isValid()) {
        logout();
        navigate('/login', { replace: true });
        return;
      }

      if (!user && tokenUtils.isValid()) {
        try {
          await loadUserProfile();
        } catch (error) {
          console.error('Error loading user profile:', error);
          logout();
          navigate('/login', { replace: true });
        }
      }
    };

    checkAuth();
  }, [user, loadUserProfile, logout, navigate]);

  return {
    user,
    isAuthenticated: isAuthenticated && tokenUtils.isValid(),
    logout: () => {
      logout();
      navigate('/login', { replace: true });
    }
  };
};

export default useAuth;
