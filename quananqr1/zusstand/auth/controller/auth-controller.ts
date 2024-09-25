import { authApplication } from "../application/auth-application";
import { LoginBodyType } from "../domain/auth.schema";
import { IAuthStore } from "./auth-controller-interface";
import { create } from "zustand";

export const useAuthStore = create<IAuthStore>((set, get) => ({
  // Initial state
  user: null,
  accessToken: null,
  refreshToken: null,  // State property for the refresh token
  loading: false,
  error: null,

  // Actions
  login: async (body: LoginBodyType) => {
    set({ loading: true, error: null });
    try {
      const result = await authApplication.login(body);
      if (result.success && result.data) {
        set({
          user: result.data.account,
          accessToken: result.data.accessToken,
          refreshToken: result.data.refreshToken,
          error: null
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

  // Renamed action to avoid conflict with state
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

  clearError: () => set({ error: null })
}));
