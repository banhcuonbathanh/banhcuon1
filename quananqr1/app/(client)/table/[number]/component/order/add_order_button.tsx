"use client";

import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import useOrderStore from "@/zusstand/order/order_zustand";
import { useOrderCreationStore } from "./logic";
import { useApiStore } from "@/zusstand/api/api-controller";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
import { useWebSocketStore } from "@/zusstand/web-socket/websocketStore";
import { WebSocketMessage } from "@/schemaValidations/interface/type_websocker";
import { logWithLevel } from "@/lib/log";

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
  // Log component initialization
  logWithLevel(
    { event: "component_init", table_number, table_token },
    LOG_PATH,
    "debug",
    1
  );

  const { isLoading, createOrder } = useOrderCreationStore();
  const {
    getOrderSummary,
    clearOrder,
    canhKhongRau,
    canhCoRau,
    smallBowl,
    wantChili,
    selectedFilling
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
      useAuthStore.getState().syncAuthState();
      setAuthChecked(true);

      logWithLevel(
        { event: "auth_initialized", isLogin, userId },
        LOG_PATH,
        "info",
        5
      );
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
          getOrderSummary,
          clearOrder
        },
        websocket: { disconnect, isConnected, sendMessage },
        openLoginDialog
      });

      logWithLevel(
        { event: "order_created", orderId: order?.id },
        LOG_PATH,
        "info",
        3
      );

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
          fromUserId: "1",
          toUserId: "2",
          type: "order",
          action: "new_order",
          payload: {
            guest_id: null,
            user_id: 1,
            is_guest: false,
            table_number: 1,
            order_handler_id: 1,
            status: "pending",
            created_at: "2024-10-21T12:00:00Z",
            updated_at: "2024-10-21T12:00:00Z",
            total_price: 5000,
            order_name: "test",
            dish_items: [
              { dish_id: 1, quantity: 2 },
              { dish_id: 2, quantity: 2 },
              { dish_id: 3, quantity: 4 }
            ],
            set_items: [
              { set_id: 1, quantity: 3 },
              { set_id: 2, quantity: 3 }
            ],
            bow_chili: 1,
            bow_no_chili: 2,
            take_away: true,
            chili_number: 3,
            Table_token: "MTp0YWJsZTo0ODgzMjc3NDQy.2AZhkuCtKB0"
          }
        },
        role: "User",
        roomId: "1"
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

  return (
    <div className="mt-4">
      <Button
        className="w-full"
        onClick={handleCreateOrder}
        disabled={isButtonDisabled()}
      >
        {getButtonText()}
      </Button>

      <Button
        className="w-full"
        onClick={sendMessage1}
        disabled={isButtonDisabled()}
      >
        {"send message"}
      </Button>
    </div>
  );
};

export default OrderCreationComponent;
