import React, { useState, useMemo } from "react";
import { ChevronDown, ChevronRight, Info } from "lucide-react";
import {
  OrderDetailedDish,
  OrderDetailedResponse
} from "../component/new-order-column";
import GroupToppings from "./toppping-display";
import DishSummary from "./dishes-summary";
import { logWithLevel } from "@/lib/log";

interface RestaurantSummaryProps {
  restaurantLayoutProps: OrderDetailedResponse[];
}

interface AggregatedDish extends OrderDetailedDish {}

interface GroupedOrder {
  orderName: string;
  characteristic?: string;
  tableNumber: number;
  orders: OrderDetailedResponse[];
  hasTakeAway: boolean;
}

// Helper Components
const CollapsibleSection: React.FC<{
  title: string;
  children: React.ReactNode;
}> = ({ title, children }) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="mt-2">
      <div
        className="flex items-center cursor-pointer select-none p-2 rounded"
        onClick={() => setIsOpen(!isOpen)}
      >
        <h3 className="text-md font-semibold">{title}</h3>
        {isOpen ? (
          <ChevronDown className="ml-2 h-4 w-4" />
        ) : (
          <ChevronRight className="ml-2 h-4 w-4" />
        )}
      </div>
      {isOpen && (
        <div className="p-2 border-l border-r border-b rounded-b">
          {children}
        </div>
      )}
    </div>
  );
};

const DishSummary1: React.FC<{ dish: AggregatedDish }> = ({ dish }) => {
  const [showDetails, setShowDetails] = useState(false);

  return (
    <div className="p-2 mb-2">
      <div className="flex items-center justify-between">
        <div
          className="flex-1 cursor-pointer"
          onClick={() => setShowDetails(!showDetails)}
        >
          <span className="font-bold">
            {dish.name} :{dish.quantity} - {"delivery"}
          </span>
          {/* <span className="ml-4">:{dish.quantity} 4</span> */}
        </div>
      </div>

      {showDetails && (
        <div className="mt-2 pl-4 text-gray-600">
          <div className="grid grid-cols-2 gap-1">
            <div className="font-medium">Price per Unit:</div>
            <div>${dish.price.toFixed(2)}</div>
            <div className="font-medium">Total Price:</div>
            <div>${(dish.price * dish.quantity).toFixed(2)}</div>
          </div>
        </div>
      )}
    </div>
  );
};

const OrderDetails: React.FC<{
  order: OrderDetailedResponse;
}> = ({ order }) => (
  <div className="border-b last:border-b-0 py-4">
    <div className="grid grid-cols-2 gap-2">
      <div className="font-semibold">Table Number:</div>
      <div>{order.table_number}</div>
      <div className="font-semibold">Status:</div>
      <div className={order.takeAway ? "text-red-600 font-bold" : ""}>
        {order.takeAway ? "Take Away" : order.status}
      </div>
      <div className="font-semibold">Total Price:</div>
      <div>${order.total_price.toFixed(2)}</div>
      <div className="font-semibold">Tracking Order:</div>
      <div>{order.tracking_order}</div>
      <div className="font-semibold">Chili Number:</div>
      <div>{order.chiliNumber}</div>
      {order.topping && (
        <>
          <div className="font-semibold">Toppings:</div>
          <div>{order.topping}</div>
        </>
      )}
    </div>

    <div className="mt-4">
      <h4 className="font-semibold mb-2">Individual Dishes:</h4>
      {order.data_dish.map((dish, index) => (
        <div key={`${dish.dish_id}-${index}`} className="ml-4 mb-2">
          <div>
            {dish.name} x{dish.quantity} (${dish.price.toFixed(2)} each)
          </div>
        </div>
      ))}
    </div>

    {order.data_set.length > 0 && (
      <div className="mt-4">
        <h4 className="font-semibold mb-2">Order Sets:</h4>
        {order.data_set.map((set, index) => (
          <div key={`${set.id}-${index}`} className="ml-4 mb-2">
            <div>
              {set.name} x{set.quantity} (${set.price.toFixed(2)} each)
            </div>
            <div className="ml-4 text-gray-600">
              Includes:
              {set.dishes.map((d, i) => (
                <React.Fragment key={d.dish_id}>
                  {i > 0 && ", "}
                  <span className="inline">
                    {d.name} (x{d.quantity})
                  </span>
                </React.Fragment>
              ))}
            </div>
          </div>
        ))}
      </div>
    )}
  </div>
);

const parseOrderName = (orderName: string): string => {
  const parts = orderName.split("-");
  return parts[0].trim();
};

const getOrdinalSuffix = (num: number): string => {
  const j = num % 10;
  const k = num % 100;
  if (j === 1 && k !== 11) return "st";
  if (j === 2 && k !== 12) return "nd";
  if (j === 3 && k !== 13) return "rd";
  return "th";
};

