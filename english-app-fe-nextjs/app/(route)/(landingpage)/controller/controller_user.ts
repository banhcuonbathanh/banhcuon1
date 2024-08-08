
// hooks/useUsers.ts
import { useState, useEffect, useCallback } from 'react';
import { UserService } from '../application/userService';
import { IUser } from '../domain/interface_User';
import { IUserController } from './interface_User_controller';


const userService = new UserService();

export const useUsers = (): IUserController => {
  const [users, setUsers] = useState<IUser[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    try {
      const fetchedUsers = await userService.getUsers();
      setUsers(fetchedUsers);
      setError(null);
    } catch (err) {
      setError('Failed to load users');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const fetchUserById = async (id: string): Promise<IUser> => {
    setLoading(true);
    try {
      const user = await userService.getUserById(id);
      setError(null);
      return user;
    } catch (err) {
      setError(`Failed to fetch user with id ${id}`);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const createUser = async (userData: { name: string; email: string }): Promise<IUser> => {
    setLoading(true);
    try {
      const newUser = await userService.createUser(userData);
      setUsers(prevUsers => [...prevUsers, newUser]);
      setError(null);
      return newUser;
    } catch (err) {
      setError('Failed to create user');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const updateUser = async (id: string, userData: { name?: string; email?: string }): Promise<IUser> => {
    setLoading(true);
    try {
      const updatedUser = await userService.updateUser(id, userData);
      setUsers(prevUsers => prevUsers.map(user => user.id === id ? updatedUser : user));
      setError(null);
      return updatedUser;
    } catch (err) {
      setError(`Failed to update user with id ${id}`);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const deleteUser = async (id: string): Promise<void> => {
    setLoading(true);
    try {
      await userService.deleteUser(id);
      setUsers(prevUsers => prevUsers.filter(user => user.id !== id));
      setError(null);
    } catch (err) {
      setError(`Failed to delete user with id ${id}`);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const searchUsersByName = async (name: string): Promise<IUser[]> => {
    setLoading(true);
    try {
      const matchedUsers = await userService.getUsersByName(name);
      setError(null);
      return matchedUsers;
    } catch (err) {
      setError(`Failed to search users by name ${name}`);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { users, loading, error, fetchUsers, fetchUserById, createUser, updateUser, deleteUser, searchUsersByName };
};