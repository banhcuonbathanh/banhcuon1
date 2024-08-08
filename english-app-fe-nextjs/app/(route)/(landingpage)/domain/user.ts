// domain/User.ts
import { IUser } from "./interface_User";
import { IUserFactory } from "./interface_User_Factory";

export class User implements IUser {
  constructor(public id: string, public name: string, public email: string) {}
}

export class UserFactory implements IUserFactory {
  createUser(id: string, name: string, email: string): IUser {
    return new User(id, name, email);
  }

  fromJSON(json: any): IUser {
    return new User(json.id, json.name, json.email);
  }
}
