import axios, { AxiosInstance } from "axios";
import { IAuthRepository } from "./interface_auth_repository";
import envConfig from "@/config";
import {
  LoginBodyType,
  LoginResType,
  LogoutBodyType,
  RefreshTokenBodyType,
  RefreshTokenResType
} from '@/schemaValidations/auth.schema';

class AuthRepository implements IAuthRepository {
  private baseUrl = envConfig.NEXT_PUBLIC_API_ENDPOINT;
  private createUserEndpoint = envConfig.NEXT_PUBLIC_API_Create_User;
  private http: AxiosInstance;
  private refreshTokenRequest: Promise<{ status: number; payload: RefreshTokenResType }> | null = null;

  constructor() {
    this.http = axios.create({
      baseURL: this.baseUrl,
    });
  }

  async sLogin(body: LoginBodyType): Promise<LoginResType> {
    const response = await this.http.post<LoginResType>('/auth/login', body);
    return response.data;
  }

  async login(body: LoginBodyType): Promise<LoginResType> {
    const response = await axios.post<LoginResType>('/api/auth/login', body, {
      baseURL: '',
    });
    return response.data;
  }

  async sLogout(body: LogoutBodyType & { accessToken: string }): Promise<void> {
    await this.http.post(
      '/auth/logout',
      { refreshToken: body.refreshToken },
      {
        headers: {
          Authorization: `Bearer ${body.accessToken}`
        }
      }
    );
  }

  async logout(): Promise<void> {
    await axios.post('/api/auth/logout', null, { baseURL: '' });
  }

  async sRefreshToken(body: RefreshTokenBodyType): Promise<RefreshTokenResType> {
    const response = await this.http.post<RefreshTokenResType>('/auth/refresh-token', body);
    return response.data;
  }

  async refreshToken(): Promise<{ status: number; payload: RefreshTokenResType }> {
    if (this.refreshTokenRequest) {
      return this.refreshTokenRequest;
    }
    
    this.refreshTokenRequest = axios.post<RefreshTokenResType>(
      '/api/auth/refresh-token',
      null,
      { baseURL: '' }
    ).then(response => ({
      status: response.status,
      payload: response.data
    }));

    const result = await this.refreshTokenRequest;
    this.refreshTokenRequest = null;
    return result;
  }
}

// Export an instance of the class
export const authApi = new AuthRepository();
