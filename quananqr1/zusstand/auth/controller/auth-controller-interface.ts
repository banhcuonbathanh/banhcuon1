import { LoginBodyType, LoginResType } from "../domain/auth.schema";

export interface IAuthState {
    user: LoginResType['data']['account'] | null;
    accessToken: string | null;
    refreshToken: string | null; // State property for refresh token
    loading: boolean;
    error: string | null;
  }
  
  export interface IAuthActions {
    login: (body: LoginBodyType) => Promise<void>;
    logout: () => Promise<void>;
    refreshAccessTokenAction: () => Promise<void>; // Renamed action for refreshing token
    clearError: () => void;
  }
  
  export interface IAuthStore extends IAuthState, IAuthActions {}
  