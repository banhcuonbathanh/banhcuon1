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

interface AuthState {
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

      register: async (body: RegisterBodyType) => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts register function called with body:",
          body
        );
        set({ loading: true, error: null });

        const formattedName = body.name;

        try {
          const response = await useApiStore
            .getState()
            .http.post<User>(`${envConfig.NEXT_PUBLIC_API_Create_User}`, {
              ...body,
              name: formattedName
            });

          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts register successful, response:",
            response.data
          );

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
          console.error(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts register error:",
            error
          );
          set({
            error:
              error instanceof Error ? error.message : "Registration failed",
            loading: false,
            isLogin: false
          });
        }
      },

      login: async (body: LoginBodyType, fromPath?: string | null) => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts login function called with body:",
          body,
          "fromPath:",
          fromPath
        );
        set({ loading: true, error: null });
        try {
          const response = await useApiStore
            .getState()
            .http.post<LoginResType>(
              `${envConfig.NEXT_PUBLIC_API_Login}`,
              body
            );

          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts login successful, response:",
            response.data
          );

          set({
            user: {
              ...response.data.user,
              password: body.password
            },
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

          // window.location.href = fromPath || "/";
        } catch (error) {
          console.error(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts login error:",
            error
          );
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
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts guestLogin function called with body:",
          body,
          "fromPath:",
          fromPath
        );
        set({ loading: true, error: null });
        try {
          useApiStore.getState().setTableToken(body.token);

          const guest_login_link =
            envConfig.NEXT_PUBLIC_API_ENDPOINT +
            envConfig.NEXT_PUBLIC_API_Guest_Login;
          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts guest_login_link:",
            guest_login_link
          );
          const response = await useApiStore
            .getState()
            .http.post<GuestLoginResponse>(`${guest_login_link}`, {
              name: body.name,
              table_number: body.tableNumber,
              token: body.token
            });

          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts guestLogin successful, response:",
            response.data
          );

          set({
            userId: response.data.guest.id.toString(),
            user: null,
            guest: response.data.guest,
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
          console.error(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts guestLogin error:",
            error
          );
          set({
            error:
              error instanceof Error ? error.message : "Guest login failed",
            loading: false,
            isLogin: false
          });
        }
      },

      logout: async () => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts logout function called"
        );
        set({ loading: true, error: null });
        try {
          Cookies.remove("accessToken", { path: "/" });
          Cookies.remove("refreshToken", { path: "/" });

          useApiStore.getState().setAccessToken(null);

          await useApiStore
            .getState()
            .http.post(`${envConfig.NEXT_PUBLIC_API_Logout}`);

          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts logout successful"
          );

          set({
            userId: null,
            user: null,
            guest: null,
            isGuest: false,
            accessToken: null,
            refreshToken: null,
            error: null,
            loading: false,
            isLogin: false
          });
        } catch (error) {
          console.error(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts logout error:",
            error
          );
          Cookies.remove("accessToken", { path: "/" });
          Cookies.remove("refreshToken", { path: "/" });
          useApiStore.getState().setAccessToken(null);

          set({
            userId: null,
            user: null,
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
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts guestLogout function called with body:",
          body
        );
        const guest_logout_link =
          envConfig.NEXT_PUBLIC_API_ENDPOINT +
          envConfig.NEXT_PUBLIC_API_Guest_Logout;

        set({ loading: true, error: null });
        try {
          Cookies.remove("accessToken", { path: "/" });
          Cookies.remove("refreshToken", { path: "/" });

          useApiStore.getState().setAccessToken(null);

          await useApiStore.getState().http.post(`${guest_logout_link}`, body);

          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts guestLogout successful"
          );

          set({
            userId: null,
            user: null,
            guest: null,
            isGuest: false,
            accessToken: null,
            refreshToken: null,
            error: null,
            loading: false,
            isLogin: false
          });
        } catch (error) {
          console.error(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts guestLogout error:",
            error
          );
          Cookies.remove("accessToken", { path: "/" });
          Cookies.remove("refreshToken", { path: "/" });
          useApiStore.getState().setAccessToken(null);

          set({
            userId: null,
            user: null,
            guest: null,
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
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts refreshAccessToken function called"
        );
        set({ loading: true, error: null });
        try {
          await useApiStore.getState().refreshToken();
          const newAccessToken = useApiStore.getState().accessToken;

          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts refreshAccessToken successful, new token:",
            newAccessToken
          );

          set({
            accessToken: newAccessToken,
            error: null,
            loading: false
          });
        } catch (error) {
          console.error(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts refreshAccessToken error:",
            error
          );
          set({
            error:
              error instanceof Error ? error.message : "Token refresh failed",
            loading: false
          });
        }
      },

      clearError: () => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts clearError function called"
        );
        set({ error: null });
      },

      openLoginDialog: () => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts openLoginDialog function called"
        );
        set({
          isLoginDialogOpen: true,
          isGuestDialogOpen: false,
          isRegisterDialogOpen: false
        });
      },

      closeLoginDialog: () => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts closeLoginDialog function called"
        );
        set({ isLoginDialogOpen: false });
      },

      openGuestDialog: () => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts openGuestDialog function called"
        );
        set({
          isGuestDialogOpen: true,
          isLoginDialogOpen: false,
          isRegisterDialogOpen: false
        });
      },

      closeGuestDialog: () => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts closeGuestDialog function called"
        );
        set({ isGuestDialogOpen: false });
      },

      openRegisterDialog: () => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts openRegisterDialog function called"
        );
        set({
          isRegisterDialogOpen: true,
          isLoginDialogOpen: false,
          isGuestDialogOpen: false
        });
      },

      closeRegisterDialog: () => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts closeRegisterDialog function called"
        );
        set({ isRegisterDialogOpen: false });
      },

      syncAuthState: () => {
        // console.log(
        //   "quananqr1/zusstand/new_auth/new_auth_controller.ts syncAuthState function called"
        // );
        const accessToken = Cookies.get("accessToken");
        const refreshToken = Cookies.get("refreshToken");

        // console.log(
        //   "quananqr1/zusstand/new_auth/new_auth_controller.ts syncAuthState accessToken, refreshToken:",
        //   accessToken,
        //   refreshToken
        // );

        if (accessToken && refreshToken) {
          try {
            const decoded = decodeToken(accessToken);

            // console.log(
            //   "quananqr1/zusstand/new_auth/new_auth_controller.ts syncAuthState decoded token:",
            //   decoded
            // );

            const newState = {
              accessToken,
              refreshToken,
              isLogin: true,

              isGuest: decoded.role === "Guest",
              userId: decoded.id.toString()
            };

            // console.log(
            //   "quananqr1/zusstand/new_auth/new_auth_controller.ts syncAuthState new state:",
            //   newState
            // );

            set(newState);
          } catch (error) {
            console.error(
              "quananqr1/zusstand/new_auth/new_auth_controller.ts syncAuthState error:",
              error
            );
            // If token is invalid, clear auth state
            set({
              userId: null,
              user: null,
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
          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts syncAuthState: No tokens found"
          );
          // No tokens found, clear auth state
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
        // console.log(
        //   "quananqr1/zusstand/new_auth/new_auth_controller.ts initializeAuthFromCookies function called"
        // );
        const accessToken = Cookies.get("accessToken");
        const refreshToken = Cookies.get("refreshToken");

        // console.log(
        //   "quananqr1/zusstand/new_auth/new_auth_controller.ts initializeAuthFromCookies accessToken, refreshToken:",
        //   accessToken,
        //   refreshToken
        // );

        if (accessToken && refreshToken) {
          get().syncAuthState();
        }
      }
    }),
    {
      name: "auth-storage",
      skipHydration: true
    }
  )
);
