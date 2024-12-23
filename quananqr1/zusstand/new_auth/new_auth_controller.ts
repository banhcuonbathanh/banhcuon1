// First add the logging utility
const logAuth = (logId: number, message: string, data?: any) => {
  const logPath = loggerPaths.find(
    (path) => path.path === "quananqr1/zusstand/new_auth/new_auth_controller.ts"
  );

  if (!logPath || !logPath.enabled || !logPath.enabledLogIds.includes(logId)) {
    return;
  }

  const logInfo = logPath.logDescriptions[logId];
  if (!logInfo) {
    return;
  }

  console.log(`[AUTH][${logInfo.location}] ${message}`, data || "");
};

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
  GuestLoginResponse,
  LogoutRequest
} from "@/schemaValidations/interface/type_guest";
import { GuestLoginBodyType } from "@/schemaValidations/guest.schema";
import { decodeToken } from "@/lib/utils";
import { loggerPaths } from "@/lib/logger/loggerConfig";

export interface AuthState {
  userId: string | null;
  user: User | null;
  guest: GuestInfo | null;
  accessToken: string | null;
  refreshToken: string | null;
  loading: boolean;
  error: string | null;
  isLoginDialogOpen: boolean;
  isGuestDialogOpen: boolean;
  isRegisterDialogOpen: boolean;
  isGuest: boolean;
  isLogin: boolean;
  persistedUser: User | null;
}

interface AuthActions {
  register: (body: RegisterBodyType) => Promise<void>;
  login: (body: LoginBodyType, fromPath?: string | null) => Promise<void>;
  logout: () => Promise<void>;
  refreshAccessToken: () => Promise<void>;
  guestLogin: (
    body: GuestLoginBodyType,
    fromPath?: string | null
  ) => Promise<void>;
  guestLogout: (body: LogoutRequest) => Promise<void>;
  clearError: () => void;
  openLoginDialog: () => void;
  closeLoginDialog: () => void;
  openGuestDialog: () => void;
  closeGuestDialog: () => void;
  openRegisterDialog: () => void;
  closeRegisterDialog: () => void;
  syncAuthState: () => void;
  initializeAuthFromCookies: () => void;
}

