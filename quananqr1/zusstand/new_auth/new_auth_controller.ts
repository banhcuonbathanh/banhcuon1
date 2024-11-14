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
import { generateFormattedName } from "@/lib/utils";
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
}

interface AuthActions {
  register: (body: RegisterBodyType) => Promise<void>;
  login: (body: LoginBodyType) => Promise<void>;
  logout: () => Promise<void>;
  refreshAccessToken: () => Promise<void>;
  guestLogin: (body: GuestLoginBodyType) => Promise<void>;
  guestLogout: (body: LogoutRequest) => Promise<void>;
  clearError: () => void;
  openLoginDialog: () => void;
  closeLoginDialog: () => void;
  openGuestDialog: () => void;
  closeGuestDialog: () => void;
  openRegisterDialog: () => void;
  closeRegisterDialog: () => void;
}

type AuthStore = AuthState & AuthActions;

export const useAuthStore = create<AuthStore>()(
  persist(
    (set) => ({
      user: null,
      guest: null,
      accessToken: null,
      refreshToken: null,
      loading: false,
      error: null,
      isLoginDialogOpen: false,
      isGuestDialogOpen: false, // New initial state
      isRegisterDialogOpen: false, // New initial state
      isGuest: false,
      userId: null,

      register: async (body: RegisterBodyType) => {
        set({ loading: true, error: null });

        // Generate the formatted name before making the API call
        const formattedName = generateFormattedName(body.name);

        try {
          const response = await useApiStore
            .getState()
            .http.post<User>(`${envConfig.NEXT_PUBLIC_API_Create_User}`, {
              ...body,
              name: formattedName // Use the generated formatted name here
            });

          set({
            user: response.data,
            userId: response.data.id.toString(), // Set userId from user.id
            error: null,
            isLoginDialogOpen: false,
            isRegisterDialogOpen: false, // Close register dialog on success
            loading: false,
            isGuest: true,
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
          set({
            user: {
              ...response.data.user,
              password: body.password
            },
            userId: response.data.user.id.toString(), // Set userId from user.id
            guest: null,
            isGuest: false,
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

      guestLogin: async (body: GuestLoginBodyType) => {
        console.log(
          "quananqr1/zusstand/new_auth/new_auth_controller.ts GuestLoginBodyType",
          body
        );
        set({ loading: true, error: null });
        try {
          // Set the table token before making the request
          useApiStore.getState().setTableToken(body.token);

          const guest_login_link =
            envConfig.NEXT_PUBLIC_API_ENDPOINT +
            envConfig.NEXT_PUBLIC_API_Guest_Login;

          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts guest_login_link",
            guest_login_link
          );

          const response = await useApiStore
            .getState()
            .http.post<GuestLoginResponse>(`${guest_login_link}`, {
              name: generateFormattedName(body.name), // Add generated name here
              table_number: body.tableNumber,
              token: body.token
            });

          console.log(
            "quananqr1/zusstand/new_auth/new_auth_controller.ts response.data",
            response.data
          );

          set({
            userId: response.data.guest.id.toString(), // Set userId from guest.id
            user: null,
            guest: response.data.guest,
            isGuest: true,
            accessToken: response.data.access_token,
            refreshToken: response.data.refresh_token,
            error: null,
            isLoginDialogOpen: false,
            isGuestDialogOpen: false,
            loading: false,
            isRegisterDialogOpen: false
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
            userId: null, // Reset userId
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
            userId: "", // Reset userId
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
      closeRegisterDialog: () => set({ isRegisterDialogOpen: false })
    }),
    {
      name: "auth-storage",
      skipHydration: true
    }
  )
);

// // Example usage
// const { openGuestDialog, closeGuestDialog, openRegisterDialog, closeRegisterDialog } = useAuthStore();

// // Open guest dialog
// openGuestDialog();

// // Close register dialog
// closeRegisterDialog();
