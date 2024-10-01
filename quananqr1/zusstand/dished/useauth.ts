import { useAuthStore } from "./controller/dished-controller";


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
    refreshAccessTokenAction,
    clearError,
    isLoginDialogOpen,
    openLoginDialog,
    closeLoginDialog
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
    refreshAccessToken: refreshAccessTokenAction,
    clearError,
    isLoginDialogOpen,
    openLoginDialog,
    closeLoginDialog
  };
};