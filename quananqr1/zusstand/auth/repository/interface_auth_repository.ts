import { LoginBodyType, LogoutBodyType, RefreshTokenBodyType, RegisterBodyType } from "../domain/auth.schema";

export interface LoginResType {
  message: string;
  data: {
    accessToken: string;
    refreshToken: string;
    account: {
      name: string;
      email: string;
      id: number;
      role: "Owner" | "Employee" | "Guest";
      avatar: string | null;
    };
  };
}

export interface RefreshTokenResType {
  message: string;
  data: {
    accessToken: string;
    refreshToken: string;
  };
}

export interface IAuthRepository {
  login(body: LoginBodyType): Promise<LoginResType>;
  logout(): Promise<void>;
  // refreshToken(): Promise<RefreshTokenResType>;
  register(body: RegisterBodyType): Promise<RegisterBodyType>;
}