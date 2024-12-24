import { create } from "zustand";
import { persist } from "zustand/middleware";
import { toast } from "@/components/ui/use-toast";
import envConfig from "@/config";
import { CreateOrderRequest } from "@/schemaValidations/interface/type_order";
import { logWithLevel } from "@/lib/logger/log";
import {
  DishOrderItem,
  Order,
  OrderDetailedDish,
  OrderSetDetailed,
  PaginationInfo,
  SetOrderItem
} from "@/schemaValidations/interface/type_order";
interface DishState {
  [key: number]: {
    name: string;
    price: number;
    description: string;
    image: string;
    status: string;
  }
}

const LOG_PATH =
  "quananqr1/app/(client)/table/[number]/component/order/logic.ts";

interface BowlOptions {
  canhKhongRau: number;
  canhCoRau: number;
  smallBowl: number;
  wantChili: boolean;
  selectedFilling: {
    mocNhi: boolean;
    thit: boolean;
    thitMocNhi: boolean;
  };
}

interface OrderSummary {
  totalItems: number;
  totalPrice: number;
  dishes: OrderDetailedDish[];
  sets: OrderSetDetailed[];
}

interface OrderState extends BowlOptions {
  //

  dishState: DishState;
  setDishDetails: (id: number, details: Omit<DishState[number], 'id'>) => void;
  //
  orders: Record<string, Order[]>;
  currentOrder: Order | null;
  isLoading: boolean;
  error: string | null;
  pagination: PaginationInfo;
  tableNumber: number | null;
  tableToken: string | null;

  // Table management
  addTableNumber: (number: number) => void;
  addTableToken: (token: string) => void;

  // Bowl options management
  updateCanhKhongRau: (count: number) => void;
  updateCanhCoRau: (count: number) => void;
  updateSmallBowl: (count: number) => void;
  updateWantChili: (value: boolean) => void;
  updateSelectedFilling: (type: "mocNhi" | "thit" | "thitMocNhi") => void;

  // Order management
  setOrders: (key: string, orders: Order[]) => void;
  addOrderToCategory: (key: string, order: Order) => void;
  updateOrderInCategory: (
    key: string,
    orderId: number,
    updates: Partial<Order>
  ) => void;
  removeOrderFromCategory: (key: string, orderId: number) => void;
  getOrdersByCategory: (key: string) => Order[];
  initializeCategory: (key: string) => void;
  removeCategory: (key: string) => void;

  // Current order management
  setCurrentOrder: (order: Order | null) => void;
  getOrderSummary: () => OrderSummary;

  // Current order item management
  addDishToCurrentOrder: (dishItem: DishOrderItem) => void;
  removeDishFromCurrentOrder: (dishId: number) => void;
  updateDishQuantityInCurrentOrder: (dishId: number, quantity: number) => void;
  addSetToCurrentOrder: (setItem: SetOrderItem) => void;
  removeSetFromCurrentOrder: (setId: number) => void;
  updateSetQuantityInCurrentOrder: (setId: number, quantity: number) => void;

  // Utilities
  clearCurrentOrder: () => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  setPagination: (pagination: PaginationInfo) => void;

  // Order creation
  createOrder: (params: {
    topping: string;
    Table_token: string;
    http: any;
    auth: {
      guest: any;
      user: any;
      isGuest: boolean;
    };
    orderStore: {
      tableNumber: number;
      getOrderSummary: () => any;
    };
    websocket: {
      disconnect: () => void;
      isConnected: boolean;
      sendMessage: (message: any) => void;
    };
    openLoginDialog: () => void;
  }) => Promise<any>;
}

const INITIAL_ORDER: Order = {
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
  tracking_order: "",
  takeAway: false,
  chiliNumber: 0,
  table_token: "",
  order_name: ""
};

