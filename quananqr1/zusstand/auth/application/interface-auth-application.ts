import { LoginBodyType, LoginResType, LogoutBodyType, RefreshTokenBodyType, RefreshTokenResType, RegisterBodyType } from "../domain/auth.schema";


export interface IAuthApplication {
    sLogin(body: LoginBodyType): Promise<{ success: boolean; data?: LoginResType['data']; error?: string }>;
    login(body: LoginBodyType): Promise<{ success: boolean; data?: LoginResType['data']; error?: string }>;
    sLogout(body: LogoutBodyType & { accessToken: string }): Promise<{ success: boolean; error?: string }>;
    logout(): Promise<{ success: boolean; error?: string }>;
    sRefreshToken(body: RefreshTokenBodyType): Promise<{ success: boolean; data?: RefreshTokenResType['data']; error?: string }>;
    refreshToken(): Promise<{ success: boolean; data?: RefreshTokenResType['data']; error?: string }>;

    register(body: RegisterBodyType): Promise<{ success: boolean; data?: RegisterBodyType['data']; error?: string }>;
  }
  