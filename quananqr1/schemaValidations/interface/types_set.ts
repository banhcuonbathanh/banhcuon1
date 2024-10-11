import { Dish } from "./type_dish";

export interface Set {
    id: number;
    name: string;
    description?: string;
    dishes: Dish[];
    created_at: Date;
    updated_at: Date;
    user_id: number;
    is_favourite: number;
    like_by: number[];
  }
  
  export type SetListRes = Set[];
  export type SetListResType = SetListRes;
  
  export type SetType = Set;
  
  // FavoriteSet interfaces and types
  
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