import React, { useMemo } from "react";
import { Card } from "@/components/ui/card";
import { OrderDetailedResponse } from "../component/new-order-column";

interface ToppingDisplayProps {
  orders: OrderDetailedResponse[];
}

const GroupToppings: React.FC<ToppingDisplayProps> = ({ orders }) => {
  const parseToppings = (toppingString: string): string[] => {
    // Split the string by commas and clean up whitespace
    return toppingString
      .split(",")
      .map((topping) => topping.trim())
      .filter((topping) => topping.length > 0);
  };

  const toppingAnalysis = useMemo(() => {
    const toppingCounts = new Map<string, number>();
    let totalOrders = 0;
    let ordersWithToppings = 0;

    orders.forEach((order) => {
      if (order.topping) {
        ordersWithToppings++;
        const toppings = parseToppings(order.topping);
        toppings.forEach((topping) => {
          toppingCounts.set(topping, (toppingCounts.get(topping) || 0) + 1);
        });
      }
      totalOrders++;
    });

    // Sort toppings by frequency
    const sortedToppings = Array.from(toppingCounts.entries())
      .sort((a, b) => b[1] - a[1])
      .map(([topping, count]) => ({
        name: topping,
        count,
        percentage: ((count / ordersWithToppings) * 100).toFixed(1)
      }));

    return {
      sortedToppings,
      totalOrders,
      ordersWithToppings
    };
  }, [orders]);

  return (
    <Card className="p-4">
      <div className="space-y-4">
        <div className="mb-4">
          <h3 className="text-lg font-semibold mb-2">Topping Summary</h3>
          <div className="text-sm text-gray-600">
            Orders with toppings: {toppingAnalysis.ordersWithToppings} of{" "}
            {toppingAnalysis.totalOrders}
          </div>
        </div>

        <div className="space-y-2">
          {toppingAnalysis.sortedToppings.map((topping, index) => (
            <div
              key={index}
              className="flex items-center justify-between p-2  "
            >
              <div className="flex-1">
                <span className="font-medium">{topping.name}</span>
              </div>
              <div className="flex items-center space-x-4">
                <span className="text-sm text-gray-600">
                  {topping.count} orders
                </span>
                <span className="text-sm text-gray-600">
                  ({topping.percentage}%)
                </span>
              </div>
            </div>
          ))}
        </div>

        {toppingAnalysis.sortedToppings.length === 0 && (
          <div className="text-gray-500 text-center py-4">
            No toppings found in these orders
          </div>
        )}
      </div>
    </Card>
  );
};

export default GroupToppings;
