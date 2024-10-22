import { DishInterface } from "./type_dish";
import { SetInterface, SetProtoDish } from "./types_set";

// Base interface for order items
export interface OrderItemBase {
  id: number;
  quantity: number;
}

// Order item for a dish
export interface DishOrderItem extends OrderItemBase {
  id: number;
  quantity: number;
  dish: DishInterface; // Assuming you have this interface defined elsewhere
}

// Order item for a set
export interface SetOrderItem extends OrderItemBase {
  id: number;
  quantity: number;
  set: SetInterface;
  modifiedDishes: SetProtoDish[];
}

// Main order interface
export interface Order {
  id: number;
  guestId: number;

  tableNumber: number;
  dishSnapshotId: number;

  orderHandlerId: number;
  status: "pending" | "processing" | "completed" | "cancelled";
  createdAt: string;
  updatedAt: string;
  totalPrice: number;
  dishOrderItems: DishOrderItem[];

  setOrderItems: SetOrderItem[];
}

// Interface for creating a new order
export interface CreateOrderBody {
  guestId: number;


  tableNumber: number;
  dishSnapshotId: number;

  orderHandlerId: number;
  status: "pending" | "processing" | "completed" | "cancelled";
  createdAt: string;
  updatedAt: string;
  totalPrice: number;
  dishOrderItems: DishOrderItem[];

  setOrderItems: SetOrderItem[];
}

// export interface Dish {
//   id: number;
//   name: string;
//   price: number;
//   image: string;
//   description: string;
// }

// export interface Set {
//   id: number;
//   name: string;
//   price: number;
//   image: string;
//   description: string;
//   dishes: Dish[];
// }

// export interface SetProtoDish {
//   id: number;
//   name: string;
//   price: number;
//   // Add other necessary fields as needed
// }

// export interface DishOrderItem {
//   id: number;
//   quantity: number;
//   dish: Dish;
// }

// export interface SetOrderItem {
//   id: number;
//   quantity: number;
//   set: Set;
//   modifiedDishes: SetProtoDish[];
// }

// export interface CreateOrderItem {
//   id: number;
//   quantity: number;
//   modifiedDishes?: SetProtoDish[];
// }

// export interface CreateOrderRequest {
//   userId: number;
//   items: CreateOrderItem[];
// }

// export interface Order {
//   id: number;
//   guestId: number;
//   tableNumber: number;
//   dishSnapshotId: number;
//   quantity: number;
//   orderHandlerId: number;
//   status: string;
//   createdAt: Date;
//   updatedAt: Date;
//   totalPrice: number;
//   dishItems: DishOrderItem[];
//   setItems: SetOrderItem[];
// }

// export interface Guest {
//   id: number;
//   name: string;
//   tableNumber: number;
//   createdAt: Date;
//   updatedAt: Date;
// }

// export interface DishSnapshot {
//   id: number;
//   name: string;
//   price: number;
//   image: string;
//   description: string;
//   status: string;
//   dishId: number;
//   createdAt: Date;
//   updatedAt: Date;
// }

// export interface Account {
//   id: number;
//   name: string;
//   email: string;
//   role: string;
//   avatar: string;
// }

// export interface Table {
//   number: number;
//   capacity: number;
//   status: string;
//   token: string;
//   createdAt: Date;
//   updatedAt: Date;
// }

// export interface CreateOrdersRequest {
//   guestId: number;
//   orders: CreateOrderItem[];
// }

// export interface UpdateOrderRequest {
//   orderId: number;
//   status: string;
//   dishId: number;
//   quantity: number;
// }

// export interface PayGuestOrdersRequest {
//   guestId: number;
// }

// export interface GetOrdersRequest {
//   fromDate: Date;
//   toDate: Date;
// }

// export interface OrderResponse {
//   data: Order;
// }

// export interface OrderListResponse {
//   data: Order[];
// }

// export interface OrderIdParam {
//   id: number;
// }

// export interface OrderDetailIdParam {
//   id: number;
// }
