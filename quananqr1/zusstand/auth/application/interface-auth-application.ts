import { GuestLoginBodyType, GuestLoginResType } from "@/schemaValidations/guest.schema";
import { LoginBodyType, LoginResType, LogoutBodyType, RefreshTokenBodyType, RefreshTokenResType, RegisterBodyType } from "../domain/auth.schema";


export interface IAuthApplication {
    sLogin(body: LoginBodyType): Promise<{ success: boolean; data?: LoginResType; error?: string }>;
    login(body: LoginBodyType): Promise<{ success: boolean; data?: LoginResType; error?: string }>;
    sLogout(body: LogoutBodyType & { accessToken: string }): Promise<{ success: boolean; error?: string }>;
    logout(): Promise<{ success: boolean; error?: string }>;
    sRefreshToken(body: RefreshTokenBodyType): Promise<{ success: boolean; data?: RefreshTokenResType['data']; error?: string }>;
    refreshToken(): Promise<{ success: boolean; data?: RefreshTokenResType['data']; error?: string }>;

    register(body: RegisterBodyType): Promise<{ success: boolean; data?: RegisterBodyType; error?: string }>;

    guestLogin(body: GuestLoginBodyType): Promise<{ success: boolean; data?: GuestLoginResType; error?: string }>;
  }
  