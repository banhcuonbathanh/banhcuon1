import React, { useState, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { toast } from "@/components/ui/use-toast";
import { CreateOrderRequest } from "@/schemaValidations/interface/type_order";
import { useApiStore } from "@/zusstand/api/api-controller";
import useOrderStore from "@/zusstand/order/order_zustand";
import envConfig from "@/config";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
interface OrderCreationComponentProps {
  bowlChili: number;
  bowlNoChili: number;
}

const OrderCreationComponent: React.FC<OrderCreationComponentProps> = ({
  bowlChili,
  bowlNoChili
}) => {
  const link_order = `${envConfig.NEXT_PUBLIC_API_ENDPOINT}${envConfig.Order_External_End_Point}`;
  const { http } = useApiStore();
  const [isLoading, setIsLoading] = useState(false);

  const { tableNumber, tabletoken, getOrderSummary, clearOrder } =
    useOrderStore();
  const { guest, user, isGuest } = useAuthStore();

  const createOrder = useCallback(
    (orderData: CreateOrderRequest) => {
      return new Promise((resolve, reject) => {
        http
          .post(link_order, orderData)
          .then((response) => {
            resolve(response.data);
          })
          .catch((error) => {
            reject(
              new Error(
                error.response?.data?.message || "Failed to create order"
              )
            );
          });
      });
    },
    [http]
  );

  const handleCreateOrder = useCallback(() => {
    if (isLoading) return;

    const orderSummary = getOrderSummary();

    const dish_items = orderSummary.dishes.map((dish) => ({
      dish_id: dish.id,
      quantity: dish.quantity
    }));

    const set_items = orderSummary.sets.map((set) => ({
      set_id: set.id,
      quantity: set.quantity
    }));

    // Determine user_id and guest_id based on authentication state
    let user_id: number | null = null;
    let guest_id: number | null = null;

    if (isGuest) {
      // User is logged in
      user_id = user?.id ?? null;
      guest_id = null;
    } else {
      // Guest is logged in
      user_id = null;
      guest_id = guest?.id ?? null;
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
      bow_chili: bowlChili,
      bow_no_chili: bowlNoChili,
      takeAway: false,
      chiliNumber: 0
    };

    setIsLoading(true);

    createOrder(orderData)
      .then((result) => {
        toast({
          title: "Success",
          description: "Order has been created successfully"
        });
        clearOrder();
      })
      .catch((error) => {
        console.error("Order creation failed:", error);
        toast({
          variant: "destructive",
          title: "Error",
          description:
            error instanceof Error ? error.message : "Failed to create order"
        });
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, [
    createOrder,
    isLoading,
    tableNumber,
    getOrderSummary,
    clearOrder,
    bowlChili,
    bowlNoChili,
    user,
    guest,
    isGuest
  ]);

  const orderSummary = getOrderSummary();
  const isDisabled = isLoading || !tableNumber || orderSummary.totalItems === 0;

  return (
    <div className="mt-4">
      <Button
        className="w-full"
        onClick={handleCreateOrder}
        disabled={isDisabled}
      >
        {isLoading ? "Creating Order..." : "Place Order"}
      </Button>
    </div>
  );
};

export default OrderCreationComponent;
