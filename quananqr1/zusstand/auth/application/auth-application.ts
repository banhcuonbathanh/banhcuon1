import {
  LoginBodyType,
  LoginResType,
  LogoutBodyType,
  RefreshTokenBodyType,
  RefreshTokenResType,
  RegisterBodyType
} from "../domain/auth.schema";
import { AxiosError } from "axios";
import { authApi } from "../repository/auth-repository";
import { IAuthRepository } from "../repository/interface_auth_repository";
import { IAuthApplication } from "./interface-auth-application";

// Updated AuthApplication class implementing IAuthApplication
class AuthApplication implements IAuthApplication {
  private authRepository: IAuthRepository;
  // private authRepository: IAuthRepository;

  // constructor() {
  //   this.authRepository = new AuthRepository();
  //   this.authRepository.initialize();
  // }

  constructor(authRepository: IAuthRepository) {
    this.authRepository = authRepository;

  }
  sLogout(body: LogoutBodyType & { accessToken: string; }): Promise<{ success: boolean; error?: string; }> {
    throw new Error("Method not implemented.");
  }
  sRefreshToken(body: RefreshTokenBodyType): Promise<{ success: boolean; data?: RefreshTokenResType["data"]; error?: string; }> {
    throw new Error("Method not implemented.");
  }

  async sLogin(body: LoginBodyType): Promise<{
    success: boolean;
    data?: LoginResType["data"];
    error?: string;
  }> {
    try {
      const response = await this.authRepository.login(body);
      return {
        success: true,
        data: response.data
      };
    } catch (error) {
      return this.handleError(error, "Server login failed");
    }
  }

  async login(body: LoginBodyType): Promise<{
    success: boolean;
    data?: LoginResType["data"];
    error?: string;
  }> {
    try {
      const response = await this.authRepository.login(body);
      return {
        success: true,
        data: response.data
      };
    } catch (error) {
      return this.handleError(error, "Login failed");
    }
  }

  // async sLogout(
  //   body: LogoutBodyType & { accessToken: string }
  // ): Promise<{ success: boolean; error?: string }> {
  //   try {
  //     await this.authRepository.logout(body);
  //     return { success: true };
  //   } catch (error) {
  //     return this.handleError(error, "Server logout failed");
  //   }
  // }

  async logout(): Promise<{ success: boolean; error?: string }> {
    try {
      await this.authRepository.logout();
      return { success: true };
    } catch (error) {
      return this.handleError(error, "Logout failed");
    }
  }

  // async sRefreshToken(body: RefreshTokenBodyType): Promise<{
  //   success: boolean;
  //   data?: RefreshTokenResType["data"];
  //   error?: string;
  // }> {
  //   try {
  //     const response = await this.authRepository.sRefreshToken(body);
  //     return {
  //       success: true,
  //       data: response.data
  //     };
  //   } catch (error) {
  //     return this.handleError(error, "Server token refresh failed");
  //   }
  // }

  async refreshToken(): Promise<{
    success: boolean;
    data?: RefreshTokenResType["data"];
    error?: string;
  }> {
    try {
      // const response = await this.authRepository.refreshToken();
      return {
        success: true,
        // data: response.payload.data
      };
    } catch (error) {
      return this.handleError(error, "Token refresh failed");
    }
  }

  private handleError(
    error: unknown,
    defaultMessage: string
  ): { success: false; error: string } {
    if (error instanceof AxiosError) {
      const errorMessage = error.response?.data?.message || error.message;
      return {
        success: false,
        error: errorMessage
      };
    }
    return {
      success: false,
      error: defaultMessage
    };
  }

  // Additional helper methods can be added here if needed

  async register(body: RegisterBodyType): Promise<{
    success: boolean;
    data?: RegisterBodyType;
    error?: string;
  }> {

    console.log("register AuthApplication quananqr1/zusstand/auth/application/auth-application.ts")
    try {
      const response = await this.authRepository.register(body);
      return {
        success: true,
        data: response
      };
    } catch (error) {
      return this.handleError(error, "register failed");
    }
  }
}

// Export an instance of the application layer
export const authApplication = new AuthApplication(authApi);
