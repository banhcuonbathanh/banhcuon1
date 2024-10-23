import { DishInterface } from "./type_dish";
import { SetInterface, SetProtoDish } from "./types_set";

// Base interface for order items

// Order item for a dish
export interface DishOrderItem {
  dish_id: number;
  quantity: number;
}

// Order item for a set
export interface SetOrderItem {
  set_id: number;
  quantity: number;
}

// Main order interface
export interface Order {
  id: number;
  guest_id: number;
  user_id: number;
  is_guest: boolean;
  table_number: number;
  order_handler_id: number;
  status: string;
  created_at: string;
  updated_at: string;
  total_price: number;
  dish_items: DishOrderItem[];
  set_items: SetOrderItem[];
  bow_chili: number;
  bow_no_chili: number;

  // new
  takeAway: boolean;
  chiliNumber: number
}

export interface CreateOrderRequest {
  guest_id?: number | null;
  user_id: number;
  is_guest: boolean;
  table_number: number;
  order_handler_id: number;
  status: string;
  created_at: string;
  updated_at: string;
  total_price: number;
  dish_items: DishOrderItem[];
  set_items: SetOrderItem[];
  bow_chili: number;
  bow_no_chili: number;
//
  takeAway: boolean;
  chiliNumber: number
}

export interface UpdateOrderRequest {
  id: number;
  guest_id: number;
  user_id: number;
  table_number: number;
  order_handler_id: number;
  status: string;
  total_price: number;
  dish_items: DishOrderItem[];
  set_items: SetOrderItem[];
  is_guest: boolean;
  bow_chili: number;
  bow_no_chili: number;
}

export interface GetOrdersRequest {
  from_date: string;
  to_date: string;
  user_id?: number;
  guest_id?: number;
}

export interface PayOrdersRequest {
  guest_id?: number;
  user_id?: number;
}

// Response interfaces
export interface OrderResponse {
  data: Order;
}

export interface OrderListResponse {
  data: Order[];
}

export interface OrderDetailedListResponse {
  data: OrderSetDetailed[];
}

// Parameter interfaces
export interface OrderIDParam {
  id: number;
}

export interface OrderDetailIDParam {
  id: number;
}

// Guest interface
export interface Guest {
  id: number;
  name: string;
  table_number: number;
  created_at: string;
  updated_at: string;
}

export interface OrderSetDetailed {
  id: number;
  name: string;
  description: string;
  dishes: DishOrderItem[];
  userId: number;
  created_at: string;
  updated_at: string;
  is_favourite: boolean;
  like_by: number[];
  is_public: boolean;
  image: string;
  price: number;
}
