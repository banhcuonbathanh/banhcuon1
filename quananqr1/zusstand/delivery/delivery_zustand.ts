// delivery-store.ts
import { create } from "zustand";
import { persist } from "zustand/middleware";
import { toast } from "@/components/ui/use-toast";
import envConfig from "@/config";
import { logWithLevel } from "@/lib/log";
import { useApiStore } from "../api/api-controller";

const LOG_PATH = "quananqr1/zusstand/delivery/delivery_zustand.ts";

export const DeliveryStatusValues = {
  Pending: "Pending",
  Assigned: "Assigned",
  PickedUp: "Picked Up",
  InTransit: "In Transit",
  Delivered: "Delivered",
  Failed: "Failed",
  Cancelled: "Cancelled"
} as const;

export type DeliveryStatus =
  (typeof DeliveryStatusValues)[keyof typeof DeliveryStatusValues];

interface DishDeliveryItem {
  dish_id: number;
  quantity: number;
}

interface DeliveryState {
  isLoading: boolean;
  guest_id: number | null;
  user_id: number;
  is_guest: boolean;
  table_number: number;
  order_handler_id: number;
  status: string;
  total_price: number;
  dish_items: DishDeliveryItem[];
  bow_chili: number;
  bow_no_chili: number;
  take_away: boolean;
  chili_number: number;
  table_token: string;
  client_name: string;
  delivery_address: string;
  delivery_contact: string;
  delivery_notes: string;
  scheduled_time: string;
  order_id: number;
  delivery_fee: number;
  delivery_status: DeliveryStatus;
  driver_id?: number;
  estimated_delivery_time?: string;
  actual_delivery_time?: string;
}

interface DeliveryActions {
  updateDeliveryInfo: (info: Partial<DeliveryState>) => void;
  updateStatus: (status: DeliveryStatus) => void;
  updateDriverInfo: (driverId: number, estimatedTime: string) => void;
  completeDelivery: (actualTime: string) => void;
  addDishItem: (item: DishDeliveryItem) => void;
  removeDishItem: (dishId: number) => void;
  updateDishQuantity: (dishId: number, quantity: number) => void;
  clearDelivery: () => void;
  getFormattedTotal: () => string;
  createDelivery: (params: {
    guest: any;
    user: any;
    isGuest: boolean;
    orderStore: {
      tableNumber: number;
      getOrderSummary: () => any;
      clearOrder: () => void;
    };
    deliveryDetails: {
      deliveryAddress: string;
      deliveryContact: string;
      deliveryNotes: string;
      scheduledTime: string;
      deliveryFee: number;
    };
  }) => Promise<any>;
}

const INITIAL_STATE: DeliveryState = {
  isLoading: false,
  guest_id: null,
  user_id: 0,
  is_guest: false,
  table_number: 0,
  order_handler_id: 0,
  status: "Pending",
  total_price: 0,
  dish_items: [],
  bow_chili: 0,
  bow_no_chili: 0,
  take_away: false,
  chili_number: 0,
  table_token: "",
  client_name: "",
  delivery_address: "",
  delivery_contact: "",
  delivery_notes: "",
  scheduled_time: "",
  order_id: 0,
  delivery_fee: 0,
  delivery_status: "Pending"
};

function getPriceForDish(dishId: number): number {
  logWithLevel({ dishId }, LOG_PATH, "info", 4);
  return 1000; // Replace with actual price logic
}

function calculateTotalPrice(
  items: DishDeliveryItem[],
  deliveryFee: number
): number {
  logWithLevel({ items, deliveryFee }, LOG_PATH, "info", 4);
  const itemsTotal = items.reduce(
    (total, item) => total + item.quantity * getPriceForDish(item.dish_id),
    0
  );
  const finalTotal = itemsTotal + deliveryFee;
  return finalTotal;
}

function formatCurrency(amount: number): string {
  return amount.toLocaleString("en-US", {
    style: "currency",
    currency: "USD"
  });
}

