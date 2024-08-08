// interfaces/IUserService.ts
import { IUser } from "../domain/interface_User";

export interface IUserService {
  getUsers(): Promise<IUser[]>;
  getUserById(id: string): Promise<IUser>;
  createUser(userData: { name: string; email: string }): Promise<IUser>;
  updateUser(id: string, userData: { name?: string; email?: string }): Promise<IUser>;
  deleteUser(id: string): Promise<void>;
  getUsersByName(name: string): Promise<IUser[]>;
}
