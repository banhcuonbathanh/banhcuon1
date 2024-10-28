import { create } from "zustand";
import { persist } from "zustand/middleware";
import Cookies from "js-cookie";

import envConfig from "@/config";

import { User } from "@/schemaValidations/user.schema";
import { useApiStore } from "../api/api-controller";
import {
  RegisterBodyType,
  LoginBodyType,
  LoginResType
} from "@/schemaValidations/auth.schema";
import {
  GuestInfo,
  GuestLoginRequest,
  GuestLoginResponse,
  LogoutRequest
} from "@/schemaValidations/interface/type_guest";
interface AuthState {
  user: User | null;
  guest: GuestInfo | null;
  accessToken: string | null;
  refreshToken: string | null;
  loading: boolean;
  error: string | null;
  isLoginDialogOpen: boolean;
  isGuest: boolean;
}

interface AuthActions {
  register: (body: RegisterBodyType) => Promise<void>;
  login: (body: LoginBodyType) => Promise<void>;
  logout: () => Promise<void>;
  refreshAccessToken: () => Promise<void>;
  guestLogin: (body: GuestLoginRequest) => Promise<void>;
  guestLogout: (body: LogoutRequest) => Promise<void>;
  clearError: () => void;
  openLoginDialog: () => void;
  closeLoginDialog: () => void;
}

type AuthStore = AuthState & AuthActions;

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      user: null,
      guest: null,
      accessToken: null,
      refreshToken: null,
      loading: false,
      error: null,
      isLoginDialogOpen: false,
      isGuest: false,

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
            loading: false,
            isGuest: true, // Set isGuest to true when user is not null
            guest: null
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
          // console.log(
          //   "quananqr1/zusstand/new_auth/new_auth_controller.ts login asdkfjhaskdjf",
          //   response.data.user
          // );
          set({
            user: {
              ...response.data.user,
              password: body.password
            },
            guest: null,
            isGuest: false, // Set isGuest to true when user is not null
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

      guestLogin: async (body: GuestLoginRequest) => {
        set({ loading: true, error: null });
        try {
          const guest_login_link =
            envConfig.NEXT_PUBLIC_API_ENDPOINT +
            envConfig.NEXT_PUBLIC_API_Guest_Login;
          const response = await useApiStore
            .getState()
            .http.post<GuestLoginResponse>(`${guest_login_link}`, body);
          set({
            user: null,
            guest: response.data.guest,
            isGuest: true, // Set isGuest to false when user is null and guest is not null
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
            error:
              error instanceof Error ? error.message : "Guest login failed",
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
            guest: null,
            isGuest: false,
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

      guestLogout: async (body: LogoutRequest) => {
        const guest_logout_link =
          envConfig.NEXT_PUBLIC_API_ENDPOINT +
          envConfig.NEXT_PUBLIC_API_Guest_Logout;
        set({ loading: true, error: null });
        try {
          await useApiStore.getState().http.post(`${guest_logout_link}`, body);
          set({
            user: null,
            guest: null,
            isGuest: false,
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
            error:
              error instanceof Error ? error.message : "Guest logout failed",
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
