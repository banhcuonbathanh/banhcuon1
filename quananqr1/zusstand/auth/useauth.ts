import { useAuthStore } from "./controller/auth-controller";

// Custom hook for easier use in components
export const useAuth = () => {
    const {
      user,
      accessToken,
      refreshToken,
      loading,
      error,
      login,
      logout,
      refreshAccessTokenAction, // Updated method name
      clearError
    } = useAuthStore();
  
    return {
      user,
      accessToken,
      refreshToken,
      isAuthenticated: !!user,
      loading,
      error,
      login,
      logout,
      refreshAccessToken: refreshAccessTokenAction, // Use renamed action
      clearError
    };
  };
  