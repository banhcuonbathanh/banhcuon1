import axios, { AxiosInstance } from "axios";
import { IAuthRepository } from "./interface_auth_repository";
import envConfig from "@/config";
import {
  RefreshTokenResType,
  LoginBodyType,
  LoginResType,
  LogoutBodyType,
  RefreshTokenBodyType,
  RegisterBodyType
} from "../domain/auth.schema";

class AuthRepository implements IAuthRepository {
  private baseUrl = envConfig.NEXT_PUBLIC_API_ENDPOINT;
  private createUserEndpoint = envConfig.NEXT_PUBLIC_API_Create_User;
  private http: AxiosInstance;
  private refreshTokenRequest: Promise<{
    status: number;
    payload: RefreshTokenResType;
  }> | null = null;

  constructor() {
    this.http = axios.create({
      baseURL: this.baseUrl
    });
  }

  async register(body: RegisterBodyType): Promise<RegisterBodyType> {
    console.log(
      "quananqr1/app/(public)/public-component/register-dialog.tsx hander use 2222"
    );
    const userData = {
      name: "Alice ",
      email: "alice.johnson@example.com22222",
      password: "password1231234",
      is_admin: false,
      phone: 1234567890,
      image: "alice.jpg",
      address: "123 Main St, Anytown, USA",
      created_at: "2024-08-19T16:17:16+07:00",
      updated_at: "2024-08-19T16:17:16+07:00"
    };

    try {
      const response = await axios.post(
        "http://localhost:8888/users",
        userData
      );
      console.log("User added successfully:", response.data);

      const mappedData: RegisterBodyType = {
        name: response.data.name,
        email: response.data.email,
        password: body.password, // Use the original password from the input body
        is_admin: response.data.is_admin,
        phone: String(response.data.phone), // Ensure phone is a string
        image: response.data.image,
        address: response.data.address,
        created_at: response.data.created_at,
        updated_at: response.data.updated_at
      };
      return mappedData;
    } catch (error) {
      console.error("Error adding user:", error);

      throw error; // Re-throw the error instead of returning undefined
    }
  }

  async sLogin(body: LoginBodyType): Promise<LoginResType> {
    const response = await this.http.post<LoginResType>("/auth/login", body);
    return response.data;
  }

  async login(body: LoginBodyType): Promise<LoginResType> {
    const response = await axios.post<LoginResType>("/api/auth/login", body, {
      baseURL: ""
    });
    return response.data;
  }

  async sLogout(body: LogoutBodyType & { accessToken: string }): Promise<void> {
    await this.http.post(
      "/auth/logout",
      { refreshToken: body.refreshToken },
      {
        headers: {
          Authorization: `Bearer ${body.accessToken}`
        }
      }
    );
  }

  async logout(): Promise<void> {
    await axios.post("/api/auth/logout", null, { baseURL: "" });
  }

  async sRefreshToken(
    body: RefreshTokenBodyType
  ): Promise<RefreshTokenResType> {
    const response = await this.http.post<RefreshTokenResType>(
      "/auth/refresh-token",
      body
    );
    return response.data;
  }

  async refreshToken(): Promise<{
    status: number;
    payload: RefreshTokenResType;
  }> {
    if (this.refreshTokenRequest) {
      return this.refreshTokenRequest;
    }

    this.refreshTokenRequest = axios
      .post<RefreshTokenResType>("/api/auth/refresh-token", null, {
        baseURL: ""
      })
      .then((response) => ({
        status: response.status,
        payload: response.data
      }));

    const result = await this.refreshTokenRequest;
    this.refreshTokenRequest = null;
    return result;
  }
}

// Export an instance of the class

async function register(body: RegisterBodyType): Promise<RegisterBodyType> {
  console.log(
    "quananqr1/app/(public)/public-component/register-dialog.tsx handler use"
  );

  const userData = {
    name: "Alice ",
    email: "alice.johnson@example.com1111111111111",
    password: "password1231234",
    is_admin: false,
    phone: "1234567890", // Changed to string to match RegisterBodyType
    image: "alice.jpg",
    address: "123 Main St, Anytown, USA",
    created_at: "2024-08-19T16:17:16+07:00",
    updated_at: "2024-08-19T16:17:16+07:00"
  };

  try {
    const response = await axios.post("http://localhost:8888/users", userData);
    console.log("User added successfully:", response.data);

    // Map response.data to RegisterBodyType
    const mappedData: RegisterBodyType = {
      name: response.data.name,
      email: response.data.email,
      password: body.password, // Use the original password from the input body
      is_admin: response.data.is_admin,
      phone: String(response.data.phone), // Ensure phone is a string
      image: response.data.image,
      address: response.data.address,
      created_at: response.data.created_at,
      updated_at: response.data.updated_at
    };

    return mappedData;
  } catch (error) {
    console.error("Error adding user:", error);
    throw error; // Re-throw the error to be handled by the caller
  }
}
export const authApi = new AuthRepository();
