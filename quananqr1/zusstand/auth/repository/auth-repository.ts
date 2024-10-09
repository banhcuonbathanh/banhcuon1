import { useApiStore } from "@/zusstand/api/api-controller";
import axios, { AxiosInstance } from "axios";
import Cookies from "js-cookie";
import { IAuthRepository } from "./interface_auth_repository";
import envConfig from "@/config";
import {
  LoginBodyType,
  LoginResType,
  RegisterBodyType
} from "../domain/auth.schema";
import {
  GuestLoginBodyType,
  GuestLoginResType
} from "@/schemaValidations/guest.schema";

class AuthRepository implements IAuthRepository {
  private baseUrl = envConfig.NEXT_PUBLIC_API_ENDPOINT;
  private createUserEndpoint = envConfig.NEXT_PUBLIC_API_Create_User;
  private loginEndpoint = envConfig.NEXT_PUBLIC_API_Login;
  private logoutEndpoint = envConfig.NEXT_PUBLIC_API_Logout;
  private Guest_login = envConfig.NEXT_PUBLIC_Add_Guest_login;
  private http: AxiosInstance | null = null;

  private initialize(): AxiosInstance {
    const { http } = useApiStore.getState();
    this.http = http;
    return http;
  }

  private getHttp(): AxiosInstance {
    if (!this.http) {
      return this.initialize();
    }
    return this.http;
  }
  // http://localhost:8888/users
  async register(body: RegisterBodyType): Promise<RegisterBodyType> {
    const http = this.getHttp();
    console.log(
      "quananqr1/zusstand/auth/repository/auth-repository.ts body",
      body,
      `${this.baseUrl}${this.createUserEndpoint}`,
      "http",
      http
    );
    try {
      console.log(
        "quananqr1/zusstand/auth/repository/auth-repository.ts body inside try"
      );

      const response = await http.post(`http://localhost:8888/users`, body);
      console.log("User added successfully:", response.data);
      return response.data;
    } catch (error) {
      console.error("Error adding user:", error);
      throw error;
    }
  }

  async login(body: LoginBodyType): Promise<LoginResType> {
    const http = this.getHttp();
    const loginLink = `${this.baseUrl}${this.loginEndpoint}`;

    try {
      const response = await http.post<LoginResType>(loginLink, body);
      console.log(
        "quananqr1/zusstand/auth/repository/auth-repository.ts response",
        response.data.access_token
      );
      const store = useApiStore.getState();
      store.setAccessToken(response.data.access_token);

      // Store refresh token in a secure HTTP-only cookie
      Cookies.set("refreshToken", response.data.access_token, {
        secure: true,
        sameSite: "strict"
      });

      console.log(
        "Login successful, access token set:",
        response.data.access_token
      );

      return response.data;
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  }

  async logout(): Promise<void> {
    const link = this.baseUrl + this.logoutEndpoint;
    const store = useApiStore.getState();

    try {
      console.log("Attempting to logout. Sending request to:", link);

      // Check if there's an access token
      if (!store.accessToken) {
        console.log("No access token found. User might already be logged out.");
        return;
      }

      // Use the http instance from the store, which has the interceptors set up
      await store.http.post(link);
      console.log("Logout request successful");

      // Clear the access token from the store
      store.setAccessToken(null);

      // Remove the refresh token cookie
      Cookies.remove("refreshToken");

      console.log("Logout process completed. Tokens cleared.");
    } catch (error) {
      console.error("Error during logout:", error);
      // Even if the server request fails, clear tokens on the client side
      store.setAccessToken(null);
      Cookies.remove("refreshToken");
      throw error; // Re-throw the error for the caller to handle if needed
    }
  }

  async guestLogin(body: GuestLoginBodyType): Promise<GuestLoginResType> {
    const http = this.getHttp();
    const guestLoginLink = `${this.baseUrl}${this.Guest_login}`;

    const response = await http.post<GuestLoginResType>(guestLoginLink, body);

    const { setAccessToken } = useApiStore.getState();
    setAccessToken(response.data.data.accessToken);

    // Store refresh token in a secure HTTP-only cookie
    Cookies.set("refreshToken", response.data.data.refreshToken, {
      secure: true,
      sameSite: "strict"
    });

    return response.data;
  }
}

export const authApi = new AuthRepository();
