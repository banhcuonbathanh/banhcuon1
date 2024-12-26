"use client";

import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import useOrderStore from "@/zusstand/order/order_zustand";
// import { useOrderCreationStore } from "./logic";
import { useApiStore } from "@/zusstand/api/api-controller";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
import { useWebSocketStore } from "@/zusstand/web-socket/websocketStore";
import { WebSocketMessage } from "@/schemaValidations/interface/type_websocker";
import { logWithLevel } from "@/lib/logger/log";
import { useRouter } from "next/navigation";

const LOG_PATH =
  "quananqr1/app/(client)/table/[number]/component/order/add_order_button.tsx";

interface OrderCreationComponentProps {
  table_token: string;
  table_number: string;
}

const OrderCreationComponent: React.FC<OrderCreationComponentProps> = ({
  table_number,
  table_token
}) => {
  const router = useRouter();
  logWithLevel(
    { event: "component_init", table_number, table_token },
    LOG_PATH,
    "debug",
    1
  );

  const {
    getOrderSummary,
    createOrder,
    canhKhongRau,
    canhCoRau,
    smallBowl,
    wantChili,
    selectedFilling,
    isLoading,

    //

    addToListOfOrders

    //
  } = useOrderStore();
  const { http } = useApiStore();
  const { guest, user, isGuest, openLoginDialog, userId, isLogin } =
    useAuthStore();
  const {
    connect,
    disconnect,
    isConnected,
    sendMessage,
    wsToken,
    fetchWsToken
  } = useWebSocketStore();

  const [authChecked, setAuthChecked] = useState(false);

  const getFillingString = (filling: {
    mocNhi: boolean;
    thit: boolean;
    thitMocNhi: boolean;
  }) => {
    if (filling.mocNhi) return "Mọc Nhĩ";
    if (filling.thit) return "Thịt";
    if (filling.thitMocNhi) return "Thịt Mọc Nhĩ";
    return "Không";
  };

  let topping = `canhKhongRau ${canhKhongRau} - canhCoRau ${canhCoRau} - bat be ${smallBowl} - ot tuoi ${wantChili} - nhan ${getFillingString(
    selectedFilling
  )} -`;

  const orderSummary = getOrderSummary();

  useEffect(() => {
    const initializeAuth = async () => {
      const authStore = useAuthStore.getState();

      // Set up a one-time subscription to catch the state update
      const unsubscribe = useAuthStore.subscribe((state) => {
        setAuthChecked(true);
        logWithLevel(
          {
            event: "auth_initialized",
            isLogin: state.isLogin,
            userId: state.userId,
            isGuest: state.isGuest
          },
          LOG_PATH,
          "info",
          5
        );

        // Unsubscribe after first update
        unsubscribe();
      });

      // Trigger the state sync
      authStore.syncAuthState();
    };

    initializeAuth();
  }, []);

  useEffect(() => {
    if (authChecked && isLogin && userId) {
      logWithLevel(
        { event: "websocket_init_attempt", userId, isGuest },
        LOG_PATH,
        "debug",
        2
      );
      initializeWebSocket();
    }
  }, [authChecked, isLogin, userId, user?.email, guest?.name]);

  const getEmailIdentifier = () => {
    const result =
      isGuest && guest
        ? guest.name
        : useAuthStore.getState().persistedUser?.email || user?.email;

    logWithLevel(
      { event: "email_identifier_lookup", isGuest, result },
      LOG_PATH,
      "debug",
      9
    );

    return result;
  };

  const initializeWebSocket = async () => {
    const emailIdentifier = getEmailIdentifier();

    logWithLevel(
      { event: "websocket_init", userId, emailIdentifier },
      LOG_PATH,
      "info",
      2
    );

    if (!isLogin || !userId || !emailIdentifier) {
      logWithLevel(
        { event: "websocket_init_failed", reason: "missing_credentials" },
        LOG_PATH,
        "error",
        2
      );
      return;
    }

    try {
      const wstoken1 = await fetchWsToken({
        userId: Number(userId),
        email: emailIdentifier,
        role: isGuest ? "Guest" : "User"
      });

      logWithLevel(
        { event: "ws_token_fetch", success: !!wstoken1 },
        LOG_PATH,
        "debug",
        7
      );

      if (!wstoken1) {
        throw new Error("Failed to obtain WebSocket token");
      }

      await connect({
        userId: userId.toString(),
        isGuest,
        userToken: wstoken1.token,
        tableToken: table_token,
        role: isGuest ? "Guest" : "User",
        email: emailIdentifier
      });

      logWithLevel(
        { event: "websocket_connected", userId },
        LOG_PATH,
        "info",
        2
      );
    } catch (error) {
      logWithLevel(
        { event: "websocket_error", error: "error.message" },
        LOG_PATH,
        "error",
        2
      );
    }
  };

  const handleCreateOrder = async () => {
    useAuthStore.getState().syncAuthState();
    const currentAuthState = useAuthStore.getState();

    logWithLevel(
      { event: "create_order_attempt", isLogin: currentAuthState.isLogin },
      LOG_PATH,
      "info",
      3
    );

    if (!currentAuthState.isLogin) {
      logWithLevel(
        { event: "order_creation_blocked", reason: "not_logged_in" },
        LOG_PATH,
        "warn",
        3
      );
      openLoginDialog();
      return;
    }

    if (!isConnected) {
      logWithLevel({ event: "reconnection_attempt" }, LOG_PATH, "debug", 2);
      await initializeWebSocket();
    }

    if (orderSummary.totalItems === 0) {
      logWithLevel(
        { event: "order_validation_failed", reason: "no_items" },
        LOG_PATH,
        "warn",
        6
      );
      return;
    }

    try {
      const order = await createOrder({
        topping,
        Table_token: table_token,
        http,
        auth: { guest, user, isGuest },
        orderStore: {
          tableNumber: Number(table_number),
          getOrderSummary
        },
        websocket: { disconnect, isConnected, sendMessage },
        openLoginDialog
      });
      logWithLevel(
        { event: "order_created 111", order: order.data },
        LOG_PATH,
        "info",
        3
      );
      addToListOfOrders(order.data);

      //------------------

      router.refresh();
      //----------

      await sendMessage1();
    } catch (error) {
      logWithLevel(
        { event: "order_creation_failed", error: "error.message" },
        LOG_PATH,
        "error",
        3
      );
    }
  };

  const sendMessage1 = async () => {
    logWithLevel({ event: "send_message_attempt" }, LOG_PATH, "debug", 4);

    useAuthStore.getState().syncAuthState();
    const currentAuthState = useAuthStore.getState();

    if (!currentAuthState.isLogin) {
      logWithLevel(
        { event: "message_send_blocked", reason: "not_logged_in" },
        LOG_PATH,
        "warn",
        4
      );
      openLoginDialog();
      return;
    }

    if (!isConnected) {
      logWithLevel({ event: "ws_reconnection_attempt" }, LOG_PATH, "debug", 2);
      await initializeWebSocket();
      if (!isConnected) {
        logWithLevel(
          { event: "message_send_failed", reason: "connection_failed" },
          LOG_PATH,
          "error",
          4
        );
        return;
      }
    }

    try {
      const { dishes, sets, totalPrice } = getOrderSummary();

      const messagePayload: WebSocketMessage = {
        type: "order",
        action: "create_message",
        payload: {
          fromUserId: userId?.toString() || "1",
          toUserId: "2",
          type: "order",
          action: "new_order",
          payload: {
            guest_id: isGuest ? guest?.id || null : null,
            user_id: !isGuest ? user?.id || 1 : null,
            is_guest: isGuest,
            table_number: Number(table_number),
            order_handler_id: 1,
            status: "pending",
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            total_price: totalPrice,
            order_name: isGuest ? guest?.name || "Guest" : user?.name || "User",
            dish_items: dishes.map((dish) => ({
              dish_id: dish.dish_id,
              quantity: dish.quantity
            })),
            set_items: sets.map((set) => ({
              set_id: set.id,
              quantity: set.quantity
            })),
            bow_chili: wantChili ? 1 : 0,
            bow_no_chili: canhKhongRau,
            take_away: false,
            chili_number: wantChili ? 1 : 0,
            Table_token: table_token
          }
        },
        role: isGuest ? "Guest" : "User",
        roomId: table_number.toString()
      };

      sendMessage(messagePayload);
      logWithLevel(
        { event: "message_sent_success", messageType: messagePayload.type },
        LOG_PATH,
        "info",
        4
      );
    } catch (error) {
      logWithLevel(
        { event: "message_send_error", error: "error.message" },
        LOG_PATH,
        "error",
        4
      );
    }
  };

  const getButtonText = () => {
    if (!authChecked) return "Loading...";
    if (!isLogin) return "Login to Order";
    if (orderSummary.totalItems === 0) return "Add Items to Order";
    return "Create Order";
  };

  const isButtonDisabled = () => {
    if (!authChecked) return true;
    if (!isLogin) return false;
    if (orderSummary.totalItems === 0) return true;
    return isLoading;
  };

  useEffect(() => {
    return () => {
      logWithLevel({ event: "component_cleanup" }, LOG_PATH, "debug", 8);
      disconnect();
    };
  }, []);

  return (
    <div className="mt-4">
      <Button
        className="w-full"
        onClick={handleCreateOrder}
        disabled={isButtonDisabled()}
      >
        {getButtonText()}
      </Button>
    </div>
  );
};

export default OrderCreationComponent;