type AuthStore = AuthState & AuthActions;

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      // Initial state
      user: null,
      guest: null,
      accessToken: null,
      refreshToken: null,
      loading: false,
      error: null,
      isLoginDialogOpen: false,
      isGuestDialogOpen: false,
      isRegisterDialogOpen: false,
      isGuest: false,
      userId: null,
      isLogin: false,
      persistedUser: null,

      register: async (body: RegisterBodyType) => {
        logAuth(1, "Registration attempt initiated", { email: body.email });
        set({ loading: true, error: null });

        try {
          const response = await useApiStore
            .getState()
            .http.post<User>(`${envConfig.NEXT_PUBLIC_API_Create_User}`, {
              ...body,
              name: body.name
            });

          logAuth(2, "Registration successful", {
            userId: response.data.id,
            name: response.data.name
          });

          set({
            user: response.data,
            userId: response.data.id.toString(),
            error: null,
            isLoginDialogOpen: false,
            isRegisterDialogOpen: false,
            loading: false,
            isGuest: true,
            guest: null,
            isLogin: true
          });
        } catch (error) {
          logAuth(12, "Registration error", { error });
          set({
            error:
              error instanceof Error ? error.message : "Registration failed",
            loading: false,
            isLogin: false
          });
        }
      },

      login: async (body: LoginBodyType, fromPath?: string | null) => {
        logAuth(3, "Login attempt initiated", { email: body.email });
        set({ loading: true, error: null });

        try {
          const response = await useApiStore
            .getState()
            .http.post<LoginResType>(
              `${envConfig.NEXT_PUBLIC_API_Login}`,
              body
            );

          const userData = {
            ...response.data.user,
            password: body.password
          };

          logAuth(4, "Login successful", {
            userId: response.data.user.id,
            email: response.data.user.email
          });

          set({
            user: userData,
            persistedUser: userData,
            userId: response.data.user.id.toString(),
            guest: null,
            isGuest: false,
            accessToken: response.data.access_token,
            refreshToken: response.data.refresh_token,
            error: null,
            isLoginDialogOpen: false,
            loading: false,
            isLogin: true
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
          logAuth(12, "Login error", { error });
          set({
            error: error instanceof Error ? error.message : "Login failed",
            loading: false,
            isLogin: false
          });
        }
      },

      guestLogin: async (
        body: GuestLoginBodyType,
        fromPath?: string | null
      ) => {
        logAuth(5, "Guest login attempt initiated", {
          name: body.name,
          tableNumber: body.tableNumber
        });
        set({ loading: true, error: null });

        try {
          useApiStore.getState().setTableToken(body.token);
          const guest_login_link = `${envConfig.NEXT_PUBLIC_API_ENDPOINT}${envConfig.NEXT_PUBLIC_API_Guest_Login}`;

          const response = await useApiStore
            .getState()
            .http.post<GuestLoginResponse>(guest_login_link, {
              name: body.name,
              table_number: body.tableNumber,
              token: body.token
            });

          const guestData = response.data.guest;
          logAuth(6, "Guest login successful", {
            guestId: guestData.id,
            name: guestData.name,
            tableNumber: guestData.table_number
          });

          set({
            userId: guestData.id.toString(),
            user: null,
            guest: guestData,
            persistedUser: null,
            isGuest: true,
            accessToken: response.data.access_token,
            refreshToken: response.data.refresh_token,
            error: null,
            isLoginDialogOpen: false,
            isGuestDialogOpen: false,
            loading: false,
            isRegisterDialogOpen: false,
            isLogin: true
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

          window.location.href = fromPath || "/";
        } catch (error) {
          logAuth(12, "Guest login error", { error });
          set({
            error:
              error instanceof Error ? error.message : "Guest login failed",
            loading: false,
            isLogin: false
          });
        }
      },

      logout: async () => {
        logAuth(7, "User logout initiated");
        set({ loading: true, error: null });

        try {
          Cookies.remove("accessToken", { path: "/" });
          Cookies.remove("refreshToken", { path: "/" });
          useApiStore.getState().setAccessToken(null);

          await useApiStore
            .getState()
            .http.post(`${envConfig.NEXT_PUBLIC_API_Logout}`);

          logAuth(7, "User logout successful");

          set({
            userId: null,
            user: null,
            persistedUser: null,
            guest: null,
            isGuest: false,
            accessToken: null,
            refreshToken: null,
            error: null,
            loading: false,
            isLogin: false
          });
        } catch (error) {
          logAuth(12, "Logout error", { error });
          Cookies.remove("accessToken", { path: "/" });
          Cookies.remove("refreshToken", { path: "/" });
          useApiStore.getState().setAccessToken(null);

          set({
            userId: null,
            user: null,
            persistedUser: null,
            guest: null,
            isGuest: false,
            accessToken: null,
            refreshToken: null,
            error: error instanceof Error ? error.message : "Logout failed",
            loading: false,
            isLogin: false
          });
        }
      },

      guestLogout: async (body: LogoutRequest) => {
        logAuth(8, "Guest logout initiated");
        set({ loading: true, error: null });
        try {
          const guest_logout_link =
            envConfig.NEXT_PUBLIC_API_ENDPOINT +
            envConfig.NEXT_PUBLIC_API_Guest_Logout;

          Cookies.remove("accessToken", { path: "/" });
          Cookies.remove("refreshToken", { path: "/" });

          useApiStore.getState().setAccessToken(null);

          await useApiStore.getState().http.post(`${guest_logout_link}`, body);

          logAuth(8, "Guest logout successful");

          set({
            userId: null,
            user: null,
            guest: null,
            persistedUser: null,
            isGuest: false,
            accessToken: null,
            refreshToken: null,
            error: null,
            loading: false,
            isLogin: false
          });
        } catch (error) {
          logAuth(12, "Guest logout error", { error });
          Cookies.remove("accessToken", { path: "/" });
          Cookies.remove("refreshToken", { path: "/" });
          useApiStore.getState().setAccessToken(null);

          set({
            userId: null,
            user: null,
            guest: null,
            persistedUser: null,
            isGuest: false,
            accessToken: null,
            refreshToken: null,
            error:
              error instanceof Error ? error.message : "Guest logout failed",
            loading: false,
            isLogin: false
          });
        }
      },

      refreshAccessToken: async () => {
        logAuth(9, "Token refresh initiated");
        set({ loading: true, error: null });

        try {
          await useApiStore.getState().refreshToken();
          const newAccessToken = useApiStore.getState().accessToken;

          logAuth(9, "Token refresh successful", { newAccessToken });

          set({
            accessToken: newAccessToken,
            error: null,
            loading: false
          });
        } catch (error) {
          logAuth(12, "Token refresh error", { error });
          set({
            error:
              error instanceof Error ? error.message : "Token refresh failed",
            loading: false
          });
        }
      },

      clearError: () => {
        set({ error: null });
      },

      openLoginDialog: () => {
        set({
          isLoginDialogOpen: true,
          isGuestDialogOpen: false,
          isRegisterDialogOpen: false
        });
      },

      closeLoginDialog: () => {
        set({ isLoginDialogOpen: false });
      },

      openGuestDialog: () => {
        set({
          isGuestDialogOpen: true,
          isLoginDialogOpen: false,
          isRegisterDialogOpen: false
        });
      },

      closeGuestDialog: () => {
        set({ isGuestDialogOpen: false });
      },

      openRegisterDialog: () => {
        set({
          isRegisterDialogOpen: true,
          isLoginDialogOpen: false,
          isGuestDialogOpen: false
        });
      },

      closeRegisterDialog: () => {
        set({ isRegisterDialogOpen: false });
      },

      syncAuthState: () => {
        logAuth(10, "Auth state synchronization initiated");
        const accessToken = Cookies.get("accessToken");
        const refreshToken = Cookies.get("refreshToken");
        const currentState = get();

        if (accessToken && refreshToken) {
          try {
            const decoded = decodeToken(accessToken);

            logAuth(10, "Auth state synchronized successfully", {
              userId: decoded.id,
              role: decoded.role
            });

            set({
              accessToken,
              refreshToken,
              isLogin: true,
              isGuest: decoded.role === "Guest",
              userId: decoded.id.toString(),
              user: currentState.persistedUser
            });
          } catch (error) {
            logAuth(12, "Auth state sync error", { error });
            console.error("Token validation failed during sync:", error);
            set({
              userId: null,
              user: null,
              persistedUser: null,
              guest: null,
              accessToken: null,
              refreshToken: null,
              isLogin: false,
              isGuest: false
            });
            Cookies.remove("accessToken");
            Cookies.remove("refreshToken");
          }
        } else {
          logAuth(10, "No tokens found during auth sync");
          set({
            userId: null,
            user: null,
            guest: null,
            accessToken: null,
            refreshToken: null,
            isLogin: false,
            isGuest: false
          });
        }
      },

      initializeAuthFromCookies: () => {
        logAuth(11, "Cookie-based auth initialization started");
        const accessToken = Cookies.get("accessToken");
        const refreshToken = Cookies.get("refreshToken");

        if (accessToken && refreshToken) {
          logAuth(11, "Tokens found in cookies, syncing auth state");
          get().syncAuthState();
        } else {
          logAuth(11, "No tokens found in cookies during initialization");
        }
      }
    }),
    {
      name: "auth-storage"
    }
  )
);

export default useAuthStore;
