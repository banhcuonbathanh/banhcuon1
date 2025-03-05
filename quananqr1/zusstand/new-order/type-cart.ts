export interface DishDeliveryCart {
    id: number;
    order_id: number;
    order_name: string;
    guest_id?: number;
    user_id?: number;
    table_number?: number;
    quantity_delivered: number;
    delivery_status: string;
    delivered_at?: string;
    delivered_by_user_id?: number;
    created_at: string;
    updated_at: string;
    dish_id: number;
    is_guest: boolean;
    modification_number: number;
  }
  
  export type DeliveryStatusCart = 
    | 'PENDING'
    | 'PARTIALLY_DELIVERED'
    | 'FULLY_DELIVERED'
    | 'CANCELLED';
  
  export interface OrderDetailedResponseWithDelivery {
    id: number;
    guest_id: number;
    user_id: number;
    table_number: number;
    order_handler_id: number;
    status: string;
    created_at: string;
    updated_at: string;
    is_guest: boolean;
    topping: string;
    tracking_order: string;
    takeAway: boolean;
    chiliNumber: number;
    table_token: string;
    order_name: string;
    current_version: number;
    version_history: OrderVersionSummaryCart[];
    delivery_history: DishDeliveryCart[];
    current_delivery_status: DeliveryStatusCart;
    total_items_delivered: number;
    last_delivery_at: string;
  }
  
  export interface OrderDetailedListResponseWithDelivery {
    data: OrderDetailedResponseWithDelivery[];
    pagination: PaginationInfoCart;
  }
  
  // Add these missing interfaces if they don't exist:
  export interface OrderVersionSummaryCart {
    version_number: number;
    modification_type: string;
    modified_at: string;
    dishes_ordered: OrderDetailedDishCart[];
    set_ordered: OrderSetDetailedCart[];
  }
  
  // Update existing interfaces with missing fields if needed:
  export interface OrderSetDetailedCart {
    id: number;
    name: string;
    description: string;
    dishes: OrderDetailedDishCart[];
    userId: number;
    created_at: string;
    updated_at: string;
    is_favourite: boolean;
    like_by: number[];
    is_public: boolean;
    image: string;
    price: number;
    quantity: number;
  }
  
  export interface OrderDetailedDishCart {
    dish_id: number;
    quantity: number;
    name: string;
    price: number;
    description: string;
    image: string;  // Fixed typo from 'iamge'
    status: string;
  }
  
  export interface PaginationInfoCart {
    current_page: number;
    total_pages: number;
    total_items: number;
    page_size: number;
  }