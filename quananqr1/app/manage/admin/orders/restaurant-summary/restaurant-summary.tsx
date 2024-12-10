import React, { useState, useMemo } from "react";
import { ChevronDown, ChevronRight, Info } from "lucide-react";
import {
  OrderDetailedDish,
  OrderDetailedResponse
} from "../component/new-order-column";
import GroupToppings from "./toppping-display";

interface RestaurantSummaryProps {
  restaurantLayoutProps: OrderDetailedResponse[];
}

interface AggregatedDish extends OrderDetailedDish {}

interface GroupedOrder {
  orderName: string;
  characteristic?: string;
  tableNumber: number;
  orders: OrderDetailedResponse[];
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

const DishSummary: React.FC<{ dish: AggregatedDish }> = ({ dish }) => {
  const [showDetails, setShowDetails] = useState(false);

  return (
    <div className="p-2 rounded mb-2 border-b last:border-b-0">
      <div className="flex items-center justify-between">
        <div className="flex-1">
          <span className="font-medium">{dish.name}</span>
          <span className="ml-4">x{dish.quantity}</span>
        </div>
        <button
          onClick={() => setShowDetails(!showDetails)}
          className="flex items-center text-gray-600 hover:text-gray-800"
        >
          <Info className="h-4 w-4 mr-1" />
          {showDetails ? "Hide Details" : "Show Details"}
        </button>
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
          orders: []
        });
      }
      groups.get(groupKey)!.orders.push(order);
    });

    return Array.from(groups.values());
  }, [restaurantLayoutProps]);

  return (
    <div className="p-4">
      <h2 className="text-2xl font-bold mb-4">Restaurant Order Summary</h2>

      <CollapsibleSection title="All Orders">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {groupedOrders.map((group) => {
            const aggregatedDishes = aggregateDishes(group.orders);

            return (
              <div
                key={`${group.orderName}-${group.tableNumber}`}
                className="shadow-md rounded-lg p-4"
              >
                <h3 className="text-xl font-semibold mb-4">
                  {group.orderName} - Table {group.tableNumber}
                </h3>

                <div className="rounded-lg shadow-sm p-4">
                  <CollapsibleSection title="Toppings Summary">
                    <GroupToppings orders={group.orders} />
                  </CollapsibleSection>

                  <CollapsibleSection title="Aggregated Dishes">
                    {aggregatedDishes.map((dish, index) => (
                      <DishSummary
                        key={`${dish.dish_id}-${index}`}
                        dish={dish}
                      />
                    ))}
                  </CollapsibleSection>

                  <CollapsibleSection title="Individual Orders">
                    {group.orders.map((order, index) => (
                      <div key={order.id} className="mb-4 last:mb-0">
                        <div className="font-medium text-lg mb-2">
                          {`${index + 1}${getOrdinalSuffix(index + 1)} Order`}
                        </div>
                        <OrderDetails order={order} />
                      </div>
                    ))}
                  </CollapsibleSection>
                </div>
              </div>
            );
          })}
        </div>
      </CollapsibleSection>
    </div>
  );
};

export default RestaurantSummary;
