import { create } from "zustand";
import { persist } from "zustand/middleware";
import { authApplication } from "../application/auth-application";
import { LoginBodyType, RegisterBodyType } from "../domain/auth.schema";
import { IAuthStore } from "./auth-controller-interface";
import { GuestLoginBodyType } from "@/schemaValidations/guest.schema";

export const useAuthStore = create<IAuthStore>()(
  persist(
    (set, get) => ({
      // Initial state
      user: null,
      accessToken: null,
      refreshToken: null,
      loading: false,
      error: null,
      isLoginDialogOpen: false,

      // Actions
      register: async (body: RegisterBodyType) => {
        console.log(
          "register useAuthStore quananqr1/zusstand/auth/controller/auth-controller.ts"
        );
        set({ loading: true, error: null });
        try {
          const result = await authApplication.register(body);
          if (result.success && result.data) {
            set({
              user: {
                id: 0,
                name: result.data.name,
                email: result.data.email,
                role: "Guest",
                avatar: null
              },
              error: null,
              isLoginDialogOpen: false
            });
          } else {
            throw new Error(result.error || "Registration failed");
          }
        } catch (error) {
          set({
            error:
              error instanceof Error
                ? error.message
                : "An unknown error occurred"
          });
        } finally {
          set({ loading: false });
        }
      },

      login: async (body: LoginBodyType) => {
        set({ loading: true, error: null });
        try {
          const result = await authApplication.login(body);

          // console.log(
          //   "quananqr1/zusstand/auth/controller/auth-controller.ts auth controller login result.data",
          //   result.data
          // );
          if (result.success && result.data) {
            set({
              user: result.data,
              accessToken: result.data.accessToken,
              refreshToken: result.data.refreshToken,
              error: null,
              isLoginDialogOpen: false
            });
          } else {
            throw new Error(result.error || "Login failed");
          }
        } catch (error) {
          set({
            error:
              error instanceof Error
                ? error.message
                : "An unknown error occurred"
          });
        } finally {
          set({ loading: false });
        }
      },

      logout: async () => {
        set({ loading: true, error: null });
        try {
          const result = await authApplication.logout();
          if (result.success) {
            set({
              user: null,
              accessToken: null,
              refreshToken: null,
              error: null
            });
          } else {
            throw new Error(result.error || "Logout failed");
          }
        } catch (error) {
          set({
            error:
              error instanceof Error
                ? error.message
                : "An unknown error occurred"
          });
        } finally {
          set({ loading: false });
        }
      },

      refreshAccessTokenAction: async () => {
        set({ loading: true, error: null });
        try {
          const result = await authApplication.refreshToken();
          if (result.success && result.data) {
            set({
              accessToken: result.data.accessToken,
              refreshToken: result.data.refreshToken,
              error: null
            });
          } else {
            throw new Error(result.error || "Token refresh failed");
          }
        } catch (error) {
          set({
            error:
              error instanceof Error
                ? error.message
                : "An unknown error occurred"
          });
        } finally {
          set({ loading: false });
        }
      },
      guestLogin: async (body: GuestLoginBodyType) => {
        set({ loading: true, error: null });
        try {
          const result = await authApplication.guestLogin(body);
          if (result.success && result.data) {
            set({
              user: {
                id: result.data.guest.id,
                name: result.data.guest.name,
                role: result.data.guest.role,
                tableNumber: result.data.guest.tableNumber,
                createdAt: result.data.guest.createdAt,
                updatedAt: result.data.guest.updatedAt
              },
              accessToken: result.data.accessToken,
              refreshToken: result.data.refreshToken,
              error: null,
              isLoginDialogOpen: false
            });
          } else {
            throw new Error(result.error || "Guest login failed");
          }
        } catch (error) {
          set({
            error:
              error instanceof Error
                ? error.message
                : "An unknown error occurred"
          });
        } finally {
          set({ loading: false });
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
