import { toast } from "react-hot-toast";
import { create } from "zustand";
import { persist } from "zustand/middleware";

// Define the new interfaces as provided
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

export interface OrderVersionSummaryCart {
  version_number: number;
  modification_type: string;
  modified_at: string;
  dishes_ordered: OrderDetailedDishCart[];
  set_ordered: OrderSetDetailedCart[];
}

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
  image: string;
  status: string;
}

export interface PaginationInfoCart {
  current_page: number;
  total_pages: number;
  total_items: number;
  page_size: number;
}

interface OrderSummary {
  totalItems: number;
  totalPrice: number;
  dishes: OrderDetailedDishCart[];
  sets: OrderSetDetailedCart[];
}

interface CartState {
  dishTotal: OrderDetailedDishCart[];
  deliveryData: Record<number, OrderDetailedDishCart>;
  remainingData: OrderDetailedDishCart;
  new_order: OrderDetailedResponseWithDelivery[];
  current_order: OrderDetailedResponseWithDelivery | null;
  isLoading: boolean;
  error: string | null;
  pagination: PaginationInfoCart;
  tableNumber: number | null;
  tableToken: string | null;
  dishState: Record<number, OrderDetailedDishCart>;
  setStore: Record<number, OrderSetDetailedCart>;
  getOrderSummary: () => OrderSummary;
  setIsLoading: (loading: boolean) => void;
  addDishToCart: (dish: OrderDetailedDishCart) => void;
  addSetToCart: (setItem: OrderSetDetailedCart) => void;
  removeDishFromCart: (dishId: number) => void;
  removeSetFromCart: (setId: number) => void;
  updateDishQuantity: (type: "increment" | "decrement", dishId: number) => void;
  updateSetQuantity: (type: "increment" | "decrement", setId: number) => void;
  updateTopping: (topping: string) => void;
  clearCart: () => void;
  initializeOrder: (
    tableNumber: number,
    isGuest: boolean,
    userId?: number,
    guestId?: number
  ) => void;
  addToNewOrder: (order: OrderDetailedResponseWithDelivery) => void;
  removeFromNewOrder: (orderId: number) => void;
  clearNewOrders: () => void;
  updateNewOrderStatus: (orderId: number, status: string) => void;
  getNewOrderById: (orderId: number) => OrderDetailedResponseWithDelivery | undefined;
  addTableNumber: (number: number) => void;
  addTableToken: (token: string) => void;
}

const defaultPagination: PaginationInfoCart = {
  current_page: 1,
  total_pages: 1,
  total_items: 0,
  page_size: 10
};

const defaultOrder: OrderDetailedResponseWithDelivery = {
  id: 0,
  guest_id: 0,
  user_id: 0,
  table_number: 0,
  order_handler_id: 0,
  status: "pending",
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  is_guest: true,
  topping: "",
  tracking_order: "created",
  takeAway: false,
  chiliNumber: 0,
  table_token: "",
  order_name: "",
  current_version: 1,
  version_history: [
    {
      version_number: 1,
      modification_type: "CREATED",
      modified_at: new Date().toISOString(),
      dishes_ordered: [],
      set_ordered: []
    }
  ],
  delivery_history: [],
  current_delivery_status: "PENDING",
  total_items_delivered: 0,
  last_delivery_at: new Date().toISOString()
};

const dishtotlaexample: OrderDetailedDishCart = {
  dish_id: 1,
  quantity: 10,
  name: "dishtotlaexample",
  price: 10,
  description: "dishtotlaexample",
  image: "dishtotlaexample",
  status: "active"
};

const deliveryExample: Record<number, OrderDetailedDishCart> = {
  1: {
    dish_id: 1,
    quantity: 3,
    name: "Pasta Delivery",
    price: 10,
    description: "Fast delivery pasta",
    image: "pasta-delivery",
    status: "active"
  },
  2: {
    dish_id: 2,
    quantity: 1,
    name: "Pizza Delivery",
    price: 5,
    description: "Express pizza delivery",
    image: "pizza-delivery",
    status: "active"
  },
  3: {
    dish_id: 3,
    quantity: 2,
    name: "Sushi Delivery",
    price: 7,
    description: "Fresh sushi delivery",
    image: "sushi-delivery",
    status: "active"
  }
};

