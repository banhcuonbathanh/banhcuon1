"use client";

import React, { useEffect, useState, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { toast } from "@/components/ui/use-toast";
import useCartStore from "@/zusstand/new-order/new-order-zustand";
import { useApiStore } from "@/zusstand/api/api-controller";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
import { useWebSocketStore } from "@/zusstand/web-socket/websocketStore";
import { WebSocketMessage } from "@/schemaValidations/interface/type_websocker";
import { logWithLevel } from "@/lib/logger/log";
import { useRouter } from "next/navigation";
import envConfig from "@/config";
import {
  CreateOrderRequest,
  DishOrderItem,
  Order,
  SetOrderItem
} from "@/schemaValidations/interface/type_order";

const LOG_PATH =
  "quananqr1/app/(client)/table/[number]/component/order/add_order_button.tsx";

const validateOrder = (
  order: Partial<CreateOrderRequest>
): { isValid: boolean; missingFields: string[] } => {
  const requiredFields = [
    "is_guest",
    "table_number",
    "order_handler_id",
    "status",
    "total_price",
    "dish_items",
    "set_items",
    "tracking_order",
    "takeAway",
    "chiliNumber",
    "table_token",
    "order_name"
  ];

  const missingFields = requiredFields.filter((field) => !(field in order));

  if (order.dish_items && order.dish_items.length > 0) {
    order.dish_items.forEach((item: DishOrderItem, index) => {
      ["dish_id", "quantity", "name", "price", "status"].forEach((field) => {
        if (!(field in item)) {
          missingFields.push(`dish_items[${index}].${field}`);
        }
      });
    });
  }

  if (order.set_items && order.set_items.length > 0) {
    order.set_items.forEach((item: SetOrderItem, index) => {
      ["set_id", "quantity", "name", "price", "dishes"].forEach((field) => {
        if (!(field in item)) {
          missingFields.push(`set_items[${index}].${field}`);
        }
      });
    });
  }

  logWithLevel(
    {
      event: "order_validation",
      logId: 4,
      details: {
        isValid: missingFields.length === 0,
        missingFieldsCount: missingFields.length,
        missingFields: missingFields
      }
    },
    LOG_PATH,
    "debug",
    4
  );

  return {
    isValid: missingFields.length === 0,
    missingFields
  };
};

const OrderCreationComponent: React.FC = () => {
  const router = useRouter();
  const [authChecked, setAuthChecked] = useState(false);

  const {
    new_order,
    isLoading,
    tableToken,
    tableNumber,
    current_order,
    setIsLoading,
    addToNewOrder,
    clearCart
  } = useCartStore();
  const { http } = useApiStore();
  const { guest, user, isGuest, openLoginDialog, userId, isLogin } =
    useAuthStore();
  const { connect, disconnect, isConnected, sendMessage, fetchWsToken } =
    useWebSocketStore();

  useEffect(() => {
    const initializeAuth = async () => {
      logWithLevel(
        {
          event: "auth_initialization",
          logId: 1,
          details: {
            startTime: new Date().toISOString()
          }
        },
        LOG_PATH,
        "info",
        1
      );

      const authStore = useAuthStore.getState();
      const unsubscribe = useAuthStore.subscribe((state) => {
        setAuthChecked(true);
        logWithLevel(
          {
            event: "auth_state_updated",
            logId: 1,
            details: {
              isLogin: state.isLogin,
              userId: state.userId,
              isGuest: state.isGuest
            }
          },
          LOG_PATH,
          "info",
          1
        );
        unsubscribe();
      });

      authStore.syncAuthState();
    };

    initializeAuth();
  }, []);

  const handleAddToNewOrders = () => {
    if (!current_order) {
      toast({
        title: "Error",
        description: "No current order to add",
        variant: "destructive"
      });
      return;
    }

    // Create a new order object based on current_order
    const newOrder: Order = {
      ...current_order,
      id: Date.now(), // Generate a temporary ID
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    };

    addToNewOrder(newOrder);

    new_order.length;

    console.log(
      "quananqr1/app/(client)/table/[number]/component/order/add_order_button.tsx",
      new_order.length
    );
    // clearCart();

    toast({
      title: "Success",
      description: "Order added to pending orders"
    });
  };
  useEffect(() => {
    if (authChecked && isLogin && userId) {
      logWithLevel(
        {
          event: "websocket_initialization",
          logId: 2,
          details: {
            userId,
            isGuest,
            authChecked
          }
        },
        LOG_PATH,
        "info",
        2
      );
      initializeWebSocket();
    }
  }, [authChecked, isLogin, userId, user?.email, guest?.name]);

  useEffect(() => {
    return () => {
      logWithLevel(
        {
          event: "component_cleanup",
          logId: 8,
          details: {
            timestamp: new Date().toISOString()
          }
        },
        LOG_PATH,
        "debug",
        8
      );
      disconnect();
    };
  }, [disconnect]);

  const getEmailIdentifier = useCallback(() => {
    logWithLevel(
      {
        event: "user_identification",
        logId: 9,
        details: {
          isGuest,
          guestName: guest?.name,
          userEmail: user?.email
        }
      },
      LOG_PATH,
      "debug",
      9
    );

    return isGuest && guest
      ? guest.name
      : useAuthStore.getState().persistedUser?.email || user?.email;
  }, [isGuest, guest, user?.email]);

  const initializeWebSocket = async () => {
    const emailIdentifier = getEmailIdentifier();

    if (!isLogin || !userId || !emailIdentifier) {
      logWithLevel(
        {
          event: "websocket_initialization_failed",
          logId: 3,
          details: {
            reason: "missing_credentials",
            isLogin,
            userId: userId || null,
            hasEmailIdentifier: !!emailIdentifier
          }
        },
        LOG_PATH,
        "error",
        3
      );
      return;
    }

    try {
      const wstoken = await fetchWsToken({
        userId: Number(userId),
        email: emailIdentifier,
        role: isGuest ? "Guest" : "User"
      });

      if (!wstoken || !tableToken) {
        throw new Error("Failed to obtain required tokens");
      }

      await connect({
        userId: userId.toString(),
        isGuest,
        userToken: wstoken.token,
        tableToken: tableToken,
        role: isGuest ? "Guest" : "User",
        email: emailIdentifier
      });

      logWithLevel(
        {
          event: "websocket_connected",
          logId: 3,
          details: {
            userId,
            isGuest,
            connectionTime: new Date().toISOString()
          }
        },
        LOG_PATH,
        "info",
        3
      );
    } catch (error) {
      logWithLevel(
        {
          event: "websocket_connection_error",
          logId: 3,
          error: error instanceof Error ? error.message : "Unknown error"
        },
        LOG_PATH,
        "error",
        3
      );
      toast({
        title: "Connection Error",
        description:
          "Failed to establish WebSocket connection. Please try again.",
        variant: "destructive"
      });
    }
  };

  const prepareOrder = (): CreateOrderRequest => {
    if (!current_order) {
      logWithLevel(
        {
          event: "order_preparation_failed",
          logId: 10,
          details: {
            reason: "no_current_order"
          }
        },
        LOG_PATH,
        "error",
        10
      );
      throw new Error("No current order found");
    }

    const now = new Date().toISOString();
    const dishItemsTotal =
      current_order.dish_items?.reduce(
        (sum, item) => sum + item.price * item.quantity,
        0
      ) || 0;
    const setItemsTotal =
      current_order.set_items?.reduce(
        (sum, item) => sum + item.price * item.quantity,
        0
      ) || 0;

    const orderData: CreateOrderRequest = {
      created_at: now,
      updated_at: now,
      table_token: tableToken || "",
      table_number: tableNumber || 0,
      is_guest: isGuest,
      guest_id: isGuest ? guest?.id : null,
      user_id: !isGuest ? user?.id : null,
      order_name: isGuest ? guest?.name || "Guest" : user?.name || "User",
      status: "pending",
      order_handler_id: 1,
      total_price: dishItemsTotal + setItemsTotal,
      dish_items: current_order.dish_items || [],
      set_items: current_order.set_items || [],
      topping: current_order.topping || "",
      tracking_order: current_order.tracking_order || "created",
      takeAway: current_order.takeAway || false,
      chiliNumber: current_order.chiliNumber || 0
    };

    logWithLevel(
      {
        event: "order_prepared",
        logId: 10,
        details: {
          totalPrice: orderData,
          dishItemsCount: orderData.dish_items.length,
          setItemsCount: orderData.set_items.length,
          tableNumber: orderData.table_number
        }
      },
      LOG_PATH,
      "info",
      10
    );

    return orderData;
  };

  const sendOrderNotification = async (orderData: CreateOrderRequest) => {
    if (!isConnected) {
      await initializeWebSocket();
    }

    try {
      const messagePayload: WebSocketMessage = {
        type: "order",
        action: "create_message",
        payload: {
          fromUserId: userId?.toString() || "1",
          toUserId: "2",
          type: "order",
          action: "new_order",
          payload: {
            ...orderData
          }
        },
        role: isGuest ? "Guest" : "User"
      };

      await sendMessage(messagePayload);

      logWithLevel(
        {
          event: "notification_sent",
          logId: 7,
          details: {
            orderId: orderData.order_handler_id,
            userId: userId,
            timestamp: new Date().toISOString()
          }
        },
        LOG_PATH,
        "info",
        7
      );
    } catch (error) {
      logWithLevel(
        {
          event: "notification_failed",
          logId: 7,
          error: error instanceof Error ? error.message : "Unknown error",
          details: {
            userId: userId,
            timestamp: new Date().toISOString()
          }
        },
        LOG_PATH,
        "error",
        7
      );
    }
  };

  const handleCreateOrder = async () => {
    if (!isLogin) {
      logWithLevel(
        {
          event: "order_creation_blocked",
          logId: 5,
          details: {
            reason: "not_logged_in"
          }
        },
        LOG_PATH,
        "debug",
        5
      );
      openLoginDialog();
      return;
    }

    if (!isConnected) {
      await initializeWebSocket();
    }

    try {
      setIsLoading(true);
      const orderData = prepareOrder();
      const validation = validateOrder(orderData);
      logWithLevel(
        {
          event: "order_creation_started 12121",
          logId: 5,
          details: {
            isLogin,
            orderData,
            timestamp: new Date().toISOString()
          }
        },
        LOG_PATH,
        "info",
        5
      );

      // if (!validation.isValid) {
      //   throw new Error(
      //     `Invalid order data. Missing fields: ${validation.missingFields.join(
      //       ", "
      //     )}`
      //   );
      // }

      const link_order = `${envConfig.NEXT_PUBLIC_API_ENDPOINT}${envConfig.Order_External_End_Point}`;

      logWithLevel(
        {
          event: "order_creation_response 1212121 33333",
          logId: 6,
          details: {
            success: link_order,
            orderId: orderData,
            timestamp: new Date().toISOString()
          }
        },
        LOG_PATH,
        "info",
        6
      );
      const response = await http.post(link_order, orderData);
      console.log(
        "quananqr1/app/(client)/table/[number]/component/order/add_order_button.tsx response.data.data.id",
        response.data.data.id
      );

      // addToNewOrder(orderData);
      // logWithLevel(
      //   {
      //     event: "order_creation_response",
      //     logId: 6,
      //     details: {
      //       success: !!response.data,
      //       orderId: response.data?.id,
      //       timestamp: new Date().toISOString()
      //     }
      //   },
      //   LOG_PATH,
      //   "info",
      //   6
      // );

      if (!response.data) {
        throw new Error("Failed to create order");
      }

      await sendOrderNotification(orderData);

      toast({
        title: "Order Created",
        description: "Your order has been successfully created!"
      });

      router.refresh();
    } catch (error) {
      logWithLevel(
        {
          event: "order_creation_error",
          logId: 11,
          error: error instanceof Error ? error.message : "Unknown error",
          details: {
            timestamp: new Date().toISOString()
          }
        },
        LOG_PATH,
        "error",
        6
      );

      toast({
        title: "Order Creation Failed",
        description:
          error instanceof Error ? error.message : "Failed to create order",
        variant: "destructive"
      });
    } finally {
      setIsLoading(false);
    }
  };

  const getButtonText = () => {
    if (!authChecked) return "Loading...";
    if (!isLogin) return "Login to Order";
    if (
      !current_order?.dish_items?.length &&
      !current_order?.set_items?.length
    ) {
      return "Add Items to Order";
    }
    return "Create Order";
  };

  const isButtonDisabled = () => {
    if (!authChecked) return true;
    if (!isLogin) return false;
    if (
      !current_order?.dish_items?.length &&
      !current_order?.set_items?.length
    ) {
      return true;
    }
    return isLoading;
  };

  const getValidationMessage = () => {
    if (!tableNumber || !tableToken) {
      return "Table information is missing";
    }
    if (!current_order) {
      return "No order data found";
    }
    if (!current_order.dish_items?.length && !current_order.set_items?.length) {
      return "Please add at least one item to your order";
    }
    if (isLoading) {
      return "Creating your order...";
    }
    if (!isLogin) {
      return "Please login to create an order";
    }
    return null;
  };

  const validationMessage = getValidationMessage();
  const handleCombinedActions = () => {
    handleCreateOrder();
    handleAddToNewOrders();
  };

  return (
    <div className="mt-4 space-y-2">
      {validationMessage && (
        <div className="text-sm text-center px-4 py-2 rounded-md bg-yellow-50 text-yellow-800 border border-yellow-200">
          {validationMessage}
        </div>
      )}
      <Button
        className="w-full"
        onClick={handleCombinedActions}
        disabled={isButtonDisabled()}
      >
        {getButtonText()}
      </Button>

      <Button
        variant="destructive"
        className="w-full mt-2"
        onClick={() => {
          const store = useCartStore.getState();
          store.clearCart();
          store.clearNewOrders();
          toast({
            title: "Orders Cleared",
            description: "All orders have been cleared successfully."
          });
        }}
      >
        Clear All Orders
      </Button>

      <Button
        variant="outline"
        className="w-full"
        onClick={handleAddToNewOrders}
        disabled={!current_order || isLoading}
      >
        Add to Pending Orders
      </Button>
    </div>
  );
};

export default OrderCreationComponent;
