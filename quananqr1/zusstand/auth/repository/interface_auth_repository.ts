import { GuestLoginBodyType, GuestLoginResType } from "@/schemaValidations/guest.schema";
import { LoginBodyType, LoginResType, LogoutBodyType, RefreshTokenBodyType, RegisterBodyType } from "../domain/auth.schema";



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

  guestLogin(body: GuestLoginBodyType): Promise<GuestLoginResType>;
}