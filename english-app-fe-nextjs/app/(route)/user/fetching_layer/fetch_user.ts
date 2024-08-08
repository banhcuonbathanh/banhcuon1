// api/userApi.ts
import axios from "axios";
import { IUserApi } from "./interface_User_Api";
import { mockUsers } from "../mock/mockUserAPI";

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

class UserApi implements IUserApi {
  private baseUrl = "/api/users";

  constructor() {
    this.setupAxiosMockInterceptors();
  }

  private setupAxiosMockInterceptors() {
    axios.interceptors.response.use(async (response) => {
      // Simulate network delay
      await delay(5000);
      return response;
    });

    axios.interceptors.request.use(async (config) => {
      const { method, url, data } = config;

      if (url?.startsWith(this.baseUrl)) {
        let mockResponse;

        if (method === "get" && url === this.baseUrl) {
          mockResponse = { data: mockUsers };
        } else if (method === "get" && url.startsWith(`${this.baseUrl}/`)) {
          const id = url.split("/").pop();
          const user = mockUsers.find((u) => u.id === id);
          mockResponse = user
            ? { data: user }
            : { status: 404, data: "User not found" };
        } else if (method === "post") {
          const newUser = {
            id: String(mockUsers.length + 1),
            ...JSON.parse(data)
          };
          mockUsers.push(newUser);
          mockResponse = { data: newUser };
        } else if (method === "put") {
          const id = url?.split("/").pop();
          const index = mockUsers.findIndex((u) => u.id === id);
          if (index !== -1) {
            mockUsers[index] = { ...mockUsers[index], ...JSON.parse(data) };
            mockResponse = { data: mockUsers[index] };
          } else {
            mockResponse = { status: 404, data: "User not found" };
          }
        } else if (method === "delete") {
          const id = url?.split("/").pop();
          const index = mockUsers.findIndex((u) => u.id === id);
          if (index !== -1) {
            mockUsers.splice(index, 1);
            mockResponse = { status: 204 };
          } else {
            mockResponse = { status: 404, data: "User not found" };
          }
        }

        if (mockResponse) {
          // Instead of returning the mock response directly,
          // we'll throw it as an error to be caught by the response interceptor
          throw { response: mockResponse };
        }
      }

      return config;
    });

    // Add a response interceptor to handle our mock responses
    axios.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response) {
          return Promise.resolve(error.response);
        }
        return Promise.reject(error);
      }
    );
  }
  async fetchUsers(): Promise<any[]> {
    try {
      const response = await axios.get(this.baseUrl);
      return response.data;
    } catch (error) {
      throw new Error("Failed to fetch users");
    }
  }

  async fetchUserById(id: string): Promise<any> {
    try {
      const response = await axios.get(`${this.baseUrl}/${id}`);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to fetch user with id ${id}`);
    }
  }

  async createUser(userData: { name: string; email: string }): Promise<any> {
    try {
      const response = await axios.post(this.baseUrl, userData);
      return response.data;
    } catch (error) {
      throw new Error("Failed to create user");
    }
  }

  async updateUser(
    id: string,
    userData: { name?: string; email?: string }
  ): Promise<any> {
    try {
      const response = await axios.put(`${this.baseUrl}/${id}`, userData);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to update user with id ${id}`);
    }
  }

  async deleteUser(id: string): Promise<void> {
    try {
      await axios.delete(`${this.baseUrl}/${id}`);
    } catch (error) {
      throw new Error(`Failed to delete user with id ${id}`);
    }
  }
}

// Export an instance of the class
export const userApi = new UserApi();
