import { useApiStore } from "@/zusstand/api/api-controller";
import axios, { AxiosInstance } from "axios";
import Cookies from "js-cookie";
import { IAuthRepository, LoginResType } from "./interface_auth_repository";
import envConfig from "@/config";
import { LoginBodyType, RegisterBodyType } from "../domain/auth.schema";

class AuthRepository implements IAuthRepository {
  private baseUrl = envConfig.NEXT_PUBLIC_API_ENDPOINT;
  private createUserEndpoint = envConfig.NEXT_PUBLIC_API_Create_User;
  private loginEndpoint = envConfig.NEXT_PUBLIC_API_Login;
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

      const responsetest = await axios.get(`http://localhost:8888/test`);

      console.log(" connection is ok ", responsetest);
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

    console.log(
      "quananqr1/zusstand/auth/repository/auth-repository.ts loginLink",
      loginLink
    );

    const response = await http.post<LoginResType>(loginLink, body);

    console.log(
      "login quananqr1/zusstand/auth/repository/auth-repository.ts",
      response.data
    );
    const { setAccessToken } = useApiStore.getState();
    setAccessToken(response.data.data.accessToken);
    // Store refresh token in a secure HTTP-only cookie
    Cookies.set("refreshToken", response.data.data.refreshToken, {
      secure: true,
      sameSite: "strict"
    });
    return response.data;
  }

  async logout(): Promise<void> {
    const http = this.getHttp();
    await http.post(`${this.baseUrl}/auth/logout`);
    const { setAccessToken } = useApiStore.getState();
    setAccessToken(null);
    Cookies.remove("refreshToken");
  }
}

export const authApi = new AuthRepository();
