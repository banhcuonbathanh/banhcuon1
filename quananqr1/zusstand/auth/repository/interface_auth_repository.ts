import { LoginBodyType, LoginResType, LogoutBodyType, RefreshTokenBodyType, RefreshTokenResType } from "../domain/auth.schema";


export interface IAuthRepository {
  sLogin(body: LoginBodyType): Promise<LoginResType>;
  login(body: LoginBodyType): Promise<LoginResType>;
  sLogout(body: LogoutBodyType & { accessToken: string }): Promise<void>;
  logout(): Promise<void>;
  sRefreshToken(body: RefreshTokenBodyType): Promise<RefreshTokenResType>;
  refreshToken(): Promise<{ status: number; payload: RefreshTokenResType }>;
}