const aggregateDishes = (orders: OrderDetailedResponse[]): AggregatedDish[] => {
  const dishMap = new Map<number, AggregatedDish>();
  logWithLevel(
    {
      dishMap
    },
    "quananqr1/app/manage/admin/orders/restaurant-summary/restaurant-summary.tsx",
    "info",
    1 // You can use "debug", "info", "warn", or "error"
  );
  orders.forEach((order) => {
    // Add individual dishes
    order.data_dish.forEach((dish) => {
      const existingDish = dishMap.get(dish.dish_id);
      if (existingDish) {
        existingDish.quantity += dish.quantity;
      } else {
        dishMap.set(dish.dish_id, {
          ...dish,
          quantity: dish.quantity
        });
      }
    });

    // Add dishes from sets
    order.data_set.forEach((set) => {
      set.dishes.forEach((setDish) => {
        const existingDish = dishMap.get(setDish.dish_id);
        if (existingDish) {
          existingDish.quantity += setDish.quantity * set.quantity;
        } else {
          dishMap.set(setDish.dish_id, {
            ...setDish,
            quantity: setDish.quantity * set.quantity
          });
        }
      });
    });
  });

  return Array.from(dishMap.values());
};

export const RestaurantSummary: React.FC<RestaurantSummaryProps> = ({
  restaurantLayoutProps
}) => {
  const groupedOrders = useMemo(() => {
    const groups = new Map<string, GroupedOrder>();

    restaurantLayoutProps.forEach((order) => {
      const characteristic = parseOrderName(order.order_name);
      const groupKey = `${characteristic}-${order.table_number}`;

      if (!groups.has(groupKey)) {
        groups.set(groupKey, {
          orderName: characteristic,
          tableNumber: order.table_number,
          orders: [],
          hasTakeAway: false
        });
      }
      const group = groups.get(groupKey)!;
      group.orders.push(order);
      // Update hasTakeAway if any order in the group is takeaway
      if (order.takeAway) {
        group.hasTakeAway = true;
      }
    });

    return Array.from(groups.values());
  }, [restaurantLayoutProps]);

  return (
    <div className="p-4">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {groupedOrders.map((group) => {
          const aggregatedDishes = aggregateDishes(group.orders);
          logWithLevel(
            {
              aggregatedDishes
            },
            "quananqr1/app/manage/admin/orders/restaurant-summary/restaurant-summary.tsx",
            "info",
            2 // You can use "debug", "info", "warn", or "error"
          );
          return (
            <div
              key={`${group.orderName}-${group.tableNumber}`}
              className="shadow-md rounded-lg p-4 border"
            >
              <h3 className="text-xl font-semibold mb-4">
                {group.orderName} - Bàn {group.tableNumber}
                {group.hasTakeAway && (
                  <span className="ml-2 text-red-600">(Đem đi)</span>
                )}
              </h3>

              <div className="rounded-lg shadow-sm p-4">
                <CollapsibleSection title="Canh">
                  <GroupToppings orders={group.orders} />
                </CollapsibleSection>

                {/* <CollapsibleSection title="Món Ăn">
                  {aggregatedDishes.map((dish, index) => (
                    <DishSummary1
                      key={`${dish.dish_id}-${index}`}
                      dish={dish}
                    />
                  ))}
                </CollapsibleSection> */}

                <CollapsibleSection title="Món Ăn 123412341234">
                  {aggregatedDishes.map((dish, index) => (
                    <DishSummary
                      key={`${dish.dish_id}-${index}`}
                      dish={dish}
                      http={undefined}
                      orderStore={{
                        tableNumber: 0,
                        getOrderSummary: function () {
                          throw new Error("Function not implemented.");
                        },
                        clearOrder: function (): void {
                          throw new Error("Function not implemented.");
                        }
                      }}
                    />
                  ))}
                </CollapsibleSection>

                <CollapsibleSection title="Lần Gọi Đồ">
                  {group.orders.map((order, index) => (
                    <div key={order.id} className="mb-4 last:mb-0">
                      <div className="font-medium text-lg mb-2">
                        {`${index + 1}${getOrdinalSuffix(index + 1)} Order`}
                      </div>
                      <OrderDetails order={order} />
                    </div>
                  ))}
                </CollapsibleSection>

                <GroupSummary orders={group.orders} />
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default RestaurantSummary;

const GroupSummary: React.FC<{ orders: OrderDetailedResponse[] }> = ({
  orders
}) => {
  const [isDetailsVisible, setIsDetailsVisible] = useState(false);
  const totals = useMemo(() => {
    let dishTotal = 0;
    let setTotal = 0;

    orders.forEach((order) => {
      order.data_dish.forEach((dish) => {
        dishTotal += dish.price * dish.quantity;
      });

      order.data_set.forEach((set) => {
        setTotal += set.price * set.quantity;
      });
    });

    return {
      dishTotal,
      setTotal,
      grandTotal: dishTotal + setTotal
    };
  }, [orders]);

  return (
    <div className="mt-4 pt-4 border-t">
      <div
        className="cursor-pointer select-none"
        onClick={() => setIsDetailsVisible(!isDetailsVisible)}
      >
        <div className="grid grid-cols-2 gap-2">
          <div className="font-bold text-lg">Total:</div>
          <div className="text-right font-bold text-lg">
            ${totals.grandTotal.toFixed(2)}
            <ChevronDown
              className={`inline-block ml-2 h-4 w-4 transition-transform duration-200 ${
                isDetailsVisible ? "transform rotate-180" : ""
              }`}
            />
          </div>
        </div>

        {isDetailsVisible && (
          <div className="grid grid-cols-2 gap-2 mt-2 text-sm">
            <div className="font-medium">Individual Dishes:</div>
            <div className="text-right">${totals.dishTotal.toFixed(2)}</div>

            <div className="font-medium">Set Orders:</div>
            <div className="text-right">${totals.setTotal.toFixed(2)}</div>
          </div>
        )}
      </div>
    </div>
  );
};
