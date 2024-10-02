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
import { GuestLoginBodyType, GuestLoginResType } from "@/schemaValidations/guest.schema";

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
    const linklogout = this.baseUrl + this.logoutEndpoint;

    const response = await http.post<LoginResType>(loginLink, body);

    // console.log(
    //   "login quananqr1/zusstand/auth/repository/auth-repository.ts",
    //   response.data
    // );
    const { setAccessToken } = useApiStore.getState();
    setAccessToken(response.data.accessToken);
    // Store refresh token in a secure HTTP-only cookie
    Cookies.set("refreshToken", response.data.refreshToken, {
      secure: true,
      sameSite: "strict"
    });
    return response.data;
  }

  async logout(): Promise<void> {
    const link = this.baseUrl + this.logoutEndpoint;
    const http = this.getHttp();
    console.log(
      "login quananqr1/zusstand/auth/repository/auth-repository.ts logout link",
      link
    );


    await http.post(link);

    console.log(
      "login quananqr1/zusstand/auth/repository/auth-repository.ts logout before const { setAccessToken } = useApiStore.getState();"
    );
    const { setAccessToken } = useApiStore.getState();
    setAccessToken(null);

    console.log(
      "login quananqr1/zusstand/auth/repository/auth-repository.ts logout after const { setAccessToken } = useApiStore.getState();"
    );
    Cookies.remove("refreshToken");
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