const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      dishState: {},
      setStore: {},
      new_order: [],
      current_order: null,
      isLoading: false,
      error: null,
      pagination: defaultPagination,
      tableNumber: null,
      tableToken: null,
      dishTotal: [dishtotlaexample],
      deliveryData: deliveryExample,
      remainingData: {
        dish_id: 0,
        quantity: 0,
        name: "",
        price: 0,
        description: "",
        image: "",
        status: "active"
      },

      addTableNumber: (number: number) => {
        set({ tableNumber: number });
        toast.success("Table number updated");
      },

      setIsLoading: (loading: boolean) => {
        set({ isLoading: loading });
      },

      addTableToken: (token: string) => {
        set({ tableToken: token });
        toast.success("Table token updated");
      },

      getOrderSummary: () => {
        const state = get();
        if (!state.current_order) {
          return {
            totalItems: 0,
            totalPrice: 0,
            dishes: [],
            sets: []
          };
        }

        const currentVersion = state.current_order.version_history.find(
          (v) => v.version_number === state.current_order?.current_version
        );

        if (!currentVersion) {
          return {
            totalItems: 0,
            totalPrice: 0,
            dishes: [],
            sets: []
          };
        }

        const dishes = currentVersion.dishes_ordered;
        const sets = currentVersion.set_ordered;

        const totalPrice =
          dishes.reduce((acc, dish) => acc + dish.price * dish.quantity, 0) +
          sets.reduce((acc, set) => acc + set.price * set.quantity, 0);

        const totalItems =
          dishes.reduce((acc, dish) => acc + dish.quantity, 0) +
          sets.reduce((acc, set) => acc + set.quantity, 0);

        return {
          totalItems,
          totalPrice,
          dishes,
          sets
        };
      },

      addDishToCart: (dish: OrderDetailedDishCart) => {
        set((state) => {
          const currentOrder = state.current_order || { ...defaultOrder };
          const currentVersion = currentOrder.version_history.find(
            (v) => v.version_number === currentOrder.current_version
          );

          if (!currentVersion) {
            toast.error("Current version not found");
            return state;
          }

          const existingDish = currentVersion.dishes_ordered.find(
            (d) => d.dish_id === dish.dish_id
          );

          if (existingDish) {
            toast.error("Dish already exists in cart");
            return state;
          }

          const newDish: OrderDetailedDishCart = { ...dish, quantity: 1 };
          const updatedVersion = {
            ...currentVersion,
            dishes_ordered: [...currentVersion.dishes_ordered, newDish]
          };

          const updatedVersionHistory = currentOrder.version_history.map((v) =>
            v.version_number === currentOrder.current_version ? updatedVersion : v
          );

          return {
            ...state,
            current_order: {
              ...currentOrder,
              version_history: updatedVersionHistory
            }
          };
        });
        toast.success("Dish added successfully");
      },

      addSetToCart: (setItem: OrderSetDetailedCart) => {
        set((state) => {
          const currentOrder = state.current_order || { ...defaultOrder };
          const currentVersion = currentOrder.version_history.find(
            (v) => v.version_number === currentOrder.current_version
          );

          if (!currentVersion) {
            toast.error("Current version not found");
            return state;
          }

          const existingSet = currentVersion.set_ordered.find(
            (s) => s.id === setItem.id
          );

          if (existingSet) {
            toast.error("Set already exists in cart");
            return state;
          }

          const newSet: OrderSetDetailedCart = { ...setItem, quantity: 1 };
          const updatedVersion = {
            ...currentVersion,
            set_ordered: [...currentVersion.set_ordered, newSet]
          };

          const updatedVersionHistory = currentOrder.version_history.map((v) =>
            v.version_number === currentOrder.current_version ? updatedVersion : v
          );

          return {
            ...state,
            current_order: {
              ...currentOrder,
              version_history: updatedVersionHistory
            }
          };
        });
        toast.success("Set added successfully");
      },

      updateTopping: (topping: string) => {
        set((state) => {
          const currentOrder = state.current_order;
          if (!currentOrder) return state;

          return {
            ...state,
            current_order: {
              ...currentOrder,
              topping
            }
          };
        });
        toast.success("Topping updated successfully");
      },

      removeDishFromCart: (dishId: number) => {
        set((state) => {
          const currentOrder = state.current_order;
          if (!currentOrder) return state;

          const currentVersion = currentOrder.version_history.find(
            (v) => v.version_number === currentOrder.current_version
          );

          if (!currentVersion) return state;

          const updatedDishes = currentVersion.dishes_ordered.filter(
            (dish) => dish.dish_id !== dishId
          );

          const updatedVersion = {
            ...currentVersion,
            dishes_ordered: updatedDishes
          };

          const updatedVersionHistory = currentOrder.version_history.map((v) =>
            v.version_number === currentOrder.current_version ? updatedVersion : v
          );

          return {
            ...state,
            current_order: {
              ...currentOrder,
              version_history: updatedVersionHistory
            }
          };
        });
        toast.success("Dish removed from cart");
      },

      removeSetFromCart: (setId: number) => {
        set((state) => {
          const currentOrder = state.current_order;
          if (!currentOrder) return state;

          const currentVersion = currentOrder.version_history.find(
            (v) => v.version_number === currentOrder.current_version
          );

          if (!currentVersion) return state;

          const updatedSets = currentVersion.set_ordered.filter(
            (set) => set.id !== setId
          );

          const updatedVersion = {
            ...currentVersion,
            set_ordered: updatedSets
          };

          const updatedVersionHistory = currentOrder.version_history.map((v) =>
            v.version_number === currentOrder.current_version ? updatedVersion : v
          );

          return {
            ...state,
            current_order: {
              ...currentOrder,
              version_history: updatedVersionHistory
            }
          };
        });
        toast.success("Set removed from cart");
      },

      updateDishQuantity: (type: "increment" | "decrement", dishId: number) => {
        set((state) => {
          const currentOrder = state.current_order;
          if (!currentOrder) return state;

          const currentVersion = currentOrder.version_history.find(
            (v) => v.version_number === currentOrder.current_version
          );

          if (!currentVersion) return state;

          const updatedDishes = currentVersion.dishes_ordered
            .map((dish) => {
              if (dish.dish_id === dishId) {
                const newQuantity =
                  type === "increment" ? dish.quantity + 1 : dish.quantity - 1;
                if (newQuantity <= 0) return null;
                return { ...dish, quantity: newQuantity };
              }
              return dish;
            })
            .filter((dish): dish is OrderDetailedDishCart => dish !== null);

          const updatedVersion = {
            ...currentVersion,
            dishes_ordered: updatedDishes
          };

          const updatedVersionHistory = currentOrder.version_history.map((v) =>
            v.version_number === currentOrder.current_version ? updatedVersion : v
          );

          return {
            ...state,
            current_order: {
              ...currentOrder,
              version_history: updatedVersionHistory
            }
          };
        });
      },

      updateSetQuantity: (type: "increment" | "decrement", setId: number) => {
        set((state) => {
          const currentOrder = state.current_order;
          if (!currentOrder) return state;

          const currentVersion = currentOrder.version_history.find(
            (v) => v.version_number === currentOrder.current_version
          );

          if (!currentVersion) return state;

          const updatedSets = currentVersion.set_ordered
            .map((set) => {
              if (set.id === setId) {
                const newQuantity =
                  type === "increment" ? set.quantity + 1 : set.quantity - 1;
                if (newQuantity <= 0) return null;
                return { ...set, quantity: newQuantity };
              }
              return set;
            })
            .filter((set): set is OrderSetDetailedCart => set !== null);

          const updatedVersion = {
            ...currentVersion,
            set_ordered: updatedSets
          };

          const updatedVersionHistory = currentOrder.version_history.map((v) =>
            v.version_number === currentOrder.current_version ? updatedVersion : v
          );

          return {
            ...state,
            current_order: {
              ...currentOrder,
              version_history: updatedVersionHistory
            }
          };
        });
      },

      clearCart: () => {
        set({ current_order: null });
        toast.success("Cart cleared");
      },

      initializeOrder: (
        tableNumber: number,
        isGuest: boolean,
        userId?: number,
        guestId?: number
      ) => {
        set({
          current_order: {
            ...defaultOrder,
            table_number: tableNumber,
            is_guest: isGuest,
            user_id: userId || 0,
            guest_id: guestId || 0
          }
        });
      },

      addToNewOrder: (order: OrderDetailedResponseWithDelivery) => {
        set((state) => ({
          new_order: [...state.new_order, order]
        }));
      },

      removeFromNewOrder: (orderId: number) => {
        set((state) => ({
          new_order: state.new_order.filter((order) => order.id !== orderId)
        }));
      },

      clearNewOrders: () => {
        set({ new_order: [] });
      },

      updateNewOrderStatus: (orderId: number, status: string) => {
        set((state) => ({
          new_order: state.new_order.map((order) =>
            order.id === orderId ? { ...order, status } : order
          )
        }));
      },

      getNewOrderById: (