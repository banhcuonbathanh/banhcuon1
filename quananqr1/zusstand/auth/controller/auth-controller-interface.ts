import { LoginBodyType, LoginResType, RegisterBodyType } from "../domain/auth.schema";

export interface IAuthState {
  user: LoginResType["data"]["account"] | null;
  accessToken: string | null;
  refreshToken: string | null; // State property for refresh token
  loading: boolean;
  error: string | null;
  isLoginDialogOpen: boolean;
}

export interface IAuthActions {
  login: (body: LoginBodyType) => Promise<void>;
  logout: () => Promise<void>;
  refreshAccessTokenAction: () => Promise<void>; // Renamed action for refreshing token
  clearError: () => void;
  openLoginDialog: () => void;
  closeLoginDialog: () => void;

  register: (body: RegisterBodyType) => Promise<void>;
}

export interface IAuthStore extends IAuthState, IAuthActions {}
