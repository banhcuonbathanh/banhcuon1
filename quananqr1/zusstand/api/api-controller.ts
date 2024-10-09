import { create } from "zustand";
import { persist } from "zustand/middleware";
import axios, {
  AxiosInstance,
  InternalAxiosRequestConfig,
  AxiosResponse
} from "axios";
import Cookies from "js-cookie";
import envConfig from "@/config";
import { RefreshTokenResType } from "../auth/domain/auth.schema";

interface ApiStore {
  http: AxiosInstance;
  accessToken: string | null;
  setAccessToken: (token: string | null) => void;
  refreshToken: () => Promise<void>;
}

export const useApiStore = create<ApiStore>()(
  persist(
    (set, get) => ({
      http: axios.create({
        baseURL: envConfig.NEXT_PUBLIC_API_ENDPOINT
      }),
      accessToken: null,
      setAccessToken: (token) => {
        set({ accessToken: token });
        console.log("Access token updated:", token);
      },
      refreshToken: async () => {
        try {
          const refreshToken = Cookies.get("refreshToken");
          if (!refreshToken) {
            throw new Error("No refresh token available");
          }
          const response = await axios.post<RefreshTokenResType>(
            "/api/auth/refresh-token",
            { refreshToken },
            { baseURL: "" }
          );
          set({ accessToken: response.data.data.accessToken });
          // Update the refresh token cookie if a new one is provided
          if (response.data.data.refreshToken) {
            Cookies.set("refreshToken", response.data.data.refreshToken, {
              secure: true,
              sameSite: "strict"
            });
          }
        } catch (error) {
          // Handle refresh token error (e.g., logout user)
          set({ accessToken: null });
          Cookies.remove("refreshToken");
        }
      }
    }),
    {
      name: "api-storage",
      skipHydration: true
    }
  )
);

// Setup interceptors
let http: AxiosInstance;

if (typeof window !== "undefined") {
  const store = useApiStore.getState();
  http = store.http;

  http.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      const { accessToken } = useApiStore.getState();
      if (accessToken) {
        config.headers = config.headers || {};
        config.headers.Authorization = `Bearer ${accessToken}`;
        console.log("Adding access token to request:", accessToken);
      } else {
        console.log("No access token available for request");
      }
      return config;
    },
    (error) => Promise.reject(error)
  );

  http.interceptors.response.use(
    (response: AxiosResponse) => response,
    async (error) => {
      const originalRequest = error.config;
      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;
        await useApiStore.getState().refreshToken();
        return http(originalRequest);
      }
      return Promise.reject(error);
    }
  );
}

// export const useApiStore = create<ApiStore>()(
//   persist(
//     (set, get) => ({
//       http: axios.create({
//         baseURL: envConfig.NEXT_PUBLIC_API_ENDPOINT
//       }),
//       accessToken: null,
//       setAccessToken: (token) => {
//         set({ accessToken: token });
//         console.log("Access token updated:", token);
//       },
//       // ... other methods ...
//     }),
//     {
//       name: "api-storage",
//       skipHydration: true
//     }
//   )
// );
