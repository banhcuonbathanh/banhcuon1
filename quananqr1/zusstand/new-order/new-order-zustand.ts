import {
  Order,
  DishOrderItem,
  SetOrderItem,
  PaginationInfo
} from "@/schemaValidations/interface/type_order";
import toast from "react-hot-toast";
import { create } from "zustand";
import { persist } from "zustand/middleware";
// dishTotal: DishOrderItem,
// deliveryData : DishOrderItem,
// remainingData : DishOrderItem,

interface OrderSummary {
  totalItems: number;
  totalPrice: number;
  dishes: DishOrderItem[];
  sets: SetOrderItem[];
}
interface CartState {
  // new
  // Added new fields
  dishTotal: DishOrderItem;
  deliveryData: DishOrderItem;
  remainingData: DishOrderItem;
  //
  new_order: Order[];
  current_order: Order | null;
  isLoading: boolean;
  error: string | null;
  pagination: PaginationInfo;
  tableNumber: number | null;
  tableToken: string | null;

  //
  dishState: Record<number, DishOrderItem>;
  setStore: Record<number, SetOrderItem>;
  getOrderSummary: () => OrderSummary;
  //
  setIsLoading: (loading: boolean) => void;
  addDishToCart: (dish: DishOrderItem) => void;
  addSetToCart: (setItem: SetOrderItem) => void;
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
  addToNewOrder: (order: Order) => void;
  removeFromNewOrder: (orderId: number) => void;
  clearNewOrders: () => void;
  updateNewOrderStatus: (orderId: number, status: string) => void;
  getNewOrderById: (orderId: number) => Order | undefined;
  addTableNumber: (number: number) => void;
  addTableToken: (token: string) => void;
}

const defaultPagination: PaginationInfo = {
  current_page: 1,
  total_pages: 1,
  total_items: 0,
  page_size: 10
};

const defaultOrder: Order = {
  id: 0,
  guest_id: 0,
  user_id: 0,
  is_guest: true,
  table_number: 0,
  order_handler_id: 0,
  status: "pending",
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  total_price: 0,
  dish_items: [],
  set_items: [],
  topping: "",
  tracking_order: "created",
  takeAway: false,
  chiliNumber: 0,
  table_token: "",
  order_name: ""
};
const defaultDishOrderItem: DishOrderItem = {
  dish_id: 0,
  quantity: 0,
  name: "",
  price: 0,
  description: "",
  image: "",
  status: "active",
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  is_favourite: false,
  like_by: []
};