const useDeliveryStore = create<DeliveryState & DeliveryActions>()(
  persist(
    (set, get) => ({
      ...INITIAL_STATE,

      updateDeliveryInfo: (info) => {
        logWithLevel({ info }, LOG_PATH, "info", 8);
        set((state) => {
          const newState = { ...state, ...info };
          return newState;
        });
      },

      updateStatus: (status) => {
        logWithLevel({ status }, LOG_PATH, "info", 5);
        set({ delivery_status: status });
      },

      updateDriverInfo: (driverId, estimatedTime) => {
        logWithLevel({ driverId, estimatedTime }, LOG_PATH, "info", 5);
        set({
          driver_id: driverId,
          estimated_delivery_time: estimatedTime,
          delivery_status: "Assigned"
        });
      },

      completeDelivery: (actualTime) => {
        logWithLevel({ actualTime }, LOG_PATH, "info", 5);
        set({
          actual_delivery_time: actualTime,
          delivery_status: "Delivered"
        });
      },

      addDishItem: (item) => {
        logWithLevel({ action: "addDishItem", item }, LOG_PATH, "info", 3);
        set((state) => {
          const newItems = [...state.dish_items, item];
          const newTotal = calculateTotalPrice(newItems, state.delivery_fee);
          return {
            dish_items: newItems,
            total_price: newTotal
          };
        });
      },

      removeDishItem: (dishId) => {
        logWithLevel({ action: "removeDishItem", dishId }, LOG_PATH, "info", 3);
        set((state) => {
          const updatedItems = state.dish_items.filter(
            (item) => item.dish_id !== dishId
          );
          const newTotal = calculateTotalPrice(
            updatedItems,
            state.delivery_fee
          );
          return {
            dish_items: updatedItems,
            total_price: newTotal
          };
        });
      },

      updateDishQuantity: (dishId, quantity) => {
        logWithLevel(
          { action: "updateDishQuantity", dishId, quantity },
          LOG_PATH,
          "info",
          3
        );
        set((state) => {
          const updatedItems = state.dish_items.map((item) =>
            item.dish_id === dishId ? { ...item, quantity } : item
          );
          const newTotal = calculateTotalPrice(
            updatedItems,
            state.delivery_fee
          );
          return {
            dish_items: updatedItems,
            total_price: newTotal
          };
        });
      },

      clearDelivery: () => {
        logWithLevel({ action: "clearDelivery" }, LOG_PATH, "info", 8);
        set(INITIAL_STATE);
      },

      getFormattedTotal: () => {
        const state = get();
        return formatCurrency(state.total_price);
      },

      createDelivery: async ({
        guest,
        user,
        isGuest,
        orderStore,
        deliveryDetails
      }) => {
        // Validate required fields
        if (!orderStore?.getOrderSummary) {
          throw new Error("Order summary function is required");
        }

        if (!orderStore?.tableNumber) {
          throw new Error("Table number is required");
        }

        // Get order summary first
        const orderSummary = orderStore.getOrderSummary();
        logWithLevel({ orderSummary }, LOG_PATH, "info", 9);

        if (!orderSummary?.dishes?.length) {
          throw new Error("No dishes found in order");
        }

        // Validate user/guest data
        if (isGuest && (!guest || !guest.id)) {
          throw new Error("Guest ID is required for guest orders");
        }

        if (!isGuest && (!user || !user.id)) {
          throw new Error("User ID is required for user orders");
        }

        const dish_items = orderSummary.dishes.map((dish: any) => ({
          dish_id: dish.id,
          quantity: dish.quantity
        }));

        const deliveryData = {
          guest_id: isGuest ? guest.id : null,
          user_id: isGuest ? null : user.id,
          is_guest: isGuest,
          table_number: orderStore.tableNumber,
          order_handler_id: 1,
          status: "Pending",
          total_price: calculateTotalPrice(
            dish_items,
            deliveryDetails.deliveryFee
          ),
          dish_items,
          bow_chili: 0,
          bow_no_chili: 0,
          take_away: false,
          chili_number: 0,
          table_token: get().table_token,
          client_name: isGuest ? guest?.name : user?.name,
          delivery_address: deliveryDetails.deliveryAddress,
          delivery_contact: deliveryDetails.deliveryContact,
          delivery_notes: deliveryDetails.deliveryNotes,
          scheduled_time: deliveryDetails.scheduledTime,
          order_id: orderSummary.orderId,
          delivery_fee: deliveryDetails.deliveryFee,
          delivery_status: "Pending" as DeliveryStatus
        };

        logWithLevel({ deliveryData }, LOG_PATH, "info", 1);

        set({ isLoading: true });

        try {
          const deliveryEndpoint = `${envConfig.NEXT_PUBLIC_API_ENDPOINT}${envConfig.Delivery_External_End_Point}`;
          logWithLevel(
            { endpoint: deliveryEndpoint, method: "POST", data: deliveryData },
            LOG_PATH,
            "info",
            2
          );

          const response = await useApiStore
            .getState()
            .http.post(deliveryEndpoint, deliveryData);

          set((state) => ({
            ...state,
            ...deliveryData
          }));

          toast({
            title: "Success",
            description: "Delivery has been created successfully"
          });

          orderStore.clearOrder();
          return response.data;
        } catch (error) {
          logWithLevel({ error }, LOG_PATH, "error", 7);
          toast({
            variant: "destructive",
            title: "Error",
            description:
              error instanceof Error
                ? error.message
                : "Failed to create delivery"
          });
          throw error;
        } finally {
          set({ isLoading: false });
        }
      }
    }),
    {
      name: "delivery-storage",
      partialize: (state) => {
        const partialState = {
          dish_items: state.dish_items,
          delivery_status: state.delivery_status,
          total_price: state.total_price,
          client_name: state.client_name,
          delivery_address: state.delivery_address,
          delivery_contact: state.delivery_contact,
          delivery_notes: state.delivery_notes,
          scheduled_time: state.scheduled_time,
          order_id: state.order_id
        };
        logWithLevel({ persistedState: partialState }, LOG_PATH, "info", 6);
        return partialState;
      }
    }
  )
);

export default useDeliveryStore;
