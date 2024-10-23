import React, { useState, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { toast } from "@/components/ui/use-toast";
import { CreateOrderRequest } from "@/schemaValidations/interface/type_order";
import { useApiStore } from "@/zusstand/api/api-controller";
import useOrderStore from "@/zusstand/order/order_zustand";

interface OrderCreationComponentProps {
  bowlChili: number;
  bowlNoChili: number;
}

const OrderCreationComponent: React.FC<OrderCreationComponentProps> = ({
  bowlChili,
  bowlNoChili
}) => {
  const { http } = useApiStore();
  const [isLoading, setIsLoading] = useState(false);

  const { tableNumber, tabletoken, getOrderSummary, clearOrder } =
    useOrderStore();

  const createOrder = useCallback(
    (orderData: CreateOrderRequest) => {
      return new Promise((resolve, reject) => {
        http
          .post("http://localhost:8888/orders", orderData)
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

    const orderData: CreateOrderRequest = {
      guest_id: null,
      user_id: 1,
      is_guest: false,
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

      // asokdjfjasdlfjlasdkjflkasjdlf;jasdlkfjlasdjfl;asjdlf;jsa;l
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
    bowlNoChili
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