const INITIAL_STATE: Omit<
  OrderState,
  | "setOrders"
  | "addOrderToCategory"
  | "updateOrderInCategory"
  | "removeOrderFromCategory"
  | "getOrdersByCategory"
  | "initializeCategory"
  | "removeCategory"
  | "setCurrentOrder"
  | "getOrderSummary"
  | "addDishToCurrentOrder"
  | "removeDishFromCurrentOrder"
  | "updateDishQuantityInCurrentOrder"
  | "addSetToCurrentOrder"
  | "removeSetFromCurrentOrder"
  | "updateSetQuantityInCurrentOrder"
  | "clearCurrentOrder"
  | "setLoading"
  | "setError"
  | "setPagination"
  | "addTableNumber"
  | "addTableToken"
  | "updateCanhKhongRau"
  | "updateCanhCoRau"
  | "updateSmallBowl"
  | "updateWantChili"
  | "updateSelectedFilling"
  | "createOrder"
> = {
  dishState: {} as DishState,
  orders: {},
  currentOrder: null,
  isLoading: false,
  error: null,
  pagination: {
    current_page: 1,
    total_pages: 1,
    total_items: 0,
    page_size: 10
  },
  canhKhongRau: 0,
  canhCoRau: 0,
  smallBowl: 0,
  wantChili: false,
  selectedFilling: {
    mocNhi: false,
    thit: false,
    thitMocNhi: false
  },
  tableNumber: null,
  tableToken: null
};

