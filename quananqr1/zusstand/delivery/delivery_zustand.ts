import { create } from "zustand";
import { persist } from "zustand/middleware";
import { toast } from "@/components/ui/use-toast";
import envConfig from "@/config";

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
    http: any;
    auth: {
      guest: any;
      user: any;
      isGuest: boolean;
    };
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

// Helper functions
function calculateTotalPrice(
  items: DishDeliveryItem[],
  deliveryFee: number
): number {
  console.log(`[${LOG_PATH}] calculateTotalPrice input:`, { items, deliveryFee });
  const itemsTotal = items.reduce(
    (total, item) => total + item.quantity * getPriceForDish(item.dish_id),
    0
  );
  const finalTotal = itemsTotal + deliveryFee;
  console.log(`[${LOG_PATH}] calculateTotalPrice result:`, finalTotal);
  return finalTotal;
}

function getPriceForDish(dishId: number): number {
  console.log(`[${LOG_PATH}] getPriceForDish:`, dishId);
  return 1000;
}

function formatCurrency(amount: number): string {
  console.log(`[${LOG_PATH}] formatCurrency:`, amount);
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
        console.log(`[${LOG_PATH}] updateDeliveryInfo:`, info);
        set((state) => {
          const newState = { ...state, ...info };
          console.log(`[${LOG_PATH}] Updated state:`, newState);
          return newState;
        });
      },

      updateStatus: (status) => {
        console.log(`[${LOG_PATH}] updateStatus:`, status);
        set({ delivery_status: status });
      },

      updateDriverInfo: (driverId, estimatedTime) => {
        console.log(`[${LOG_PATH}] updateDriverInfo:`, { driverId, estimatedTime });
        set({
          driver_id: driverId,
          estimated_delivery_time: estimatedTime,
          delivery_status: "Assigned"
        });
      },

      completeDelivery: (actualTime) => {
        console.log(`[${LOG_PATH}] completeDelivery:`, actualTime);
        set({
          actual_delivery_time: actualTime,
          delivery_status: "Delivered"
        });
      },

      addDishItem: (item) => {
        console.log(`[${LOG_PATH}] addDishItem:`, item);
        set((state) => {
          const newItems = [...state.dish_items, item];
          const newTotal = calculateTotalPrice(newItems, state.delivery_fee);
          const newState = {
            dish_items: newItems,
            total_price: newTotal
          };
          console.log(`[${LOG_PATH}] addDishItem result:`, newState);
          return newState;
        });
      },

      removeDishItem: (dishId) => {
        console.log(`[${LOG_PATH}] removeDishItem:`, dishId);
        set((state) => {
          const updatedItems = state.dish_items.filter(
            (item) => item.dish_id !== dishId
          );
          const newTotal = calculateTotalPrice(updatedItems, state.delivery_fee);
          const newState = {
            dish_items: updatedItems,
            total_price: newTotal
          };
          console.log(`[${LOG_PATH}] removeDishItem result:`, newState);
          return newState;
        });
      },

      updateDishQuantity: (dishId, quantity) => {
        console.log(`[${LOG_PATH}] updateDishQuantity:`, { dishId, quantity });
        set((state) => {
          const updatedItems = state.dish_items.map((item) =>
            item.dish_id === dishId ? { ...item, quantity } : item
          );
          const newTotal = calculateTotalPrice(updatedItems, state.delivery_fee);
          const newState = {
            dish_items: updatedItems,
            total_price: newTotal
          };
          console.log(`[${LOG_PATH}] updateDishQuantity result:`, newState);
          return newState;
        });
      },

      clearDelivery: () => {
        console.log(`[${LOG_PATH}] clearDelivery`);
        set(INITIAL_STATE);
      },

      getFormattedTotal: () => {
        const state = get();
        const formatted = formatCurrency(state.total_price);
        console.log(`[${LOG_PATH}] getFormattedTotal:`, formatted);
        return formatted;
      },

      createDelivery: async ({
        http,
        auth: { guest, user, isGuest },
        orderStore: { tableNumber, getOrderSummary, clearOrder },
        deliveryDetails: {
          deliveryAddress,
          deliveryContact,
          deliveryNotes,
          scheduledTime,
          deliveryFee
        }
      }) => {
        console.log(`[${LOG_PATH}] createDelivery started`, {
          auth: { isGuest },
          tableNumber,
          deliveryDetails: {
            deliveryAddress,
            deliveryContact,
            deliveryNotes,
            scheduledTime,
            deliveryFee
          }
        });

        if (!user && !guest) {
          console.error(`[${LOG_PATH}] Authentication required`);
          toast({
            variant: "destructive",
            title: "Error",
            description: "User authentication required"
          });
          return;
        }

        const orderSummary = getOrderSummary();
        console.log(`[${LOG_PATH}] Order summary:`, orderSummary);

        const currentState = get();
        const dish_items = orderSummary.dishes.map((dish: any) => ({
          dish_id: dish.id,
          quantity: dish.quantity
        }));

        const deliveryData = {
          guest_id: isGuest ? guest?.id ?? null : null,
          user_id: isGuest ? null : user?.id ?? null,
          is_guest: isGuest,
          table_number: tableNumber,
          order_handler_id: 1,
          status: "Pending",
          total_price: orderSummary.totalPrice + deliveryFee,
          dish_items,
          bow_chili: currentState.bow_chili,
          bow_no_chili: currentState.bow_no_chili,
          take_away: false,
          chili_number: currentState.chili_number,
          table_token: currentState.table_token,
          client_name: isGuest ? guest?.name : user?.name,
          delivery_address: deliveryAddress,
          delivery_contact: deliveryContact,
          delivery_notes: deliveryNotes,
          scheduled_time: scheduledTime,
          order_id: orderSummary.orderId,
          delivery_fee: deliveryFee,
          delivery_status: "Pending" as DeliveryStatus
        };

        console.log(`[${LOG_PATH}] Delivery data:`, deliveryData);
        set({ isLoading: true });

        try {
          const deliveryEndpoint = `${envConfig.NEXT_PUBLIC_API_ENDPOINT}${envConfig.Delivery_External_End_Point}`;
          console.log(`[${LOG_PATH}] Making API request to:`, deliveryEndpoint);
          
          const response = await http.post(deliveryEndpoint, deliveryData);
          console.log(`[${LOG_PATH}] API response:`, response.data);
          
          set((state) => ({
            ...state,
            ...deliveryData
          }));

          toast({
            title: "Success",
            description: "Delivery has been created successfully"
          });

          clearOrder();
          return response.data;

        } catch (error) {
          console.error(`[${LOG_PATH}] Delivery creation failed:`, error);
          toast({
            variant: "destructive",
            title: "Error",
            description: error instanceof Error ? error.message : "Failed to create delivery"
          });
          throw error;

        } finally {
          set({ isLoading: false });
          console.log(`[${LOG_PATH}] Delivery creation completed`);
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
        console.log(`[${LOG_PATH}] Persisting state:`, partialState);
        return partialState;
      }
    }
  )
);

export default useDeliveryStore;