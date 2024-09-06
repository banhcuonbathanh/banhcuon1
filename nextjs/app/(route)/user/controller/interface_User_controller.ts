// interfaces/IUserController.ts
import { IUser } from "../domain/interface_User";

export interface IUserController {
  users: IUser[];
  loading: boolean;
  error: string | null;
  fetchUsers(): Promise<void>;
  fetchUserById(id: string): Promise<IUser>;
  createUser(userData: { name: string; email: string }): Promise<IUser>;
  updateUser(id: string, userData: { name?: string; email?: string }): Promise<IUser>;
  deleteUser(id: string): Promise<void>;
  searchUsersByName(name: string): Promise<IUser[]>;
}