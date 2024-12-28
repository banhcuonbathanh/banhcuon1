// components/OrderDetailsList.tsx
import React from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table";
import { Order } from "@/schemaValidations/interface/type_order";

interface OrderDetailsListProps {
  order: Order;
  formatPrice: (price: number) => string;
}

export const OrderDetailsList: React.FC<OrderDetailsListProps> = ({
  order,
  formatPrice
}) => (
  <div className="space-y-4">
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
                <TableCell className="text-right">{item.quantity}</TableCell>
                <TableCell className="text-right">
                  {formatPrice(item.price * item.quantity)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    )}

    {order.set_items.length > 0 && (
      <div>
        <h4 className="font-medium mb-2">Set Menus</h4>
        <div className="space-y-4">
          {order.set_items.map((set) => (
            <div key={`set-${set.set_id}`} className="border rounded-lg p-4">
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
  </div>
);
