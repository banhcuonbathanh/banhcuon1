import React, { useEffect, useState } from "react";
import useOrderStore, { OrderState } from "@/zusstand/order/order_zustand";
import { Order } from "@/schemaValidations/interface/type_order";

// Custom hook for orders with state management
const useOrders = () => {
  // Direct subscription to store using selector
  return useOrderStore((state) => state.listOfOrders);
};

const OrderItem = React.memo(({ order }: { order: Order }) => {
  const formattedDate = new Date(order.created_at).toLocaleDateString();
  const itemCount = order.dish_items.length + order.set_items.length;

  return (
    <li className="border rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex justify-between items-start">
        <div>
          <h3 className="font-medium">Order #{order.id}</h3>
          <p className="text-gray-600">Table: {order.table_number}</p>
          <p className="text-gray-600">Status: {order.status}</p>
          <p className="font-medium text-green-600">
            Total: ${order.total_price}
          </p>
        </div>
        <div className="text-right">
          <p className="text-sm text-gray-500">{formattedDate}</p>
          <p className="text-sm mt-2">
            {order.takeAway ? "Takeaway" : "Dine-in"}
          </p>
        </div>
      </div>
      <div className="mt-2 text-sm text-gray-600">Items: {itemCount}</div>
    </li>
  );
});

OrderItem.displayName = "OrderItem";

const OrdersList = () => {
  // Use the custom hook for orders
  const orders = useOrders();

  if (!orders?.length) {
    return (
      <div className="container mx-auto p-4">
        <h2 className="text-xl font-bold mb-4">Orders List</h2>
        <p className="text-gray-500">No orders available</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4">
      <h2 className="text-xl font-bold mb-4">Orders List</h2>
      <ul className="space-y-4">
        {orders.map((order) => (
          <OrderItem key={`order-${order.id}`} order={order} />
        ))}
      </ul>
    </div>
  );
};

OrdersList.displayName = "OrdersList";

export default React.memo(OrdersList);
