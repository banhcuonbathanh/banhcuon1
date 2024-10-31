import React, { useState } from "react";
import { ColumnDef } from "@tanstack/react-table";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select";
import {
  OrderDetailedDish,
  OrderDetailedResponse,
  OrderSetDetailed
} from "@/schemaValidations/interface/type_order";

const ORDER_STATUSES = ["ORDERING", "SERVING", "WAITING", "DONE"] as const;
type OrderStatus = (typeof ORDER_STATUSES)[number];

const PAYMENT_METHODS = ["CASH", "TRANSFER"] as const;
type PaymentMethod = (typeof PAYMENT_METHODS)[number];

interface TableMeta {
  onStatusChange?: (orderId: number, newStatus: string) => void;
  onPaymentMethodChange?: (orderId: number, newMethod: string) => void;
}

export const columns: ColumnDef<OrderDetailedResponse, any>[] = [
  {
    accessorKey: "order_name",
    header: "Name",
    cell: ({ row }) => (
      <div className="font-medium">#{row.getValue("order_name")}</div>
    )
  },
  {
    accessorKey: "table_number",
    header: "Table/away",
    cell: ({ row }) => (
      <div
        className={`text-center ${
          row.original.takeAway ? "bg-orange-600 rounded-md px-2 py-1" : ""
        }`}
      >
        {row.getValue("table_number")}
      </div>
    )
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row, table }) => {
      const [selectedStatus, setSelectedStatus] = useState<OrderStatus>(
        row.getValue("status") as OrderStatus
      );

      const statusStyles: Record<OrderStatus, string> = {
        ORDERING: "bg-blue-100 text-blue-800",
        SERVING: "bg-yellow-100 text-yellow-800",
        WAITING: "bg-orange-100 text-orange-800",
        DONE: "bg-green-100 text-green-800"
      };

      const meta = table.options.meta as TableMeta;

      return (
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
      );
    }
  },
  {
    accessorKey: "payment_method",
    header: "Payment",
    cell: ({ row, table }) => {
      const [selectedPayment, setSelectedPayment] = useState<PaymentMethod>(
        (row.getValue("payment_method") as PaymentMethod) || "CASH"
      );

      const paymentStyles: Record<PaymentMethod, string> = {
        CASH: "bg-emerald-50 text-emerald-700",
        TRANSFER: "bg-indigo-50 text-indigo-700"
      };

      const meta = table.options.meta as TableMeta;

      return (
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
      );
    }
  },
  {
    accessorKey: "data_set",
    header: "Sets",
    cell: ({ row }) => {
      const sets = row.getValue("data_set") as OrderSetDetailed[];
      return (
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
                    {/* Show total quantity for each dish in the set */}
                    {set.quantity * dish.quantity}x {dish.name}
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      );
    }
  },
  {
    accessorKey: "data_dish",
    header: "Individual Dishes",
    cell: ({ row }) => {
      const dishes = row.getValue("data_dish") as OrderDetailedDish[];
      return (
        <div className="space-y-1">
          {dishes.map((dish, index) => (
            <div key={`${dish.dish_id}-${index}`} className="text-sm">
              {dish.quantity}x {dish.name} (${dish.price})
            </div>
          ))}
        </div>
      );
    }
  },
  {
    accessorKey: "total_dishes",
    header: "Total Dishes Breakdown",
    cell: ({ row }) => {
      const sets = row.original.data_set;
      const individualDishes = row.original.data_dish;

      // Create a map to store total quantities for each dish type
      const dishTotals = new Map<string, number>();

      // Calculate totals from sets
      sets.forEach((set) => {
        set.dishes.forEach((dish) => {
          const totalQuantity = set.quantity * dish.quantity;
          const currentTotal = dishTotals.get(dish.name) || 0;
          dishTotals.set(dish.name, currentTotal + totalQuantity);
        });
      });

      // Add individual dishes to totals
      individualDishes.forEach((dish) => {
        const currentTotal = dishTotals.get(dish.name) || 0;
        dishTotals.set(dish.name, currentTotal + dish.quantity);
      });

      // Calculate grand total
      const grandTotal = Array.from(dishTotals.values()).reduce(
        (sum, qty) => sum + qty,
        0
      );

      return (
        <div className="space-y-2">
          {/* Display individual dish totals */}
          <div className="text-sm space-y-1">
            {Array.from(dishTotals.entries()).map(([dishName, quantity]) => (
              <div key={dishName} className="flex justify-between">
                <span>
                  {dishName}: {quantity}
                </span>
              </div>
            ))}
          </div>
          {/* Display grand total */}
          <div className="text-sm font-bold border-t pt-1">
            Total Items: {grandTotal}
          </div>
        </div>
      );
    }
  },
  {
    accessorKey: "bow_details",
    header: "Bowl Details",
    cell: ({ row }) => {
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
    }
  },
  {
    accessorKey: "total_price",
    header: "Total & Payment",
    cell: ({ row }) => {
      const totalPrice = row.getValue("total_price") as number;
      const [amountPaid, setAmountPaid] = useState<string>("");
      const [change, setChange] = useState<number | null>(null);

      const handlePaymentInput = (value: string) => {
        setAmountPaid(value);
        const numericValue = parseFloat(value) || 0;
        const changeAmount = numericValue - totalPrice;
        setChange(changeAmount >= 0 ? changeAmount : null);
      };

      return (
        <div className="space-y-2">
          <div className="font-medium text-right">Total: ${totalPrice}</div>
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
            <div
              className={`text-right text-sm ${
                change >= 0 ? "text-green-600" : "text-red-600"
              }`}
            >
              Change: ${change.toFixed(2)}
            </div>
          )}
        </div>
      );
    }
  }
];
