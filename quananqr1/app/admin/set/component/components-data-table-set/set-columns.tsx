"use client";
import { Button } from "@/components/ui/button";
import { SetType } from "@/schemaValidations/dish.schema";
import { ColumnDef } from "@tanstack/react-table";
import { MoreHorizontal, ChevronDown, ChevronUp } from "lucide-react";
import React from "react";

export const columns: ColumnDef<SetType>[] = [
  {
    accessorKey: "id",
    header: "ID",
    size: 60
  },
  {
    accessorKey: "name",
    header: "Set Name",
    size: 200
  },
  {
    accessorKey: "description",
    header: "Description",
    size: 300,
    cell: ({ row }) => {
      const description = row.original.description;
      return description ? description : "N/A";
    }
  },
  {
    accessorKey: "dishes",
    header: "Dishes",
    size: 300,
    cell: ({ row }) => {
      const [isExpanded, setIsExpanded] = React.useState(false);
      const dishes = row.original.dishes;

      return (
        <div>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setIsExpanded(!isExpanded)}
          >
            {dishes.length} Dishes{" "}
            {isExpanded ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
          </Button>
          {isExpanded && (
            <ul className="mt-2 space-y-1">
              {dishes.map((dish) => (
                <li key={dish.id} className="text-sm">
                  {dish.name} - ${dish.price.toFixed(2)}
                </li>
              ))}
            </ul>
          )}
        </div>
      );
    }
  },
  {
    accessorKey: "createdAt",
    header: "Created At",
    size: 150,
    cell: ({ row }) => new Date(row.original.createdAt).toLocaleDateString()
  },
  {
    accessorKey: "updatedAt",
    header: "Updated At",
    size: 150,
    cell: ({ row }) => new Date(row.original.updatedAt).toLocaleDateString()
  },
  {
    id: "actions",
    size: 60,
    cell: ({ row }) => (
      <Button variant="ghost" size="icon">
        <MoreHorizontal className="h-4 w-4" />
      </Button>
    )
  }
];
