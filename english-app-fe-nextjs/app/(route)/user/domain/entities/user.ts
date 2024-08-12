// domain/User.ts

export class User implements IUser {
  constructor(public id: string, public name: string, public email: string) {}

  static fromJSON(json: any): IUser {
    return new User(json.id, json.name, json.email);
  }
}

export class UserFactory implements IUserFactory {
  createUser(id: string, name: string, email: string): IUser {
    return new User(id, name, email);
  }

  fromJSON(json: any): IUser {
    return new User(json.id, json.name, json.email);
  }
}

export interface IUser {
  id: string;
  name: string;
  email: string;
}
// domain/interface_UserFactory.ts

export interface IUserFactory {
  createUser(id: string, name: string, email: string): IUser;
  fromJSON(json: any): IUser;
}
