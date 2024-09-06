// src/domain/usecases/getUserUseCase.ts
import { User } from '../entities/user';

export interface GetUserUseCase {
  execute(id: number): Promise<User>;
}
