import { ColumnDef } from "@tanstack/react-table";

import { Order } from "@/schemaValidations/interface/type_order";

import React from "react";

export const columns: ColumnDef<Order>[] = [
  {
    accessorKey: "table_number",
    header: "Table"
  },
  {
    accessorKey: "dish_items",
    header: "Dishes",
    cell: ({ row }) => {
      return (
        <div className="space-y-1">
          {row.original.dish_items.map((item, index) => (
            <div
              key={`${item.dish_id}-${index}`}
              className="flex items-center gap-2"
            >
              <span className="font-medium">Dish #{item.dish_id}</span>
              <span className="text-sm text-gray-500">× {item.quantity}</span>
            </div>
          ))}
          {row.original.dish_items.length === 0 && (
            <span className="text-sm text-gray-500">No dishes ordered</span>
          )}
        </div>
      );
    }
  },
  {
    accessorKey: "set_items",
    header: "Sets",
    cell: ({ row }) => {
      return (
        <div className="space-y-1">
          {row.original.set_items.map((item, index) => (
            <div
              key={`${item.set_id}-${index}`}
              className="flex items-center gap-2"
            >
              <span className="font-medium">Set #{item.set_id}</span>
              <span className="text-sm text-gray-500">× {item.quantity}</span>
            </div>
          ))}
          {row.original.set_items.length === 0 && (
            <span className="text-sm text-gray-500">No sets ordered</span>
          )}
        </div>
      );
    }
  },
  {
    accessorKey: "total_price",
    header: "Total Price",
    cell: ({ row }) => {
      return (
        <div className="font-medium">
          ${row.original.total_price.toFixed(2)}
        </div>
      );
    }
  },
  {
    accessorKey: "bowls",
    header: "Bowls",
    cell: ({ row }) => {
      return (
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <span className="text-sm">With Chili:</span>
            <span className="font-medium">{row.original.bow_chili}</span>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-sm">No Chili:</span>
            <span className="font-medium">{row.original.bow_no_chili}</span>
          </div>
          <div className="text-sm text-gray-500">
            Total: {row.original.bow_chili + row.original.bow_no_chili}
          </div>
        </div>
      );
    }
  },
  {
    accessorKey: "takeAway",
    header: "Take Away",
    cell: ({ row }) => {
      return (
        <div
          className={`inline-flex px-2 py-1 rounded-full text-sm ${
            row.original.takeAway
              ? "bg-blue-100 text-blue-800"
              : "bg-gray-100 text-gray-800"
          }`}
        >
          {row.original.takeAway ? "Take Away" : "Dine In"}
        </div>
      );
    }
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const statusStyles = {
        pending: "bg-yellow-100 text-yellow-800",
        processing: "bg-blue-100 text-blue-800",
        completed: "bg-green-100 text-green-800",
        cancelled: "bg-red-100 text-red-800"
      };

      const status = row.original.status.toLowerCase();
      const styleClass =
        statusStyles[status as keyof typeof statusStyles] ||
        "bg-gray-100 text-gray-800";

      return (
        <div
          className={`inline-flex px-2 py-1 rounded-full text-sm ${styleClass}`}
        >
          {row.original.status}
        </div>
      );
    }
  }
];

export default columns;

// {
//   accessorKey: "id",
//   header: "Order ID"
// },
// {
//   accessorKey: "guest_id",
//   header: "Guest ID"
// },
// {
//   accessorKey: "user_id",
//   header: "User ID"
// },
// {
//   accessorKey: "table_number",
//   header: "Table"
// },
// {
//   accessorKey: "order_handler_id",
//   header: "Handler ID"
// },
// {
//   accessorKey: "status",
//   header: "Status"
// },
