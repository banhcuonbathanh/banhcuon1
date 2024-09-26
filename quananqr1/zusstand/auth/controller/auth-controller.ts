import { authApplication } from "../application/auth-application";
import { LoginBodyType, RegisterBodyType } from "../domain/auth.schema";
import { IAuthStore } from "./auth-controller-interface";
import { create } from "zustand";

export const useAuthStore = create<IAuthStore>((set, get) => ({
  // Initial state
  user: null,
  accessToken: null,
  refreshToken: null,
  loading: false,
  error: null,
  isLoginDialogOpen: false, // New state to control dialog visibility

  // Actions



   register : async (body: RegisterBodyType) => {
    set({ loading: true, error: null });
    try {
      const result = await authApplication.register(body);
      if (result.success && result.data) {
        set({
          user: {
            id: 0, // Assuming the API returns an id
            name: result.data.name,
            email: result.data.email,
            role: "Guest", // Default role for new registrations
            avatar: null, // No avatar set during registration
          },
          error: null,
          isLoginDialogOpen: false // Close the dialog on successful registration
        });
      } else {
        throw new Error(result.error || "Registration failed");
      }
    } catch (error) {
      set({ error: error instanceof Error ? error.message : "An unknown error occurred" });
    } finally {
      set({ loading: false });
    }
  },
  //------
  login: async (body: LoginBodyType) => {
    set({ loading: true, error: null });
    try {
      const result = await authApplication.login(body);
      if (result.success && result.data) {
        set({
          user: result.data.account,
          accessToken: result.data.accessToken,
          refreshToken: result.data.refreshToken,
          error: null,
          isLoginDialogOpen: false // Close the dialog on successful login
        });
      } else {
        throw new Error(result.error || "Login failed");
      }
    } catch (error) {
      set({ error: error instanceof Error ? error.message : "An unknown error occurred" });
    } finally {
      set({ loading: false });
    }
  },

  logout: async () => {
    set({ loading: true, error: null });
    try {
      const result = await authApplication.logout();
      if (result.success) {
        set({ user: null, accessToken: null, refreshToken: null, error: null });
      } else {
        throw new Error(result.error || "Logout failed");
      }
    } catch (error) {
      set({ error: error instanceof Error ? error.message : "An unknown error occurred" });
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
      set({ error: error instanceof Error ? error.message : "An unknown error occurred" });
    } finally {
      set({ loading: false });
    }
  },

  clearError: () => set({ error: null }),

  // New actions to control dialog visibility
  openLoginDialog: () => set({ isLoginDialogOpen: true }),
  closeLoginDialog: () => set({ isLoginDialogOpen: false })
}));