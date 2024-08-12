// src/lib/api.ts

import { User } from "../(route)/user/domain/entities/user";


const API_URL = 'http://localhost:8080'; // Replace with your Go backend URL

export async function fetchUser(id: number): Promise<User> {
  const response = await fetch(`${API_URL}/users/${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch user');
  }
  return response.json();
}
