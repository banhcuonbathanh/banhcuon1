
import { useApiStore } from '@/zusstand/api/api-controller';
import { AxiosInstance } from 'axios';
import Cookies from 'js-cookie';
import { IAuthRepository, LoginResType } from './interface_auth_repository';
import envConfig from '@/config';
import { LoginBodyType, RegisterBodyType } from '../domain/auth.schema';


class AuthRepository implements IAuthRepository {
  private baseUrl = envConfig.NEXT_PUBLIC_API_ENDPOINT;
  private createUserEndpoint = envConfig.NEXT_PUBLIC_API_Create_User;
  private loginEndpoint = envConfig.NEXT_PUBLIC_API_Login;
  private http: AxiosInstance;

  constructor() {
    const { http } = useApiStore.getState();
    this.http = http;
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
    const loginLink = this.baseUrl + this.loginEndpoint;

    console.log(
      "quananqr1/zusstand/auth/repository/auth-repository.ts loginLink",
      loginLink
    );

    const response = await this.http.post<LoginResType>(loginLink, body);

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
    await this.http.post("/auth/logout");
    const { setAccessToken } = useApiStore.getState();
    setAccessToken(null);
    Cookies.remove("refreshToken");
  }
}

export const authApi = new AuthRepository();