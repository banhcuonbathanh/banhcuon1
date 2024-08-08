// application/userService.ts
import { IUser } from "../domain/interface_User";
import { User } from "../domain/user";
import { userApi } from "../fetching_layer/fetch_user";
import { IUserService } from "./interface_Application";

export class UserService implements IUserService {
  async getUsers(): Promise<IUser[]> {
    try {
      const usersData = await userApi.fetchUsers();
      return usersData.map(User.fromJSON);
    } catch (error) {
      console.error("Error fetching users:", error);
      throw error;
    }
  }

  async getUserById(id: string): Promise<IUser> {
    try {
      const userData = await userApi.fetchUserById(id);
      return User.fromJSON(userData);
    } catch (error) {
      console.error(`Error fetching user with id ${id}:`, error);
      throw error;
    }
  }

  async createUser(userData: { name: string; email: string }): Promise<IUser> {
    try {
      const createdUserData = await userApi.createUser(userData);
      return User.fromJSON(createdUserData);
    } catch (error) {
      console.error("Error creating user:", error);
      throw error;
    }
  }

  async updateUser(
    id: string,
    userData: { name?: string; email?: string }
  ): Promise<IUser> {
    try {
      const updatedUserData = await userApi.updateUser(id, userData);
      return User.fromJSON(updatedUserData);
    } catch (error) {
      console.error(`Error updating user with id ${id}:`, error);
      throw error;
    }
  }

  async deleteUser(id: string): Promise<void> {
    try {
      await userApi.deleteUser(id);
    } catch (error) {
      console.error(`Error deleting user with id ${id}:`, error);
      throw error;
    }
  }

  async getUsersByName(name: string): Promise<IUser[]> {
    try {
      const users = await this.getUsers();
      return users.filter((user) =>
        user.name.toLowerCase().includes(name.toLowerCase())
      );
    } catch (error) {
      console.error(`Error fetching users by name ${name}:`, error);
      throw error;
    }
  }
}
