import { create } from "zustand";
import { persist } from "zustand/middleware";
import Cookies from "js-cookie";

import envConfig from "@/config";

import { User } from "@/schemaValidations/user.schema";
import { useApiStore } from "../api/api-controller";
import { RegisterBodyType, LoginBodyType, LoginResType } from "@/schemaValidations/auth.schema";

// AuthStore interface and implementation
interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  loading: boolean;
  error: string | null;
  isLoginDialogOpen: boolean;
}

interface AuthActions {
  register: (body: RegisterBodyType) => Promise<void>;
  login: (body: LoginBodyType) => Promise<void>;
  logout: () => Promise<void>;
  refreshAccessToken: () => Promise<void>;

  clearError: () => void;
  openLoginDialog: () => void;
  closeLoginDialog: () => void;
}

type AuthStore = AuthState & AuthActions;

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      loading: false,
      error: null,
      isLoginDialogOpen: false,

      register: async (body: RegisterBodyType) => {
        set({ loading: true, error: null });
        try {
          const response = await useApiStore
            .getState()
            .http.post<User>(`${envConfig.NEXT_PUBLIC_API_Create_User}`, body);
          set({
            user: response.data,
            error: null,
            isLoginDialogOpen: false,
            loading: false
          });
        } catch (error) {
          set({
            error:
              error instanceof Error ? error.message : "Registration failed",
            loading: false
          });
        }
      },
      login: async (body: LoginBodyType) => {
        set({ loading: true, error: null });
        try {
          const response = await useApiStore
            .getState()
            .http.post<LoginResType>(
              `${envConfig.NEXT_PUBLIC_API_Login}`,
              body
            );
          set({
            user: {
              ...response.data.user,
              password: body.password
            },
            accessToken: response.data.access_token,
            refreshToken: response.data.refresh_token,
            error: null,
            isLoginDialogOpen: false,
            loading: false
          });
          useApiStore.getState().setAccessToken(response.data.access_token);
          Cookies.set("accessToken", response.data.access_token, {
            secure: true,
            sameSite: "strict"
          });
          Cookies.set("refreshToken", response.data.refresh_token, {
            secure: true,
            sameSite: "strict"
          });
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : "Login failed",
            loading: false
          });
        }
      },

      logout: async () => {
        set({ loading: true, error: null });
        try {
          await useApiStore
            .getState()
            .http.post(`${envConfig.NEXT_PUBLIC_API_Logout}`);
          set({
            user: null,
            accessToken: null,
            refreshToken: null,
            error: null,
            loading: false
          });
          useApiStore.getState().setAccessToken(null);
          Cookies.remove("accessToken");
          Cookies.remove("refreshToken");
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : "Logout failed",
            loading: false
          });
        }
      },

      refreshAccessToken: async () => {
        set({ loading: true, error: null });
        try {
          await useApiStore.getState().refreshToken();
          const newAccessToken = useApiStore.getState().accessToken;
          set({
            accessToken: newAccessToken,
            error: null,
            loading: false
          });
        } catch (error) {
          set({
            error:
              error instanceof Error ? error.message : "Token refresh failed",
            loading: false
          });
        }
      },

      clearError: () => set({ error: null }),
      openLoginDialog: () => set({ isLoginDialogOpen: true }),
      closeLoginDialog: () => set({ isLoginDialogOpen: false })
    }),
    {
      name: "auth-storage",
      skipHydration: true
    }
  )
);

// Custom hooks for each operation
export const useRegisterMutation = () => {
  const { register, loading, error } = useAuthStore();
  return { mutateAsync: register, isPending: loading, error };
};

export const useLoginMutation = () => {
  const { login, loading, error } = useAuthStore();
  return { mutateAsync: login, isPending: loading, error };
};

export const useLogoutMutation = () => {
  const { logout, loading, error } = useAuthStore();
  return { mutateAsync: logout, isPending: loading, error };
};

export const useRefreshTokenMutation = () => {
  const { refreshAccessToken, loading, error } = useAuthStore();
  return { mutateAsync: refreshAccessToken, isPending: loading, error };
};
