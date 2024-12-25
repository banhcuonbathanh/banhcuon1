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
  };
}
interface SetState {
  [key: number]: {
    name: string;
    price: number;
    description: string;
    image: string;
    dishes: Array<{
      dish_id: number;
      name: string;
      price: number;
      quantity: number;
      description: string;
      image: string;
      status: string;
    }>;
    userId: number;
    created_at: string;
    updated_at: string;
    is_favourite: boolean;
    like_by: number[] | null;
    is_public: boolean;
  };
}
const LOG_PATH = "quananqr1/zusstand/order/order_zustand.ts";

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

export interface OrderState extends BowlOptions {
  //

  setStore: SetState;
  setSetDetails: (id: number, details: Omit<SetState[number], "id">) => void;
  //

  listOfOrders: Order[];

  addToListOfOrders: (order: Order) => void;
  deleteFromListOfOrders: (orderId: number) => void;
  clearListOfOrders: () => void;
  //

  dishState: DishState;
  setDishDetails: (id: number, details: Omit<DishState[number], "id">) => void;
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
  | "addToListOfOrders"
  | "deleteFromListOfOrders"
  | "clearListOfOrders"
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
  setStore: {} as SetState,
  listOfOrders: [],
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
  tableToken: null,
  setDishDetails: function (
    id: number,
    details: Omit<DishState[number], "id">
  ): void {
    throw new Error("Function not implemented.");
  },
  setSetDetails: function (
    id: number,
    details: Omit<SetState[number], "id">
  ): void {
    throw new Error("Function not implemented.");
  }
};

const useOrderStore = create<OrderState>()(
  persist(
    (set, get) => ({
      ...INITIAL_STATE,

      //

      setSetDetails: (id, details) =>
        set((state) => ({
          setStore: {
            ...state.setStore,
            [id]: {
              ...details,
              created_at: details.created_at || new Date().toISOString(),
              updated_at: details.updated_at || new Date().toISOString(),
              userId: details.userId || 0,
              is_favourite: details.is_favourite || false,
              like_by: details.like_by || [],
              is_public: details.is_public || false
            }
          }
        })),
      //
      setDishDetails: (id, details) =>
        set((state) => ({
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

        const dish_items = orderSummary.dishes.map(
          (dish: OrderDetailedDish) => ({
            dish_id: dish.dish_id,
            quantity: dish.quantity
          })
        );

        const set_items = orderSummary.sets.map((set: OrderSetDetailed) => ({
          set_id: set.id,
          quantity: set.quantity
        }));
        logWithLevel(
          {
            orderSummary,
            dish_items: dish_items,
            set_items: set_items
          },
          LOG_PATH,
          "info",
          11
        );

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
          get().clearCurrentOrder();
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

        // Transform DishOrderItem into OrderDetailedDish
        const detailedDishes: OrderDetailedDish[] =
          state.currentOrder.dish_items.map((item) => ({
            dish_id: item.dish_id,
            quantity: item.quantity,
            name: state.dishState[item.dish_id]?.name || "",
            price: state.dishState[item.dish_id]?.price || 0,
            description: state.dishState[item.dish_id]?.description || "",
            iamge: state.dishState[item.dish_id]?.image || "", // Note the "iamge" spelling
            status: state.dishState[item.dish_id]?.status || "active"
          }));

        // Transform SetOrderItem into OrderSetDetailed using setStore
        const detailedSets: OrderSetDetailed[] =
          state.currentOrder.set_items.map((item) => {
            const setDetails = state.setStore[item.set_id] || {};
            return {
              id: item.set_id,
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
                iamge: dish.image || "", // Note the "iamge" spelling to match interface
                status: dish.status || "active"
              })),
              userId: setDetails.userId || 0,
              created_at: setDetails.created_at || new Date().toISOString(),
              updated_at: setDetails.updated_at || new Date().toISOString(),
              is_favourite: setDetails.is_favourite || false,
              like_by: setDetails.like_by || [],
              is_public: setDetails.is_public || false
            };
          });

        // Calculate total price including sets
        const totalPrice =
          detailedDishes.reduce(
            (acc, dish) => acc + dish.price * dish.quantity,
            0
          ) +
          detailedSets.reduce((acc, set) => acc + set.price * set.quantity, 0);

        return {
          totalItems:
            state.currentOrder.dish_items.reduce(
              (acc, item) => acc + item.quantity,
              0
            ) +
            state.currentOrder.set_items.reduce(
              (acc, item) => acc + item.quantity,
              0
            ),
          totalPrice,
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
          },
          // Reset all bowl-related options
          canhKhongRau: 0,
          canhCoRau: 0,
          smallBowl: 0,
          wantChili: false,
          selectedFilling: {
            mocNhi: false,
            thit: false,
            thitMocNhi: false
          }
        }),
      //
      addToListOfOrders: (order) =>
        set((state) => ({
          listOfOrders: [...state.listOfOrders, order]
        })),

      deleteFromListOfOrders: (orderId) =>
        set((state) => ({
          listOfOrders: state.listOfOrders.filter(
            (order) => order.id !== orderId
          )
        })),

      clearListOfOrders: () => set({ listOfOrders: [] }),

      //
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
        selectedFilling: state.selectedFilling,

        setStore: state.setStore // Add this
      })
    }
  )
);

export default useOrderStore;
