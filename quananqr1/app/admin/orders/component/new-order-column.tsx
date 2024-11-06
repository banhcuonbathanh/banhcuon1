"use client";

import React, { useState, useEffect } from "react";
import { ColumnDef } from "@tanstack/react-table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";

const ORDER_STATUSES = ["ORDERING", "SERVING", "WAITING", "DONE"] as const;
type OrderStatus = (typeof ORDER_STATUSES)[number];

const PAYMENT_METHODS = ["CASH", "TRANSFER"] as const;
type PaymentMethod = (typeof PAYMENT_METHODS)[number];

interface TableMeta {
  onStatusChange?: (orderId: number, newStatus: string) => void;
  onPaymentMethodChange?: (orderId: number, newMethod: string) => void;
  onDeliveryUpdate?: (
    orderId: number,
    dishName: string,
    deliveredQuantity: number
  ) => void;
}
// -----------------------

const CombinedOrderDetails = ({
  sets,
  dishes,
  row,
  meta
}: {
  sets: OrderSetDetailed[];
  dishes: OrderDetailedDish[];
  row: any;
  meta: TableMeta;
}) => {
  const [deliveryState, setDeliveryState] = useState<Map<string, number>>(
    new Map()
  );
  const [amountPaid, setAmountPaid] = useState<string>("");
  const [change, setChange] = useState<number | null>(null);
  const totalPrice = row.original.total_price as number;

  useEffect(() => {
    if (row.original.deliveryData) {
      setDeliveryState(new Map(Object.entries(row.original.deliveryData)));
    }
  }, [row.original.deliveryData]);

  const calculateDishTotals = () => {
    const dishTotals = new Map<string, { quantity: number; price: number }>();

    // Calculate totals from sets
    sets.forEach((set) => {
      set.dishes.forEach((dish) => {
        const totalQuantity = set.quantity * dish.quantity;
        const dishPrice = (set.price / set.dishes.length) * dish.quantity;
        const current = dishTotals.get(dish.name) || { quantity: 0, price: 0 };
        dishTotals.set(dish.name, {
          quantity: current.quantity + totalQuantity,
          price: current.price + dishPrice * set.quantity
        });
      });
    });

    // Calculate totals from individual dishes
    dishes.forEach((dish) => {
      const current = dishTotals.get(dish.name) || { quantity: 0, price: 0 };
      dishTotals.set(dish.name, {
        quantity: current.quantity + dish.quantity,
        price: current.price + dish.price * dish.quantity
      });
    });

    return dishTotals;
  };

  const handleDeliveryUpdate =
    (dishName: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
      const dishTotals = calculateDishTotals();
      const totalQuantity = dishTotals.get(dishName)?.quantity || 0;
      const newDelivered = Math.min(
        parseInt(e.target.value) || 0,
        totalQuantity
      );

      const newState = new Map(deliveryState);
      newState.set(dishName, newDelivered);
      setDeliveryState(newState);

      meta?.onDeliveryUpdate?.(row.original.id, dishName, newDelivered);
    };

  const handlePaymentInput = (value: string) => {
    setAmountPaid(value);
    const numericValue = parseFloat(value) || 0;
    const changeAmount = numericValue - totalPrice;
    setChange(changeAmount >= 0 ? changeAmount : null);
  };

  const dishTotals = calculateDishTotals();

  return (
    <div className="space-y-4">
      {/* Sets Section */}
      {sets && sets.length > 0 && (
        <div className="border-b border-gray-100 pb-4">
          <div className="grid grid-cols-5 gap-2 mb-2 text-sm font-semibold text-gray-700">
            <div className="col-span-2">Sets</div>
            <div className="text-center">Quantity</div>
            <div className="text-center text-green-600">Delivered</div>
            <div className="text-center text-orange-600">Price</div>
          </div>
          {sets.map((set) => (
            <div key={set.id} className="space-y-2">
              <div className="grid grid-cols-5 gap-2 items-center">
                <div className="col-span-2 text-sm font-medium">{set.name}</div>
                <div className="text-center">{set.quantity}</div>
                <div className="text-center">-</div>
                <div className="text-center">${set.price}</div>
              </div>
              <div className="pl-4">
                {set.dishes.map((dish, index) => {
                  const totalQty = set.quantity * dish.quantity;
                  const delivered = deliveryState.get(dish.name) || 0;
                  return (
                    <div
                      key={`${dish.dish_id}-${index}`}
                      className="grid grid-cols-5 gap-2 items-center text-sm text-gray-600"
                    >
                      <div className="col-span-2">{dish.name}</div>
                      <div className="text-center">{totalQty}</div>
                      <div className="flex justify-center">
                        <Input
                          type="number"
                          value={delivered}
                          onChange={handleDeliveryUpdate(dish.name)}
                          className="w-16 h-7 text-center text-green-600"
                          min="0"
                          max={totalQty}
                        />
                      </div>
                      <div className="text-center">-</div>
                    </div>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Individual Dishes Section */}
      {dishes && dishes.length > 0 && (
        <div className="space-y-2">
          <div className="grid grid-cols-5 gap-2 text-sm font-semibold text-gray-700">
            <div className="col-span-2">Individual Dishes</div>
            <div className="text-center">Quantity</div>
            <div className="text-center text-green-600">Delivered</div>
            <div className="text-center text-orange-600">Price</div>
          </div>
          {dishes.map((dish, index) => {
            const delivered = deliveryState.get(dish.name) || 0;
            return (
              <div
                key={`${dish.dish_id}-${index}`}
                className="grid grid-cols-5 gap-2 items-center text-sm"
              >
                <div className="col-span-2">{dish.name}</div>
                <div className="text-center">{dish.quantity}</div>
                <div className="flex justify-center">
                  <Input
                    type="number"
                    value={delivered}
                    onChange={handleDeliveryUpdate(dish.name)}
                    className="w-16 h-7 text-center text-green-600"
                    min="0"
                    max={dish.quantity}
                  />
                </div>
                <div className="text-center">${dish.price}</div>
              </div>
            );
          })}
        </div>
      )}

      {/* Payment Section */}
      <div className="border-t pt-4 mt-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <span className="text-sm font-medium text-gray-600">
              Amount Paid
            </span>
            <div className="flex items-center gap-2">
              <Input
                type="number"
                placeholder="Amount"
                value={amountPaid}
                onChange={(e) => handlePaymentInput(e.target.value)}
                className="w-24 h-8 text-right"
              />
              <span className="text-sm">$</span>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <div className="text-sm font-medium text-gray-600">
              Total: ${totalPrice.toFixed(2)}
            </div>
            {change !== null && (
              <div
                className={`text-sm font-medium ${
                  change >= 0 ? "text-green-600" : "text-red-600"
                }`}
              >
                Change: ${change.toFixed(2)}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default CombinedOrderDetails;
//----------------------------
// Header Component

const OrderHeader = ({ row, meta }: { row: any; meta: TableMeta }) => {
  const orderName = row.original.order_name.split("_")[0];
  const tableNumber = row.original.table_number;
  const [selectedStatus, setSelectedStatus] = useState<OrderStatus>(
    row.original.status as OrderStatus
  );
  const [selectedPayment, setSelectedPayment] = useState<PaymentMethod>("CASH");
  const isTakeAway = row.original.takeAway;
  const [isCalculationsVisible, setIsCalculationsVisible] = useState(false);

  // Bowl details
  const withChili = row.original.bow_chili || 0;
  const noChili = row.original.bow_no_chili || 0;
  const totalBowls = withChili + noChili;
  const chiliNumber = row.original.chiliNumber || 0;

  const getDishNameColor = (isCalculationsVisible: boolean, name: string) => {
    const commonDishes = ["banh", "trung", "gio"];
    if (commonDishes.includes(name.toLowerCase())) {
      return isCalculationsVisible
        ? "bg-gray-600 text-white px-2 py-1 rounded text-xs"
        : "hidden";
    }
    return "text-gray-500";
  };

  const statusStyles: Record<OrderStatus, string> = {
    ORDERING: "bg-blue-100 text-blue-800",
    SERVING: "bg-yellow-100 text-yellow-800",
    WAITING: "bg-orange-100 text-orange-800",
    DONE: "bg-green-100 text-green-800"
  };

  const paymentStyles: Record<PaymentMethod, string> = {
    CASH: "bg-emerald-50 text-emerald-700",
    TRANSFER: "bg-indigo-50 text-indigo-700"
  };

  return (
    <div className="flex flex-col rounded-t-lg border-b">
      {/* First Row - Main Order Details */}
      <div className="flex items-center gap-6 p-4 pb-2">
        {/* Name Section */}
        <div className="flex flex-row min-w-[100px]">
          <span className="text-sm font-medium text-gray-600">Name</span>
          <span className="ml-2 text-sm font-medium text-gray-600">
            {orderName}
          </span>
        </div>

        {/* Table Section */}
        <div className="flex flex-row min-w-[100px]">
          <span className="text-sm font-medium text-gray-600">Table</span>
          <div
            className={`ml-2 text-sm font-medium ${
              isTakeAway
                ? "bg-orange-600 text-white rounded-md px-2 py-1"
                : "text-gray-600"
            }`}
          >
            {tableNumber}
          </div>
        </div>

        {/* Status Section */}
        <div className="flex flex-row min-w-[150px]">
          <span className="text-sm font-medium text-gray-600">Status</span>
          <div className="ml-2">
            <Select
              value={selectedStatus}
              onValueChange={(newStatus: OrderStatus) => {
                setSelectedStatus(newStatus);
                meta?.onStatusChange?.(row.original.id, newStatus);
              }}
            >
              <SelectTrigger
                className={`w-[120px] h-8 ${statusStyles[selectedStatus]}`}
              >
                <SelectValue>{selectedStatus}</SelectValue>
              </SelectTrigger>
              <SelectContent>
                {ORDER_STATUSES.map((orderStatus) => (
                  <SelectItem
                    key={orderStatus}
                    value={orderStatus}
                    className={statusStyles[orderStatus]}
                  >
                    {orderStatus}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Payment Method Section */}
        <div className="flex flex-row min-w-[150px]">
          <span className="text-sm font-medium text-gray-600">Payment</span>
          <div className="ml-2">
            <Select
              value={selectedPayment}
              onValueChange={(newMethod: PaymentMethod) => {
                setSelectedPayment(newMethod);
                meta?.onPaymentMethodChange?.(row.original.id, newMethod);
              }}
            >
              <SelectTrigger
                className={`w-[120px] h-8 ${paymentStyles[selectedPayment]}`}
              >
                <SelectValue>
                  {selectedPayment === "CASH" ? "Cash" : "Transfer"}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                {PAYMENT_METHODS.map((method) => (
                  <SelectItem
                    key={method}
                    value={method}
                    className={paymentStyles[method]}
                  >
                    {method === "CASH" ? "Cash" : "Transfer"}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Total Amount Section */}
        <div className="flex flex-row min-w-[150px]">
          <span className="text-sm font-medium text-gray-600">Total</span>
          <span className="ml-2 text-sm font-medium text-gray-800">
            ${row.original.total_price}
          </span>
        </div>
      </div>

      {/* Second Row - Bowl Details */}
      {(totalBowls > 0 || (isTakeAway && chiliNumber > 0)) && (
        <div className="px-4 py-2 flex items-center">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1">
              <span className="text-sm font-medium text-gray-600">Bowls:</span>
            </div>
            <div className="flex gap-3">
              {withChili > 0 && (
                <span className="text-sm px-3 py-1 bg-red-50 text-red-600 rounded-md flex items-center gap-1">
                  <span className="text-lg">🥬</span>
                  <span className="font-medium">{withChili}</span>
                </span>
              )}
              {noChili > 0 && (
                <span className="text-sm px-3 py-1 bg-green-50 text-green-600 rounded-md flex items-center gap-1">
                  <span className="text-lg">⛔</span>
                  <span className="font-medium">{noChili}</span>
                </span>
              )}
              {isTakeAway && chiliNumber > 0 && (
                <span className="text-sm px-3 py-1 bg-orange-50 text-orange-600 rounded-md flex items-center gap-1">
                  <span className="text-lg">🏃</span>
                  <span className="font-medium">{chiliNumber}</span>
                </span>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Third Row - Order Details */}
      <div className="px-4 py-2 rounded-b-lg">
        <div className="flex flex-wrap gap-x-8 gap-y-2">
          {/* Sets Section */}
          {row.original.data_set && row.original.data_set.length > 0 && (
            <div className="flex-1 min-w-[300px]">
              <div className="flex items-center gap-2">
                <span className="text-sm font-semibold text-gray-700">
                  Sets:{" "}
                </span>
                <button
                  onClick={() =>
                    setIsCalculationsVisible(!isCalculationsVisible)
                  }
                  className="text-gray-500 hover:text-gray-700"
                >
                  {isCalculationsVisible ? (
                    <ChevronUp className="w-4 h-4" />
                  ) : (
                    <ChevronDown className="w-4 h-4" />
                  )}
                </button>
              </div>
              <div className="mt-1 space-y-2">
                {row.original.data_set.map((set: OrderSetDetailed) => (
                  <div key={set.id}>
                    {/* Always visible set information */}
                    <div className="text-sm text-gray-700">
                      {set.quantity}x {set.name} (${set.price})
                    </div>

                    {/* Togglable detailed calculations */}
                    {isCalculationsVisible && (
                      <div className="pl-4 text-sm space-y-1">
                        {set.dishes.map((dish, index) => (
                          <div
                            key={`${dish.dish_id}-${index}`}
                            className="flex items-center gap-2"
                          >
                            <span className="text-gray-500">{dish.name}</span>
                            <span className="bg-gray-600 text-white px-2 py-1 rounded text-xs">
                              {set.quantity} x {dish.quantity} ={" "}
                              {set.quantity * dish.quantity}
                            </span>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Individual Dishes Section */}
          {row.original.data_dish && row.original.data_dish.length > 0 && (
            <div className="flex-1 min-w-[300px]">
              <span className="text-sm font-semibold text-gray-700">
                Individual Dishes:{" "}
              </span>
              <div className="mt-1 space-y-1">
                {row.original.data_dish.map(
                  (dish: OrderDetailedDish, index: number) => (
                    <div
                      key={`${dish.dish_id}-${index}`}
                      className="text-sm text-gray-700"
                    >
                      {dish.quantity}x {dish.name} (${dish.price})
                    </div>
                  )
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// Order Details Component
const OrderDetails = ({
  sets,
  dishes
}: {
  sets: OrderSetDetailed[];
  dishes: OrderDetailedDish[];
}) => {
  return (
    <div className="space-y-4">
      {/* Sets Section */}
      {sets && sets.length > 0 && (
        <div>
          <div className="text-sm font-semibold text-gray-700 mb-2">Sets</div>
          <div className="space-y-2">
            {sets.map((set) => (
              <div
                key={set.id}
                className="border-b border-gray-100 pb-2 last:border-0"
              >
                <div className="text-sm font-medium">
                  {set.quantity}x {set.name} (${set.price})
                </div>
                <div className="pl-4 text-sm text-gray-600">
                  {set.dishes.map((dish, index) => (
                    <div key={`${dish.dish_id}-${index}`}>
                      {set.quantity} x {dish.quantity} ={" "}
                      {set.quantity * dish.quantity} {dish.name}
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Individual Dishes Section */}
      {dishes && dishes.length > 0 && (
        <div>
          <div className="text-sm font-semibold text-gray-700 mb-2">
            Individual Dishes
          </div>
          <div className="space-y-1">
            {dishes.map((dish, index) => (
              <div key={`${dish.dish_id}-${index}`} className="text-sm">
                {dish.quantity}x {dish.name} (${dish.price})
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

// Order Tracking Component
// ... (previous code remains the same until OrderTracking component)

const OrderTracking = ({ row, meta }: { row: any; meta: TableMeta }) => {
  const [deliveryState, setDeliveryState] = useState<Map<string, number>>(
    new Map()
  );
  const [amountPaid, setAmountPaid] = useState<string>("");
  const [change, setChange] = useState<number | null>(null);
  const totalPrice = row.original.total_price as number;

  useEffect(() => {
    if (row.original.deliveryData) {
      setDeliveryState(new Map(Object.entries(row.original.deliveryData)));
    }
  }, [row.original.deliveryData]);

  // ... [previous helper functions remain the same]
  const calculateDishTotals = () => {
    const dishTotals = new Map<string, { quantity: number; price: number }>();

    row.original.data_set?.forEach((set: OrderSetDetailed) => {
      set.dishes.forEach((dish) => {
        const totalQuantity = set.quantity * dish.quantity;
        const dishPrice = dish.price * dish.quantity;
        const current = dishTotals.get(dish.name) || { quantity: 0, price: 0 };
        dishTotals.set(dish.name, {
          quantity: current.quantity + totalQuantity,
          price: current.price + dishPrice * set.quantity
        });
      });
    });

    row.original.data_dish?.forEach((dish: OrderDetailedDish) => {
      const current = dishTotals.get(dish.name) || { quantity: 0, price: 0 };
      dishTotals.set(dish.name, {
        quantity: current.quantity + dish.quantity,
        price: current.price + dish.price * dish.quantity
      });
    });

    return dishTotals;
  };

  const calculateTotals = () => {
    const dishTotals = calculateDishTotals();
    const totalOrderValue = Array.from(dishTotals.values()).reduce(
      (sum, total) => sum + total.price,
      0
    );
    const deliveredValue = Array.from(dishTotals.entries()).reduce(
      (sum, [dishName, totals]) => {
        const delivered = deliveryState.get(dishName) || 0;
        const pricePerUnit = totals.price / totals.quantity;
        return sum + delivered * pricePerUnit;
      },
      0
    );
    const remainingValue = totalOrderValue - deliveredValue;

    return {
      totalOrderValue,
      deliveredValue,
      remainingValue
    };
  };

  const handlePaymentInput = (value: string) => {
    setAmountPaid(value);
    const numericValue = parseFloat(value) || 0;
    const { remainingValue } = calculateTotals();
    const changeAmount = numericValue - remainingValue;
    setChange(changeAmount >= 0 ? changeAmount : null);
  };

  const handleDeliveryUpdate =
    (dishName: string) => async (newValue: number) => {
      const dishTotals = calculateDishTotals();
      const totalQuantity = dishTotals.get(dishName)?.quantity || 0;
      const newDelivered = Math.min(newValue, totalQuantity);

      try {
        const response = await meta?.onDeliveryUpdate?.(
          row.original.id,
          dishName,
          newDelivered
        );

        const newState = new Map(deliveryState);
        newState.set(dishName, newDelivered);
        setDeliveryState(newState);
      } catch (error) {
        console.error("Failed to update delivery:", error);
      }
    };

  const dishTotals = calculateDishTotals();
  const { totalOrderValue, deliveredValue, remainingValue } = calculateTotals();

  return (
    <div className="space-y-2 w-full overflow-x-auto">
      <div className="min-w-[280px] sm:min-w-[320px]">
        <div className="grid grid-cols-4 gap-1 sm:gap-2 px-1 sm:px-2 py-1 rounded-t text-xs sm:text-sm font-medium">
          <div className="col-span-1">Dish</div>
          <div className="text-center">Total($)</div>
          <div className="text-center text-green-600">Del($)</div>
          <div className="text-center text-orange-600">Rem($)</div>
        </div>

        {Array.from(dishTotals.entries()).map(([dishName, totals]) => {
          const delivered = deliveryState.get(dishName) || 0;
          const pricePerUnit = totals.price / totals.quantity;
          const deliveredValue = delivered * pricePerUnit;
          const remainingValue = totals.price - deliveredValue;
          const isComplete = delivered === totals.quantity;

          return (
            <div
              key={dishName}
              className={`grid grid-cols-4 gap-1 sm:gap-2 items-center py-1 border-b border-gray-100 last:border-0 ${
                isComplete ? "bg-green-50" : ""
              }`}
            >
              <div className="col-span-1 text-xs sm:text-sm font-medium truncate">
                {dishName}
              </div>
              <div className="text-center text-xs sm:text-sm">
                {totals.quantity}x${pricePerUnit.toFixed(2)}
              </div>
              <div className="flex items-center gap-0.5 sm:gap-1 justify-center">
                <NumericKeypadInput
                  value={delivered}
                  onChange={() => {}}
                  onSubmit={handleDeliveryUpdate(dishName)}
                  max={totals.quantity}
                  className="w-8 sm:w-12 h-6 sm:h-7 text-center text-green-600 text-xs sm:text-sm"
                />
                <span className="text-[10px] sm:text-xs text-green-600">
                  ${deliveredValue.toFixed(2)}
                </span>
              </div>
              <div className="text-center text-xs sm:text-sm">
                <span
                  className={isComplete ? "text-green-600" : "text-orange-600"}
                >
                  ${remainingValue.toFixed(2)}
                </span>
              </div>
            </div>
          );
        })}

        <div className="grid grid-cols-4 gap-1 sm:gap-2 items-center py-2 font-medium">
          <div className="col-span-1 text-xs sm:text-sm">Total</div>
          <div className="text-center text-xs sm:text-sm">
            ${totalOrderValue.toFixed(2)}
          </div>
          <div className="text-center text-green-600 text-xs sm:text-sm">
            ${deliveredValue.toFixed(2)}
          </div>
          <div className="text-center text-orange-600 text-xs sm:text-sm">
            ${remainingValue.toFixed(2)}
          </div>
        </div>

        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 sm:gap-4 pt-4 border-t">
          <div className="flex items-center gap-2 sm:gap-4">
            <div className="text-xs sm:text-sm font-medium">
              Amount Due:{" "}
              <span className="text-orange-600">
                ${remainingValue.toFixed(2)}
              </span>
            </div>
            <div className="flex items-center gap-1 sm:gap-2">
              <Input
                type="number"
                placeholder="0.00"
                value={amountPaid}
                onChange={(e) => handlePaymentInput(e.target.value)}
                className="w-20 sm:w-24 h-7 sm:h-8 text-right text-xs sm:text-sm"
              />
              <span className="text-xs sm:text-sm">$</span>
            </div>
          </div>
          {change !== null && (
            <div
              className={`text-xs sm:text-sm font-medium ${
                change >= 0 ? "text-green-600" : "text-red-600"
              }`}
            >
              Change: ${change.toFixed(2)}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// ... (rest of the code remains the same)

// Payment Information Component
const PaymentInformation = ({ row, meta }: { row: any; meta: TableMeta }) => {
  const [selectedPayment, setSelectedPayment] = useState<PaymentMethod>("CASH");
  const totalPrice = row.original.total_price as number;
  const [amountPaid, setAmountPaid] = useState<string>("");
  const [change, setChange] = useState<number | null>(null);

  const paymentStyles: Record<PaymentMethod, string> = {
    CASH: "bg-emerald-50 text-emerald-700",
    TRANSFER: "bg-indigo-50 text-indigo-700"
  };

  const handlePaymentInput = (value: string) => {
    setAmountPaid(value);
    const numericValue = parseFloat(value) || 0;
    const changeAmount = numericValue - totalPrice;
    setChange(changeAmount >= 0 ? changeAmount : null);
  };

  return (
    <div className="space-y-3">
      <div className="flex flex-col">
        <span className="text-sm font-medium text-gray-600">
          Payment Method
        </span>
        <div className="mt-1">
          <Select
            value={selectedPayment}
            onValueChange={(newMethod: PaymentMethod) => {
              setSelectedPayment(newMethod);
              meta?.onPaymentMethodChange?.(row.original.id, newMethod);
            }}
          >
            <SelectTrigger
              className={`w-[120px] h-8 ${paymentStyles[selectedPayment]}`}
            >
              <SelectValue>
                {selectedPayment === "CASH" ? "Cash" : "Transfer"}
              </SelectValue>
            </SelectTrigger>
            <SelectContent>
              {PAYMENT_METHODS.map((method) => (
                <SelectItem
                  key={method}
                  value={method}
                  className={paymentStyles[method]}
                >
                  {method === "CASH" ? "Cash" : "Transfer"}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="flex flex-col">
        <span className="text-sm font-medium text-gray-600">Total Amount</span>
        <span className="font-medium mt-1">${totalPrice}</span>
      </div>

      <div className="flex flex-col">
        <span className="text-sm font-medium text-gray-600">Amount Paid</span>
        <div className="mt-1 space-y-2">
          <div className="flex items-center gap-2">
            <Input
              type="number"
              placeholder="Amount paid"
              value={amountPaid}
              onChange={(e) => handlePaymentInput(e.target.value)}
              className="w-24 h-8 text-right"
            />
            <span className="text-sm">$</span>
          </div>
          {change !== null && (
            <div className="flex flex-col">
              <div
                className={`text-sm ${
                  change >= 0 ? "text-green-600" : "text-red-600"
                }`}
              >
                Change: ${change.toFixed(2)}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// Bowl Details Component
const BowlDetails = ({ row }: { row: any }) => {
  const withChili = row.original.bow_chili;
  const noChili = row.original.bow_no_chili;
  const total = withChili + noChili;
  const isTakeAway = row.original.takeAway;
  const chiliNumber = row.original.chiliNumber;

  return total > 0 || (isTakeAway && chiliNumber > 0) ? (
    <div className="space-y-1 text-sm">
      {withChili > 0 && <div>With Chili: {withChili}</div>}
      {noChili > 0 && <div>No Chili: {noChili}</div>}
      {isTakeAway && chiliNumber > 0 && (
        <div className="font-medium">Takeaway Chili: {chiliNumber}</div>
      )}
    </div>
  ) : null;
};

// Main Column Definition
export const columns: ColumnDef<OrderDetailedResponse, any>[] = [
  {
    id: "order_details",
    cell: ({ row, table }) => {
      const meta = table.options.meta as TableMeta;

      return (
        <div className="flex flex-col w-full">
          <OrderHeader row={row} meta={meta} />

          {/* <div className="space-y-4">
              <h3 className="font-semibold">Order Details</h3>
              <OrderDetails
                sets={row.original.data_set}
                dishes={row.original.data_dish}
              />
            </div> */}

          {/* <div className="space-y-4">
            <h3 className="font-semibold">Order Tracking</h3>
            <OrderTracking row={row} meta={meta} />
          </div> */}

          <div className="space-y-4">
            <h3 className="font-semibold">Order Tracking</h3>
            <OrderTracking12 row={row} meta={meta} />
          </div>
        </div>
      );
    },
    meta: {
      skipHeaderRender: true
    }
  }
];

// Example usage of the table (you can include this if needed)
const OrderTable = ({ data }: { data: OrderDetailedResponse[] }) => {
  const [rowSelection, setRowSelection] = useState({});

  const handleStatusChange = (orderId: number, newStatus: string) => {
    console.log(`Order ${orderId} status changed to ${newStatus}`);
    // Implement your status change logic here
  };

  const handlePaymentMethodChange = (orderId: number, newMethod: string) => {
    console.log(`Order ${orderId} payment method changed to ${newMethod}`);
    // Implement your payment method change logic here
  };

  const handleDeliveryUpdate = (
    orderId: number,
    dishName: string,
    deliveredQuantity: number
  ) => {
    console.log(
      `Order ${orderId} dish ${dishName} delivery updated to ${deliveredQuantity}`
    );
    // Implement your delivery update logic here
  };

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    onRowSelectionChange: setRowSelection,
    state: {
      rowSelection
    },
    meta: {
      onStatusChange: handleStatusChange,
      onPaymentMethodChange: handlePaymentMethodChange,
      onDeliveryUpdate: handleDeliveryUpdate
    }
  });

  return (
    <div className="rounded-md border">
      <Table>
        <TableBody>
          {table.getRowModel().rows.map((row) => (
            <TableRow key={row.id}>
              {row.getVisibleCells().map((cell) => (
                <TableCell key={cell.id}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
};

// Required imports for the table component
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";
import {
  flexRender,
  getCoreRowModel,
  useReactTable
} from "@tanstack/react-table";
import {
  OrderDetailedDish,
  OrderDetailedResponse,
  OrderSetDetailed
} from "@/schemaValidations/interface/type_order";
import { ChevronUp, ChevronDown } from "lucide-react";
import NumericKeypadInput from "./numberpad-dialog";
import OrderTracking12 from "./order-tracking";

// Additional types that might be needed
interface OrderTableProps {
  data: OrderDetailedResponse[];
}

// Export the components
export { OrderTable };
export type { OrderTableProps, OrderDetailedResponse, TableMeta };

// You might also want to include these utility types/interfaces if needed elsewhere
export type { OrderStatus, PaymentMethod, OrderSetDetailed, OrderDetailedDish };
