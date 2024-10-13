
import { DishInterface } from "./type_dish";
export interface SetProtoDish {
  dishId: number;  // Changed from dish_id to dishId for consistency
  quantity: number;
  dish: DishInterface;  // Assuming Dish interface is defined elsewhere
}

export interface SetInterface {
  id: number;
  name: string;
  description?: string;
  dishes: SetProtoDish[];
  userId: number; // Changed from user_id to userId
  created_at: string; // Changed from Date to string
  updated_at: string; // Changed from Date to string
  is_favourite: boolean; // Changed from number to boolean
  like_by: number[] | null; // Changed to allow null
  is_public: boolean; // Added new field
}

export interface SetListResponse {
  data: SetInterface[];
  message: string;
}

export type SetListResType = SetListResponse;

export type SetIntefaceType = SetInterface;

// FavoriteSet interfaces and types
// Note: We'll keep these as they are since you didn't provide new information about them

export interface FavoriteSet {
  id: number;
  userId: number;
  name: string;
  dishes: number[]; // Array of dish IDs
  createdAt: Date;
  updatedAt: Date;
}

export type FavoriteSetListRes = FavoriteSet[];
export type FavoriteSetListResType = FavoriteSetListRes;