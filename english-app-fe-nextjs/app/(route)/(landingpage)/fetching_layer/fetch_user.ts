// api/userApi.ts
import axios from 'axios';
import { IUserApi } from './interface_User_Api';


class UserApi implements IUserApi {

  private baseUrl = '/api/users';

  async fetchUsers(): Promise<any[]> {
    try {
      const response = await axios.get(this.baseUrl);
      return response.data;
    } catch (error) {
      throw new Error('Failed to fetch users');
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
      throw new Error('Failed to create user');
    }
  }

  async updateUser(id: string, userData: { name?: string; email?: string }): Promise<any> {
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