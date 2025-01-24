"use client";

import React, { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table";
import useCartStore from "@/zusstand/new-order/new-order-zustand";
import { Package, UtensilsCrossed } from "lucide-react";

const OrderSummaryPage = () => {
  const [isMounted, setIsMounted] = useState(false);
  const { getOrderSummary } = useCartStore();
  const summary = getOrderSummary();

  useEffect(() => {
    setIsMounted(true);
  }, []);

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD"
    }).format(price);
  };

  if (!isMounted) return null;

  // check hidration start

  // check hidrationi end

  return (
    <div className="p-4 max-w-4xl mx-auto">
      {summary.totalItems === 0 ? (
        <div className="p-8 text-center">
          <UtensilsCrossed className="mx-auto h-12 w-12 text-gray-400 mb-4" />
          <h3 className="text-lg font-medium">No items in cart</h3>
          <p className="text-gray-500">
            Add some dishes or set menus to get started
          </p>
        </div>
      ) : (
        <>
          <Card className="mb-6">
            <CardHeader>
              <CardTitle>Order Summary</CardTitle>
              <CardDescription>
                Total Items: {summary.totalItems} | Total Price:{" "}
                {formatPrice(summary.totalPrice)}
              </CardDescription>
            </CardHeader>
          </Card>

          {summary.dishes.length > 0 && (
            <Card className="mb-6">
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Package className="h-5 w-5" />
                  <CardTitle className="text-lg">Individual Dishes</CardTitle>
                </div>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Item</TableHead>
                      <TableHead>Description</TableHead>
                      <TableHead className="text-right">Quantity</TableHead>
                      <TableHead className="text-right">Unit Price</TableHead>
                      <TableHead className="text-right">Total</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {summary.dishes.map((dish) => (
                      <TableRow key={dish.dish_id}>
                        <TableCell className="font-medium">
                          {dish.name}
                        </TableCell>
                        <TableCell className="max-w-xs truncate">
                          {dish.description}
                        </TableCell>
                        <TableCell className="text-right">
                          {dish.quantity}
                        </TableCell>
                        <TableCell className="text-right">
                          {formatPrice(dish.price)}
                        </TableCell>
                        <TableCell className="text-right">
                          {formatPrice(dish.price * dish.quantity)}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          )}

          {summary.sets.length > 0 && (
            <Card>
              <CardHeader>
                <div className="flex items-center gap-2">
                  <UtensilsCrossed className="h-5 w-5" />
                  <CardTitle className="text-lg">Set Menus</CardTitle>
                </div>
              </CardHeader>
              <CardContent>
                {summary.sets.map((set) => (
                  <div key={set.set_id} className="mb-6 last:mb-0">
                    <div className="flex justify-between items-center mb-2">
                      <h3 className="text-lg font-semibold">{set.name}</h3>
                      <div className="text-right">
                        <div>Quantity: {set.quantity}</div>
                        <div className="font-medium">
                          {formatPrice(set.price * set.quantity)}
                        </div>
                      </div>
                    </div>

                    <div className="bg-gray-50 rounded-lg p-4">
                      <h4 className="font-medium mb-2">Included Dishes:</h4>
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableHead>Dish</TableHead>
                            <TableHead className="text-right">
                              Quantity
                            </TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {set.dishes.map((dish) => (
                            <TableRow key={`${set.set_id}-${dish.dish_id}`}>
                              <TableCell>{dish.name}</TableCell>
                              <TableCell className="text-right">
                                {dish.quantity}
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>
          )}

          <Card className="mt-6">
            <CardContent className="pt-6">
              <div className="flex justify-between items-center text-lg">
                <div className="font-semibold">Total Items:</div>
                <div>{summary.totalItems}</div>
              </div>
              <div className="flex justify-between items-center text-xl font-bold mt-2">
                <div>Total Price:</div>
                <div>{formatPrice(summary.totalPrice)}</div>
              </div>
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
};

export default OrderSummaryPage;
