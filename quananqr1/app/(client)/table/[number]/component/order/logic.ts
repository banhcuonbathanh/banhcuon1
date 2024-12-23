import { create } from "zustand";
import { toast } from "@/components/ui/use-toast";
import envConfig from "@/config";
import { CreateOrderRequest } from "@/schemaValidations/interface/type_order";
import { logWithLevel } from "@/lib/log";

const LOG_PATH =
  "quananqr1/app/(client)/table/[number]/component/order/logic.ts";

interface OrderCreationState {
  isLoading: boolean;
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
      clearOrder: () => void;
    };
    websocket: {
      disconnect: () => void;
      isConnected: boolean;
      sendMessage: (message: any) => void;
    };
    openLoginDialog: () => void;
  }) => Promise<any>;
}

export const useOrderCreationStore = create<OrderCreationState>((set) => ({
  isLoading: false,

  createOrder: async ({
    topping,
    Table_token,
    http,
    auth: { guest, user, isGuest },
    orderStore: { tableNumber, getOrderSummary, clearOrder },
    websocket: { disconnect, isConnected, sendMessage },
    openLoginDialog
  }) => {
    // Log initial state validation
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
      logWithLevel({ error: "No user or guest found" }, LOG_PATH, "warn", 7);
      openLoginDialog();
      return;
    }

    const orderSummary = getOrderSummary();

    // Log order preparation
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

    // Log request validation
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
      // Log API request attempt
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

      // Log successful order creation
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
        // Log websocket communication
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

      clearOrder();

      return response.data;
    } catch (error) {
      // Log error details
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
  }
}));
