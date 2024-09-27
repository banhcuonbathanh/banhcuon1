import axios, { AxiosInstance } from 'axios';
import Cookies from 'js-cookie';
import envConfig from "@/config";
import { useApiStore } from '@/zusstand/api/api-controller';
import { RegisterBodyType, LoginBodyType } from '../domain/auth.schema';
import { IAuthRepository, LoginResType, RefreshTokenResType } from './interface_auth_repository';

class AuthRepository implements IAuthRepository {
  private baseUrl = envConfig.NEXT_PUBLIC_API_ENDPOINT;
  private createUserEndpoint = envConfig.NEXT_PUBLIC_API_Create_User;
  private http: AxiosInstance;

  constructor() {
    const { http, setAccessToken } = useApiStore.getState();
    this.http = http;

    // Setup interceptors
    this.http.interceptors.request.use(
      (config) => {
        const { accessToken } = useApiStore.getState();
        if (accessToken) {
          config.headers = config.headers || {};
          config.headers.Authorization = `Bearer ${accessToken}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    this.http.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config;
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;
          try {
            const refreshResult = await this.refreshToken();
            setAccessToken(refreshResult.data.accessToken);
            originalRequest.headers.Authorization = `Bearer ${refreshResult.data.accessToken}`;
            return this.http(originalRequest);
          } catch (refreshError) {
            // Handle refresh token failure (e.g., logout user)
            this.logout();
            return Promise.reject(refreshError);
          }
        }
        return Promise.reject(error);
      }
    );
  }

  async register(body: RegisterBodyType): Promise<RegisterBodyType> {
    try {
      const response = await this.http.post(this.createUserEndpoint, body);
      console.log("User added successfully:", response.data);
      return response.data;
    } catch (error) {
      console.error("Error adding user:", error);
      throw error;
    }
  }

  async login(body: LoginBodyType): Promise<LoginResType> {
    const response = await this.http.post<LoginResType>("/auth/login", body);
    const { setAccessToken } = useApiStore.getState();
    setAccessToken(response.data.data.accessToken);
    // Store refresh token in a secure HTTP-only cookie
    Cookies.set('refreshToken', response.data.data.refreshToken, { secure: true, sameSite: 'strict' });
    return response.data;
  }

  async logout(): Promise<void> {
    await this.http.post("/auth/logout");
    const { setAccessToken } = useApiStore.getState();
    setAccessToken(null);
    Cookies.remove('refreshToken');
  }

  async refreshToken(): Promise<RefreshTokenResType> {
    const refreshToken = Cookies.get('refreshToken');
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }
    const response = await this.http.post<RefreshTokenResType>("/auth/refresh-token", { refreshToken });
    // Update the refresh token cookie if a new one is provided
    if (response.data.data.refreshToken) {
      Cookies.set('refreshToken', response.data.data.refreshToken, { secure: true, sameSite: 'strict' });
    }
    return response.data;
  }
}

export const authApi = new AuthRepository();