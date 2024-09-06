// src/application/services/userService.ts

import { User } from '../../domain/entities/user';
import { GetUserUseCase } from '../../domain/usecase/getUserUseCase';

export class UserService {
  constructor(private getUserUseCase: GetUserUseCase) {}

  async getUser(id: number): Promise<User> {
    return this.getUserUseCase.execute(id);
  }
}