const useOrderStore = create<OrderState>()(
  persist(
    (set, get) => ({
      ...INITIAL_STATE,
      setDishDetails: (id, details) => 
        set(state => ({
          dishState: {
            ...state.dishState,
            [id]: details
          }
        })),
      
      createOrder: async ({
        topping,
        Table_token,
        http,
        auth: { guest, user, isGuest },
        orderStore: { tableNumber, getOrderSummary },
        websocket: { disconnect, isConnected, sendMessage },
        openLoginDialog
      }) => {
        logWithLevel(
          {
            isGuest,
            hasUser: !!user,
            hasGuest: !!guest,
            tableNumber,
            websocketConnected: isConnected
          },
          LOG_PATH,
          "debug",
          1
        );

        if (!user && !guest) {
          logWithLevel(
            { error: "No user or guest found" },
            LOG_PATH,
            "warn",
            7
          );
          openLoginDialog();
          return;
        }

        const orderSummary = getOrderSummary();

        logWithLevel(
          {
            orderSummary,
            dishCount: orderSummary.dishes.length,
            setCount: orderSummary.sets.length
          },
          LOG_PATH,
          "info",
          11
        );

        const dish_items = orderSummary.dishes.map((dish: any) => ({
          dish_id: dish.id,
          quantity: dish.quantity
        }));

        const set_items = orderSummary.sets.map((set: any) => ({
          set_id: set.id,
          quantity: set.quantity
        }));

        const user_id = isGuest ? null : user?.id ?? null;
        const guest_id = isGuest ? guest?.id ?? null : null;
        let order_name = "";
        if (isGuest && guest) {
          order_name = guest.name;
        } else if (!isGuest && user) {
          order_name = user.name;
        }

        const orderData: CreateOrderRequest = {
          guest_id,
          user_id,
          is_guest: isGuest,
          table_number: tableNumber,
          order_handler_id: 1,
          status: "pending",
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          total_price: orderSummary.totalPrice,
          dish_items,
          set_items,
          topping,
          tracking_order: "tracking_order",
          takeAway: false,
          chiliNumber: 0,
          table_token: Table_token,
          order_name
        };

        logWithLevel(
          {
            orderData,
            validationStatus: "prepared"
          },
          LOG_PATH,
          "debug",
          12
        );

        set({ isLoading: true });

        try {
          logWithLevel(
            {
              endpoint: `${envConfig.NEXT_PUBLIC_API_ENDPOINT}${envConfig.Order_External_End_Point}`,
              requestData: orderData
            },
            LOG_PATH,
            "info",
            2
          );

          const link_order = `${envConfig.NEXT_PUBLIC_API_ENDPOINT}${envConfig.Order_External_End_Point}`;
          const response = await http.post(link_order, orderData);

          logWithLevel(
            {
              orderId: response,
              status: "order response"
            },
            LOG_PATH,
            "info",
            8
          );

          if (isConnected) {
            logWithLevel(
              {
                messageType: "NEW_ORDER",
                orderId: response.data.id
              },
              LOG_PATH,
              "debug",
              6
            );

            sendMessage({
              type: "NEW_ORDER",
              data: {
                orderId: response.data.id,
                orderData
              }
            });
          }

          toast({
            title: "Success",
            description: "Order has been created successfully"
          });

          return response.data;
        } catch (error) {
          logWithLevel(
            {
              error: error instanceof Error ? error.message : "Unknown error",
              orderData
            },
            LOG_PATH,
            "error",
            7
          );

          console.error("Order creation failed:", error);
          toast({
            variant: "destructive",
            title: "Error",
            description:
              error instanceof Error ? error.message : "Failed to create order"
          });
          throw error;
        } finally {
          set({ isLoading: false });
          disconnect();
        }
      },

      // Table management
      addTableNumber: (number) => set({ tableNumber: number }),
      addTableToken: (token) => set({ tableToken: token }),

      // Bowl options management
      updateCanhKhongRau: (count) => set({ canhKhongRau: count }),
      updateCanhCoRau: (count) => set({ canhCoRau: count }),
      updateSmallBowl: (count) => set({ smallBowl: count }),
      updateWantChili: (value) => set({ wantChili: value }),
      updateSelectedFilling: (type) =>
        set({
          selectedFilling: {
            mocNhi: type === "mocNhi",
            thit: type === "thit",
            thitMocNhi: type === "thitMocNhi"
          }
        }),

      // Order management
      setOrders: (key, orders) =>
        set((state) => ({
          orders: { ...state.orders, [key]: orders }
        })),

      addOrderToCategory: (key, order) =>
        set((state) => ({
          orders: {
            ...state.orders,
            [key]: [...(state.orders[key] || []), order]
          }
        })),

      updateOrderInCategory: (key, orderId, updates) =>
        set((state) => ({
          orders: {
            ...state.orders,
            [key]: (state.orders[key] || []).map((order) =>
              order.id === orderId ? { ...order, ...updates } : order
            )
          }
        })),

      removeOrderFromCategory: (key, orderId) =>
        set((state) => ({
          orders: {
            ...state.orders,
            [key]: (state.orders[key] || []).filter(
              (order) => order.id !== orderId
            )
          }
        })),

      getOrdersByCategory: (key) => {
        const state = get();
        return state.orders[key] || [];
      },

      initializeCategory: (key) =>
        set((state) => ({
          orders: { ...state.orders, [key]: [] }
        })),

      removeCategory: (key) =>
        set((state) => {
          const { [key]: _, ...remainingOrders } = state.orders;
          return { orders: remainingOrders };
        }),

      // Current order management
      setCurrentOrder: (order) => set({ currentOrder: order }),

      getOrderSummary: () => {
        const state = get();
        if (!state.currentOrder) {
          return {
            totalItems: 0,
            totalPrice: 0,
            dishes: [],
            sets: []
          };
        }
      
        const totalItems =
          state.currentOrder.dish_items.reduce(
            (acc, item) => acc + item.quantity,
            0
          ) +
          state.currentOrder.set_items.reduce(
            (acc, item) => acc + item.quantity,
            0
          );
      
        // Transform DishOrderItem into OrderDetailedDish
        const detailedDishes: OrderDetailedDish[] = state.currentOrder.dish_items.map(item => ({
          id: item.dish_id,
          quantity: item.quantity,
          name: dishState[item.dish_id]?.name || "",
          price: dishState[item.dish_id]?.price || 0,
          description: dishState[item.dish_id]?.description || "",
          image: dishState[item.dish_id]?.image || "",
          status: dishState[item.dish_id]?.status || "active"
        }));
      
        // Transform SetOrderItem into OrderSetDetailed
        const detailedSets: OrderSetDetailed[] = state.currentOrder.set_items.map(item => ({
          id: item.set_id,
          quantity: item.quantity,
          // Add other required set details here
          name: "",  // Add set name from state
          price: 0,  // Add set price from state
          description: "", // Add set description from state
          image: "",  // Add set image from state
          status: "active"
        }));
      
        return {
          totalItems,
          totalPrice: state.currentOrder.total_price,
          dishes: detailedDishes,
          sets: detailedSets
        };
      },
      // Current order item management
      addDishToCurrentOrder: (dishItem) =>
        set((state) => {
          if (!state.currentOrder) {
            return {
              currentOrder: {
                ...INITIAL_ORDER,
                dish_items: [dishItem]
              }
            };
          }

          const existingItemIndex = state.currentOrder.dish_items.findIndex(
            (item) => item.dish_id === dishItem.dish_id
          );

          if (existingItemIndex >= 0) {
            const updatedItems = state.currentOrder.dish_items.map(
              (item, index) =>
                index === existingItemIndex
                  ? { ...item, quantity: item.quantity + dishItem.quantity }
                  : item
            );

            return {
              currentOrder: {
                ...state.currentOrder,
                dish_items: updatedItems
              }
            };
          }

          return {
            currentOrder: {
              ...state.currentOrder,
              dish_items: [...state.currentOrder.dish_items, dishItem]
            }
          };
        }),

      removeDishFromCurrentOrder: (dishId) =>
        set((state) => {
          if (!state.currentOrder) return state;

          return {
            currentOrder: {
              ...state.currentOrder,
              dish_items: state.currentOrder.dish_items.filter(
                (item) => item.dish_id !== dishId
              )
            }
          };
        }),

      updateDishQuantityInCurrentOrder: (dishId, quantity) =>
        set((state) => {
          if (!state.currentOrder) return state;

          return {
            currentOrder: {
              ...state.currentOrder,
              dish_items: state.currentOrder.dish_items.map((item) =>
                item.dish_id === dishId ? { ...item, quantity } : item
              )
            }
          };
        }),

      addSetToCurrentOrder: (setItem) =>
        set((state) => {
          if (!state.currentOrder) {
            return {
              currentOrder: {
                ...INITIAL_ORDER,
                set_items: [setItem]
              }
            };
          }

          const existingItemIndex = state.currentOrder.set_items.findIndex(
            (item) => item.set_id === setItem.set_id
          );

          if (existingItemIndex >= 0) {
            const updatedItems = state.currentOrder.set_items.map(
              (item, index) =>
                index === existingItemIndex
                  ? { ...item, quantity: item.quantity + setItem.quantity }
                  : item
            );

            return {
              currentOrder: {
                ...state.currentOrder,
                set_items: updatedItems
              }
            };
          }

          return {
            currentOrder: {
              ...state.currentOrder,
              set_items: [...state.currentOrder.set_items, setItem]
            }
          };
        }),

      removeSetFromCurrentOrder: (setId) =>
        set((state) => {
          if (!state.currentOrder) return state;

          return {
            currentOrder: {
              ...state.currentOrder,
              set_items: state.currentOrder.set_items.filter(
                (item) => item.set_id !== setId
              )
            }
          };
        }),

      updateSetQuantityInCurrentOrder: (setId, quantity) =>
        set((state) => {
          if (!state.currentOrder) return state;

          return {
            currentOrder: {
              ...state.currentOrder,
              set_items: state.currentOrder.set_items.map((item) =>
                item.set_id === setId ? { ...item, quantity } : item
              )
            }
          };
        }),

      // Utilities
      clearCurrentOrder: () =>
        set({
          currentOrder: {
            ...INITIAL_ORDER,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString()
          }
        }),

      setLoading: (isLoading) => set({ isLoading }),
      setError: (error) => set({ error }),
      setPagination: (pagination) => set({ pagination })
    }),
    {
      name: "order-storage",
      partialize: (state) => ({
        orders: state.orders,
        currentOrder: state.currentOrder,
        pagination: state.pagination,
        tableNumber: state.tableNumber,
        tableToken: state.tableToken,
        canhKhongRau: state.canhKhongRau,
        canhCoRau: state.canhCoRau,
        smallBowl: state.smallBowl,
        wantChili: state.wantChili,
        selectedFilling: state.selectedFilling
      })
    }
  )
);

export default useOrderStore;