const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      //
      dishState: {},
      setStore: {},
      //
      new_order: [],
      current_order: null,
      isLoading: false,
      error: null,
      pagination: defaultPagination,
      tableNumber: null,
      tableToken: null,
      // Initialize the new fields with default values
      dishTotal: { ...defaultDishOrderItem },
      deliveryData: { ...defaultDishOrderItem },
      remainingData: { ...defaultDishOrderItem },
      // Table management functions
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
      //
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

        // Transform to DishOrderItem
        const detailedDishes: DishOrderItem[] =
          state.current_order.dish_items.map((item) => ({
            dish_id: item.dish_id,
            quantity: item.quantity,
            name: state.dishState[item.dish_id]?.name || "",
            price: state.dishState[item.dish_id]?.price || 0,
            description: state.dishState[item.dish_id]?.description || "",
            image: state.dishState[item.dish_id]?.image || "",
            status: state.dishState[item.dish_id]?.status || "active",
            created_at:
              state.dishState[item.dish_id]?.created_at ||
              new Date().toISOString(),
            updated_at:
              state.dishState[item.dish_id]?.updated_at ||
              new Date().toISOString(),
            is_favourite: state.dishState[item.dish_id]?.is_favourite || false,
            like_by: state.dishState[item.dish_id]?.like_by || []
          }));

        // Transform to SetOrderItem
        const detailedSets: SetOrderItem[] = state.current_order.set_items.map(
          (item) => {
            const setDetails = state.setStore[item.set_id] || {};
            return {
              set_id: item.set_id,
              quantity: item.quantity,
              name: setDetails.name || "",
              price: setDetails.price || 0,
              description: setDetails.description || "",
              image: setDetails.image || "",
              dishes: (setDetails.dishes || []).map((dish) => ({
                dish_id: dish.dish_id,
                quantity: dish.quantity,
                name: dish.name || "",
                price: dish.price || 0,
                description: dish.description || "",
                image: dish.image || "",
                status: dish.status || "active",
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
                is_favourite: false,
                like_by: []
              })),
              status: setDetails.status || "active",
              created_at: setDetails.created_at || new Date().toISOString(),
              updated_at: setDetails.updated_at || new Date().toISOString(),
              is_favourite: setDetails.is_favourite || false,
              like_by: setDetails.like_by || [],
              is_public: setDetails.is_public || false,
              userId: setDetails.userId || 0
            };
          }
        );

        // Calculate total price including sets
        const totalPrice =
          detailedDishes.reduce(
            (acc, dish) => acc + dish.price * dish.quantity,
            0
          ) +
          detailedSets.reduce((acc, set) => acc + set.price * set.quantity, 0);

        return {
          totalItems:
            state.current_order.dish_items.reduce(
              (acc, item) => acc + item.quantity,
              0
            ) +
            state.current_order.set_items.reduce(
              (acc, item) => acc + item.quantity,
              0
            ),
          totalPrice,
          dishes: detailedDishes,
          sets: detailedSets
        };
      },
      //
      // Current order management
      addDishToCart: (dish: DishOrderItem) => {
        set((state) => {
          const currentOrder = state.current_order || { ...defaultOrder };
          const existingDish = currentOrder.dish_items.find(
            (item) => item.dish_id === dish.dish_id
          );

          if (existingDish) {
            toast.error("Dish already exists in cart");
            return state;
          }

          const newDishItems = [
            ...currentOrder.dish_items,
            { ...dish, quantity: 1 }
          ];

          return {
            ...state,
            current_order: {
              ...currentOrder,
              dish_items: newDishItems
            }
          };
        });
        toast.success("Dish added successfully");
      },

      addSetToCart: (setItem: SetOrderItem) => {
        set((state) => {
          const currentOrder = state.current_order || { ...defaultOrder };

          if (!setItem.set_id) {
            toast.error("Invalid set configuration");
            return state;
          }

          const existingSet = currentOrder.set_items.find(
            (item) => item.set_id === setItem.set_id
          );

          if (existingSet) {
            toast.error("Set already exists in cart");
            return state;
          }

          const newSetItems = [
            ...currentOrder.set_items,
            { ...setItem, quantity: 1 }
          ];

          return {
            ...state,
            current_order: {
              ...currentOrder,
              set_items: newSetItems
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

          const newDishItems = currentOrder.dish_items.filter(
            (item) => item.dish_id !== dishId
          );

          return {
            ...state,
            current_order: {
              ...currentOrder,
              dish_items: newDishItems
            }
          };
        });
        toast.success("Dish removed from cart");
      },

      removeSetFromCart: (setId: number) => {
        set((state) => {
          const currentOrder = state.current_order;
          if (!currentOrder) return state;

          const newSetItems = currentOrder.set_items.filter(
            (item) => item.set_id !== setId
          );

          return {
            ...state,
            current_order: {
              ...currentOrder,
              set_items: newSetItems
            }
          };
        });
        toast.success("Set removed from cart");
      },

      updateDishQuantity: (type: "increment" | "decrement", dishId: number) => {
        set((state) => {
          const currentOrder = state.current_order;
          if (!currentOrder) return state;

          const newDishItems = currentOrder.dish_items
            .map((item) => {
              if (item.dish_id === dishId) {
                const newQuantity =
                  type === "increment" ? item.quantity + 1 : item.quantity - 1;

                if (newQuantity <= 0) {
                  return null; // Will be filtered out
                }

                return {
                  ...item,
                  quantity: newQuantity
                };
              }
              return item;
            })
            .filter((item): item is DishOrderItem => item !== null);

          return {
            ...state,
            current_order: {
              ...currentOrder,
              dish_items: newDishItems
            }
          };
        });
      },
      // Update the updateSetQuantity function in your useCartStore

      updateSetQuantity: (type: "increment" | "decrement", setId: number) => {
        set((state) => {
          const currentOrder = state.current_order;
          if (!currentOrder) return state;

          // Log for debugging
          console.log("Updating quantity for set:", {
            setId,
            type,
            currentItems: currentOrder.set_items
          });

          const newSetItems = currentOrder.set_items
            .map((item) => {
              // Only update the specific set
              if (item.set_id === setId) {
                const newQuantity =
                  type === "increment"
                    ? (item.quantity || 0) + 1
                    : (item.quantity || 0) - 1;

                // Log the quantity change
                console.log("Quantity change for set:", {
                  setId,
                  oldQuantity: item.quantity,
                  newQuantity
                });

                // Remove the item if quantity becomes 0
                if (newQuantity <= 0) return null;

                // Return updated item
                return {
                  ...item,
                  quantity: newQuantity
                };
              }
              // Return other items unchanged
              return item;
            })
            .filter((item): item is SetOrderItem => item !== null);

          // Create new state immutably
          return {
            ...state,
            current_order: {
              ...currentOrder,
              set_items: newSetItems
            }
          };
        });
      },
      clearCart: () => {
        set((state) => ({
          ...state,
          current_order: null
        }));
        toast.success("Cart cleared");
      },

      initializeOrder: (
        tableNumber: number,
        isGuest: boolean,
        userId?: number,
        guestId?: number
      ) => {
        set((state) => ({
          ...state,
          current_order: {
            ...defaultOrder,
            table_number: tableNumber,
            is_guest: isGuest,
            user_id: userId || 0,
            guest_id: guestId || 0
          }
        }));
      },

      // New order management
      addToNewOrder: (order: Order) => {
        set((state) => ({
          ...state,
          new_order: [...state.new_order, order]
        }));
      },

      removeFromNewOrder: (orderId: number) => {
        set((state) => ({
          ...state,
          new_order: state.new_order.filter((order) => order.id !== orderId)
        }));
      },

      clearNewOrders: () => {
        set((state) => ({
          ...state,
          new_order: []
        }));
      },

      updateNewOrderStatus: (orderId: number, status: string) => {
        set((state) => ({
          ...state,
          new_order: state.new_order.map((order) =>
            order.id === orderId ? { ...order, status } : order
          )
        }));
      },

      getNewOrderById: (orderId: number) => {
        return get().new_order.find((order) => order.id === orderId);
      }
    }),
    {
      name: "cart-storage",
      partialize: (state) => ({
        new_order: state.new_order,
        current_order: state.current_order,
        tableNumber: state.tableNumber,
        tableToken: state.tableToken
      })
    }
  )
);

export default useCartStore;
