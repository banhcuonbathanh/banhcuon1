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
  isLogin: boolean; // New field for login status
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

  syncAuthState: () => void; // New method to sync auth state
  initializeAuthFromCookies: () => void; // New method to initialize auth from cookies
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
      isLogin: false, // Initialize isLogin as false

      register: async (body: RegisterBodyType) => {
        set({ loading: true, error: null });

        const formattedName = body.name;

        try {
          const response = await useApiStore
            .getState()
            .http.post<User>(`${envConfig.NEXT_PUBLIC_API_Create_User}`, {
              ...body,
              name: formattedName
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
            isLogin: true // Set isLogin to true after successful registration
          });
        } catch (error) {
          set({
            error:
              error instanceof Error ? error.message : "Registration failed",
            loading: false,
            isLogin: false // Ensure isLogin is false if registration fails
          });
        }
      },

      login: async (body: LoginBodyType, fromPath?: string | null) => {
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
            userId: response.data.user.id.toString(),
            guest: null,
            isGuest: false,
            accessToken: response.data.access_token,
            refreshToken: response.data.refresh_token,
            error: null,
            isLoginDialogOpen: false,
            loading: false,
            isLogin: true // Set isLogin to true after successful user login
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

          // new part to
          // console.log(
          //   "quananqr1/zusstand/new_auth/new_auth_controller.ts login fromPath",
          //   fromPath
          // );
          window.location.href = fromPath || "/";
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : "Login failed",
            loading: false,
            isLogin: false // Ensure isLogin is false if login fails
          });
        }
      },

      guestLogin: async (
        body: GuestLoginBodyType,
        fromPath?: string | null
      ) => {
        // console.log("quananqr1/zusstand/new_auth/new_auth_controller.ts ");
        set({ loading: true, error: null });
        try {
          useApiStore.getState().setTableToken(body.token);

          const guest_login_link =
            envConfig.NEXT_PUBLIC_API_ENDPOINT +
            envConfig.NEXT_PUBLIC_API_Guest_Login;
          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts guest_login_link NEXT_PUBLIC_API_ENDPOINT",
            guest_login_link,
            envConfig.NEXT_PUBLIC_API_ENDPOINT
          );
          const response = await useApiStore
            .getState()
            .http.post<GuestLoginResponse>(`${guest_login_link}`, {
              name: body.name,
              table_number: body.tableNumber,
              token: body.token
            });

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
            isLogin: true // Set isLogin to true after successful guest login
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

          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts guest login fromPath",
            fromPath
          );
          // Get the original destination path from search params
          // const searchParams = new URLSearchParams(window.location.search);
          // const fromPath = searchParams.get("from") || "/"; // default to home if no path
          // console.log(
          //   "quananqr1/zusstand/new_auth/new_auth_controller.ts guest login searchParams fromPath",
          //   searchParams,
          //   fromPath
          // );
          // Redirect to the original destination or home
          window.location.href = fromPath || "/";
        } catch (error) {
          set({
            error:
              error instanceof Error ? error.message : "Guest login failed",
            loading: false,
            isLogin: false // Ensure isLogin is false if guest login fails
          });
        }
      },

      logout: async () => {
        set({ loading: true, error: null });
        try {
          // First clear cookies before making the logout request
          Cookies.remove("accessToken", { path: "/" });
          Cookies.remove("refreshToken", { path: "/" });

          // Clear API store token
          useApiStore.getState().setAccessToken(null);

          // Make logout request
          await useApiStore
            .getState()
            .http.post(`${envConfig.NEXT_PUBLIC_API_Logout}`);

          // Clear all auth state
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
          // Even if the logout request fails, we should still clear local state
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
        const guest_logout_link =
          envConfig.NEXT_PUBLIC_API_ENDPOINT +
          envConfig.NEXT_PUBLIC_API_Guest_Logout;

        set({ loading: true, error: null });
        try {
          // First clear cookies before making the logout request
          Cookies.remove("accessToken", { path: "/" });
          Cookies.remove("refreshToken", { path: "/" });

          // Clear API store token
          useApiStore.getState().setAccessToken(null);

          // Make guest logout request
          await useApiStore.getState().http.post(`${guest_logout_link}`, body);

          // Clear all auth state
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
          // Even if the logout request fails, we should still clear local state
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
      openLoginDialog: () =>
        set({
          isLoginDialogOpen: true,
          isGuestDialogOpen: false,
          isRegisterDialogOpen: false
        }),
      closeLoginDialog: () => set({ isLoginDialogOpen: false }),

      openGuestDialog: () =>
        set({
          isGuestDialogOpen: true,
          isLoginDialogOpen: false,
          isRegisterDialogOpen: false
        }),
      closeGuestDialog: () => set({ isGuestDialogOpen: false }),

      openRegisterDialog: () =>
        set({
          isRegisterDialogOpen: true,
          isLoginDialogOpen: false,
          isGuestDialogOpen: false
        }),
      closeRegisterDialog: () => set({ isRegisterDialogOpen: false }),

      //

      syncAuthState: () => {
        // console.log(
        //   "quananqr1/zusstand/new_auth/new_auth_controller.ts syncAuthState"
        // );
        const accessToken = Cookies.get("accessToken");
        const refreshToken = Cookies.get("refreshToken");
        // console.log(
        //   "quananqr1/zusstand/new_auth/new_auth_controller.ts syncAuthState accessToken, refreshToken",
        //   accessToken,
        //   refreshToken
        // );
        if (accessToken && refreshToken) {
          try {
            const decoded = decodeToken(accessToken);
            // console.log(
            //   "quananqr1/zusstand/new_auth/new_auth_controller.ts syncAuthState decoded",
            //   decoded
            // );
            // Log exact state before setting
            const newState = {
              accessToken,
              refreshToken,
              isLogin: true,
              isGuest: decoded.role === "Guest",
              userId: decoded.id.toString()
            };
            // console.log(
            //   "quananqr1/zusstand/new_auth/new_auth_controller.ts New state being set:",
            //   newState
            // );

            set(newState);
            // console.log(
            //   "quananqr1/zusstand/new_auth/new_auth_controller.ts Current login state:",
            //   useAuthStore.getState()
            // );
          } catch (error) {
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

      // New method to initialize auth from cookies when app loads
      initializeAuthFromCookies: () => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts initializeAuthFromCookies"
        );
        const accessToken = Cookies.get("accessToken");
        const refreshToken = Cookies.get("refreshToken");

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
