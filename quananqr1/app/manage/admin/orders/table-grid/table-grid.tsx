import React, { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog";
import {
  RestaurantSummaryProps,
  AggregatedDish
} from "../restaurant-summary/restaurant-summary/rs-type";
import { OrderDetailedResponse } from "../component/new-order-column";

const aggregateDishes = (orders: OrderDetailedResponse[]): AggregatedDish[] => {
  const dishMap = new Map<number, AggregatedDish>();

  orders.forEach((order) => {
    // Add individual dishes
    order.data_dish.forEach((dish) => {
      const existingDish = dishMap.get(dish.dish_id);
      if (existingDish) {
        existingDish.quantity += dish.quantity;
      } else {
        dishMap.set(dish.dish_id, { ...dish, quantity: dish.quantity });
      }
    });

    // Add dishes from sets
    order.data_set.forEach((set) => {
      set.dishes.forEach((setDish) => {
        const existingDish = dishMap.get(setDish.dish_id);
        if (existingDish) {
          existingDish.quantity += setDish.quantity * set.quantity;
        } else {
          dishMap.set(setDish.dish_id, {
            ...setDish,
            quantity: setDish.quantity * set.quantity
          });
        }
      });
    });
  });

  return Array.from(dishMap.values());
};

// Helper function to get all items from an order
const getAllItems = (order: OrderDetailedResponse) => {
  const setItems = order.data_set.map((set) => ({
    name: set.name,
    quantity: set.quantity,
    price: set.price
  }));

  const dishItems = order.data_dish.map((dish) => ({
    name: dish.name,
    quantity: dish.quantity,
    price: dish.price
  }));

  return [...setItems, ...dishItems];
};

const OrderSummaryModal: React.FC<{
  isOpen: boolean;
  onClose: () => void;
  order: OrderDetailedResponse | null;
}> = ({ isOpen, onClose, order }) => {
  if (!order) return null;

  const aggregatedDishesData = aggregateDishes([order]);
  const total = aggregatedDishesData.reduce(
    (sum, dish) => sum + dish.price * dish.quantity,
    0
  );

  return (
    <Dialog open={isOpen} onOpenChange={() => onClose()}>
      <DialogContent className="max-w-md bg-slate-500 shadow-lg border">
        <DialogHeader>
          <DialogTitle>Order Summary - Table {order.table_number}</DialogTitle>
        </DialogHeader>
        <ul className="list-disc pl-4 space-y-1">
          {getAllItems(order).map((item, index) => (
            <li key={index}>
              {item.name} <span className="font-medium">x{item.quantity}</span>
              <span className="text-gray-600 text-sm ml-2">
                (${item.price.toFixed(2)})
              </span>
            </li>
          ))}
        </ul>
      </DialogContent>
    </Dialog>
  );
};

export const TableGrid: React.FC<RestaurantSummaryProps> = ({
  restaurantLayoutProps
}) => {
  const [selectedOrder, setSelectedOrder] =
    useState<OrderDetailedResponse | null>(null);
  const [tappedOrderId, setTappedOrderId] = useState<number | null>(null);
  const headers = ["Order", "Delivery", "Topping", "Total"];

  const calculateTotal = (order: OrderDetailedResponse) => {
    let setTotal = order.data_set.reduce(
      (acc, set) => acc + set.price * set.quantity,
      0
    );
    let dishTotal = order.data_dish.reduce(
      (acc, dish) => acc + dish.price * dish.quantity,
      0
    );
    return setTotal + dishTotal;
  };

  const getTotalItems = (order: OrderDetailedResponse) => {
    const setQuantities = order.data_set.reduce(
      (acc, set) => acc + set.quantity,
      0
    );
    const dishQuantities = order.data_dish.reduce(
      (acc, dish) => acc + dish.quantity,
      0
    );
    return setQuantities + dishQuantities;
  };

  const getOrderTypeLabel = (takeAway: boolean) => {
    return takeAway ? (
      <span className="text-orange-600 text-sm">Takeaway</span>
    ) : (
      <span className="text-blue-600 text-sm">Dine-in</span>
    );
  };

  const getStatusBadge = (status: string) => {
    const statusColors: Record<string, string> = {
      pending: "bg-yellow-100 text-yellow-800",
      processing: "bg-blue-100 text-blue-800",
      completed: "bg-green-100 text-green-800",
      cancelled: "bg-red-100 text-red-800"
    };

    const colorClass =
      statusColors[status.toLowerCase()] || "bg-gray-100 text-gray-800";
    return (
      <span
        className={`px-2 py-1 rounded-full text-xs font-medium ${colorClass}`}
      >
        {status}
      </span>
    );
  };

  return (
    <div className="w-full max-w-5xl mx-auto">
      <div className="grid grid-cols-5">
        <div className="border font-bold p-2"></div>
        {headers.map((header) => (
          <div key={header} className="border font-bold p-2 text-center">
            {header}
          </div>
        ))}

        {restaurantLayoutProps.map((order) => (
          <React.Fragment key={order.id}>
            <div className="border p-2 space-y-2">
              <div className="font-bold">Table {order.table_number}</div>
              <div className="text-sm">
                {order.order_name || "Unnamed Order"}
              </div>
              <div>{getOrderTypeLabel(order.takeAway)}</div>
              <div>{getStatusBadge(order.status)}</div>
            </div>

            <div
              className={`border p-2 cursor-pointer transition-colors duration-200 ${
                tappedOrderId === order.id ? "bg-gray-50" : ""
              }`}
              onClick={() => {
                setSelectedOrder(order);
                setTappedOrderId(order.id);
              }}
            >
              <div className="font-medium mb-2">
                Order Summary - Table {order.table_number}
              </div>
              <div className="space-y-4 p-4">
                <div className="space-y-2">
                  {aggregateDishes([order]).map((dish, index) => (
                    <div
                      key={index}
                      className="flex justify-between items-center border-b pb-2"
                    >
                      <div>
                        <div className="font-medium">{dish.name}</div>
                        <div className="text-sm text-gray-500">
                          ${dish.price.toFixed(2)} × {dish.quantity}
                        </div>
                      </div>
                      <div className="font-bold">
                        ${(dish.price * dish.quantity).toFixed(2)}
                      </div>
                    </div>
                  ))}
                </div>
                <div className="flex justify-between items-center pt-4 border-t">
                  <div className="font-bold text-lg">Total</div>
                  <div className="font-bold text-lg text-green-600">
                    $
                    {aggregateDishes([order])
                      .reduce(
                        (sum, dish) => sum + dish.price * dish.quantity,
                        0
                      )
                      .toFixed(2)}
                  </div>
                </div>
              </div>
            </div>

            <div className="border p-2">
              <div className="text-center">
                {order.tracking_order || "Processing"}
              </div>
            </div>

            <div className="border p-2">
              <div className="text-center">
                {order.topping || "No extra toppings"}
              </div>
            </div>

            <div className="border p-2">
              <div className="text-center font-bold text-lg text-green-600">
                ${calculateTotal(order).toFixed(2)}
              </div>
              <div className="text-center text-sm text-gray-500">
                {getTotalItems(order)} items
              </div>
            </div>
          </React.Fragment>
        ))}
      </div>

      <OrderSummaryModal
        isOpen={!!selectedOrder}
        onClose={() => setSelectedOrder(null)}
        order={selectedOrder}
      />
    </div>
  );
};

export default TableGrid;
