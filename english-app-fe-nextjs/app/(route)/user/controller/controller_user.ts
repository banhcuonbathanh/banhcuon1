// stores/useUserStore.ts
import {create} from "zustand";
import { UserService } from "../application/userService";
import { IUser } from "../domain/interface_User";

const userService = new UserService();

interface UserState {
  users: IUser[];
  loading: boolean;
  error: string | null;
  fetchUsers: () => Promise<void>;
  fetchUserById: (id: string) => Promise<IUser>;
  createUser: (userData: { name: string; email: string }) => Promise<IUser>;
  updateUser: (
    id: string,
    userData: { name?: string; email?: string }
  ) => Promise<IUser>;
  deleteUser: (id: string) => Promise<void>;
  searchUsersByName: (name: string) => Promise<IUser[]>;
}

export const useUserStore = create<UserState>((set, get) => ({
  users: [],
  loading: false,
  error: null,

  fetchUsers: async () => {
    console.log(
      "this is fetchUsers in english-app-fe-nextjs/app/(route)/user/controller/controller_user.ts"
    );
    set({ loading: true });
    try {
      const fetchedUsers = await userService.getUsers();
      set({ users: fetchedUsers, error: null });
    } catch (err) {
      set({ error: "Failed to load users" });
    } finally {
      set({ loading: false });
    }
  },

  fetchUserById: async (id: string) => {
    set({ loading: true });
    try {
      const user = await userService.getUserById(id);
      set({ error: null });
      return user;
    } catch (err) {
      set({ error: `Failed to fetch user with id ${id}` });
      throw err;
    } finally {
      set({ loading: false });
    }
  },

  createUser: async (userData: { name: string; email: string }) => {
    set({ loading: true });
    try {
      const newUser = await userService.createUser(userData);
      set((state) => ({ users: [...state.users, newUser], error: null }));
      return newUser;
    } catch (err) {
      set({ error: "Failed to create user" });
      throw err;
    } finally {
      set({ loading: false });
    }
  },

  updateUser: async (
    id: string,
    userData: { name?: string; email?: string }
  ) => {
    set({ loading: true });
    try {
      const updatedUser = await userService.updateUser(id, userData);
      set((state) => ({
        users: state.users.map((user) => (user.id === id ? updatedUser : user)),
        error: null
      }));
      return updatedUser;
    } catch (err) {
      set({ error: `Failed to update user with id ${id}` });
      throw err;
    } finally {
      set({ loading: false });
    }
  },

  deleteUser: async (id: string) => {
    set({ loading: true });
    try {
      await userService.deleteUser(id);
      set((state) => ({
        users: state.users.filter((user) => user.id !== id),
        error: null
      }));
    } catch (err) {
      set({ error: `Failed to delete user with id ${id}` });
      throw err;
    } finally {
      set({ loading: false });
    }
  },

  searchUsersByName: async (name: string) => {
    set({ loading: true });
    try {
      const matchedUsers = await userService.getUsersByName(name);
      set({ error: null });
      return matchedUsers;
    } catch (err) {
      set({ error: `Failed to search users by name ${name}` });
      throw err;
    } finally {
      set({ loading: false });
    }
  }
}));
