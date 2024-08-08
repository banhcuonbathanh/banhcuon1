// domain/interface_UserFactory.ts
import { IUser } from "./interface_User";

export interface IUserFactory {
  createUser(id: string, name: string, email: string): IUser;
  fromJSON(json: any): IUser;
}
