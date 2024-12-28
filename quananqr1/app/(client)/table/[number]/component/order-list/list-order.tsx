import React, { useState } from "react";
import { Order } from "@/schemaValidations/interface/type_order";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table";
import {
  ChevronDown,
  ChevronUp,
  CircleDollarSign,
  Flame,
  List,
  Package
} from "lucide-react";
import { Button } from "@/components/ui/button";
import useCartStore from "@/zusstand/new-order/new-order-zustand";

interface OrderSummaryProps {
  order: Order;
  showDetails: boolean;
  onToggleDetails: () => void;
}

const OrderSummary: React.FC<OrderSummaryProps> = ({
  order,
  showDetails,
  onToggleDetails
}) => {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD"
    }).format(price);
  };

  const totalIndividualDishes = order.dish_items.reduce(
    (sum, item) => sum + item.quantity,
    0
  );
  const totalSetDishes = order.set_items.reduce(
    (sum, item) => sum + item.quantity,
    0
  );

  return (
    <Card className="w-full">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-xl">Order #{order.id}</CardTitle>
            <CardDescription>
              {totalIndividualDishes + totalSetDishes} items -{" "}
              {formatPrice(order.total_price)}
            </CardDescription>
          </div>
          {order.takeAway && (
            <Badge className="bg-purple-500 text-white">Takeaway</Badge>
          )}
        </div>
        {order.chiliNumber > 0 && (
          <div className="flex items-center gap-2 mt-2">
            <Flame className="h-4 w-4 text-red-500" />
            <span>Spice Level: {order.chiliNumber}</span>
          </div>
        )}
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Summary Section - Always Visible */}
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <span className="font-medium">Individual Dishes:</span>
            <span>{totalIndividualDishes} items</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="font-medium">Set Menus:</span>
            <span>{totalSetDishes} sets</span>
          </div>
        </div>

        <Button onClick={onToggleDetails} variant="outline" className="w-full">
          {showDetails ? (
            <div className="flex items-center">
              Hide Details <ChevronUp className="ml-2 h-4 w-4" />
            </div>
          ) : (
            <div className="flex items-center">
              Show Details <ChevronDown className="ml-2 h-4 w-4" />
            </div>
          )}
        </Button>

        {/* Detailed Section */}
        {showDetails && (
          <div className="space-y-4">
            {/* Individual Dishes */}
            {order.dish_items.length > 0 && (
              <div>
                <h4 className="font-medium mb-2">Individual Dishes</h4>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Item</TableHead>
                      <TableHead className="text-right">Qty</TableHead>
                      <TableHead className="text-right">Price</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {order.dish_items.map((item) => (
                      <TableRow key={`dish-${item.dish_id}`}>
                        <TableCell>{item.name}</TableCell>
                        <TableCell className="text-right">
                          {item.quantity}
                        </TableCell>
                        <TableCell className="text-right">
                          {formatPrice(item.price * item.quantity)}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}

            {/* Set Menus */}
            {order.set_items.length > 0 && (
              <div>
                <h4 className="font-medium mb-2">Set Menus</h4>
                <div className="space-y-4">
                  {order.set_items.map((set) => (
                    <div
                      key={`set-${set.set_id}`}
                      className="border rounded-lg p-4"
                    >
                      <div className="flex justify-between items-center mb-2">
                        <span className="font-medium">{set.name}</span>
                        <div className="text-right">
                          <div>Qty: {set.quantity}</div>
                          <div>{formatPrice(set.price * set.quantity)}</div>
                        </div>
                      </div>
                      <div className="text-sm text-gray-600">
                        <div className="mb-1">Included dishes:</div>
                        <ul className="list-disc pl-4">
                          {set.dishes.map((dish) => (
                            <li key={`set-${set.set_id}-dish-${dish.dish_id}`}>
                              {dish.name} x{dish.quantity}
                            </li>
                          ))}
                        </ul>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Special Instructions */}
            {order.topping && (
              <div className="mt-4 text-sm text-gray-600">
                <span className="font-medium">Special Instructions:</span>{" "}
                {order.topping}
              </div>
            )}
          </div>
        )}

        {/* Total Price - Always Visible */}
        <div className="flex items-center justify-between border-t pt-4">
          <div className="flex items-center gap-2">
            <CircleDollarSign className="h-5 w-5" />
            <span className="font-semibold">Total:</span>
          </div>
          <span className="text-lg font-bold">
            {formatPrice(order.total_price)}
          </span>
        </div>
      </CardContent>
    </Card>
  );
};

const OrderListPage = () => {
  const { new_order, isLoading, tableToken, tableNumber } = useCartStore();
  const [expandedOrders, setExpandedOrders] = useState<Set<number>>(new Set());

  // Toggle individual order
  const toggleOrderDetails = (orderId: number) => {
    const newExpanded = new Set(expandedOrders);
    if (newExpanded.has(orderId)) {
      newExpanded.delete(orderId);
    } else {
      newExpanded.add(orderId);
    }
    setExpandedOrders(newExpanded);
  };

  // Show all orders' details
  const showAllOrders = () => {
    const allOrderIds = new Set(new_order.map((order) => order.id));
    setExpandedOrders(allOrderIds);
  };

  // Hide all orders' details
  const hideAllOrders = () => {
    setExpandedOrders(new Set());
  };

  // Toggle show/hide all
  const toggleAllOrders = () => {
    if (expandedOrders.size === new_order.length) {
      hideAllOrders();
    } else {
      showAllOrders();
    }
  };

  if (isLoading) {
    return (
      <div className="p-4">
        <div className="text-center">Loading orders...</div>
      </div>
    );
  }

  return (
    <div className="p-4">
      {tableToken && tableNumber && (
        <div className="mb-4">
          <h2 className="text-2xl font-bold">Table #{tableNumber}</h2>
          <p className="text-sm text-gray-500">Token: {tableToken}</p>
        </div>
      )}

      {/* Global Controls */}
      <div className="flex gap-4 mb-6">
        <Button
          onClick={toggleAllOrders}
          className="flex items-center gap-2"
          variant="outline"
        >
          <Package className="h-4 w-4" />
          {expandedOrders.size === new_order.length
            ? "Hide All Details"
            : "Show All Details"}
        </Button>
      </div>

      {/* Orders Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {new_order.map((order) => (
          <OrderSummary
            key={order.id}
            order={order}
            showDetails={expandedOrders.has(order.id)}
            onToggleDetails={() => toggleOrderDetails(order.id)}
          />
        ))}
      </div>
    </div>
  );
};

export default OrderListPage;
