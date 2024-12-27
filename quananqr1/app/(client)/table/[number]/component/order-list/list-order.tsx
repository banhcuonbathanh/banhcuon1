import React from "react";

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
  Clock,
  UtensilsCrossed,
  UserCircle,
  CircleDollarSign,
  Flame
} from "lucide-react";
import useCartStore from "@/zusstand/new-order/new-order-zustand";
import { Order } from "@/schemaValidations/interface/type_order";

const OrderListPage = () => {
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

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case "pending":
        return "bg-yellow-500";
      case "completed":
        return "bg-green-500";
      case "cancelled":
        return "bg-red-500";
      default:
        return "bg-blue-500";
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD"
    }).format(price);
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

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {/* Current Order Section */}
        {current_order && (
          <Card className="w-full border-2 border-blue-500">
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="text-lg font-bold">
                  Current Order
                </CardTitle>
                <Badge className="bg-blue-500 text-white">Active</Badge>
              </div>
              <OrderContent order={current_order} />
            </CardHeader>
          </Card>
        )}

        {/* New Orders Section */}
        {new_order.map((order) => (
          <Card key={order.id} className="w-full">
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="text-lg font-bold">
                  Order #{order.id}
                </CardTitle>
                <Badge className={`${getStatusColor(order.status)} text-white`}>
                  {order.status}
                </Badge>
              </div>
              {/* <CardDescription className="space-y-2">
                <div className="flex items-center gap-2">
                  <UserCircle className="h-4 w-4" />
                  {order.is_guest
                    ? `Guest #${order.guest_id}`
                    : `User #${order.user_id}`}
                </div>
                <div className="flex items-center gap-2">
                  <Clock className="h-4 w-4" />
                  {formatDate(order.created_at)}
                </div>
                {order.chiliNumber > 0 && (
                  <div className="flex items-center gap-2">
                    <Flame className="h-4 w-4 text-red-500" />
                    Spice Level: {order.chiliNumber}
                  </div>
                )}
              </CardDescription> */}
            </CardHeader>
            <CardContent>
              <OrderContent order={order} />
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
};

// Separate component for order content to avoid repetition
interface OrderContentProps {
  order: Order;
}

const OrderContent = ({ order }: OrderContentProps) => {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD"
    }).format(price);
  };

  return (
    <div className="space-y-4">
      <div className="max-h-64 overflow-y-auto">
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
                <TableCell className="text-right">{item.quantity}</TableCell>
                <TableCell className="text-right">
                  {formatPrice(item.price * item.quantity)}
                </TableCell>
              </TableRow>
            ))}
            {order.set_items.map((set) => (
              <TableRow key={`set-${set.set_id}`} className="bg-slate-50">
                <TableCell>
                  <div>
                    <span className="font-medium">{set.name}</span>
                    <p className="text-xs text-gray-500">Set Menu</p>
                  </div>
                </TableCell>
                <TableCell className="text-right">{set.quantity}</TableCell>
                <TableCell className="text-right">
                  {formatPrice(set.price * set.quantity)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
      <div className="flex items-center justify-between border-t pt-4">
        <div className="flex items-center gap-2">
          <CircleDollarSign className="h-5 w-5" />
          <span className="font-semibold">Total:</span>
        </div>
        <span className="text-lg font-bold">
          {formatPrice(order.total_price)}
        </span>
      </div>
      {order.takeAway && (
        <Badge className="mt-2 bg-purple-500 text-white">Takeaway</Badge>
      )}
      {order.topping && (
        <div className="mt-2 text-sm text-gray-500">
          <span className="font-medium">Special Instructions:</span>{" "}
          {order.topping}
        </div>
      )}
    </div>
  );
};

export default OrderListPage;
