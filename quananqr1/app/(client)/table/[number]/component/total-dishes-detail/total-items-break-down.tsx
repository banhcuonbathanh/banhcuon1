import React from "react";
import TitleButton from "./title-button";

interface DishTotalDetails {
  quantity: number;
  totalPrice: number;
  dishId: number;
}

interface OrderGridProps {
  dishTotals: Map<string, DishTotalDetails>;
  deliveryData: {
    [key: string]: number;
  };
  remainingData: {
    [key: string]: number;
  };
}

const ItemsBreakdown = ({
  dishTotals,
  deliveryData,
  remainingData
}: OrderGridProps) => {
  // Calculate prices per unit for each dish
  const getPricePerUnit = (details: DishTotalDetails) =>
    details.totalPrice / details.quantity;

  return (
    <div className="w-full">
      <div className="grid grid-cols-4 gap-4">
        {/* Header */}
        <div className="p-3 text-left font-medium text-gray-600">Title</div>
        <div className="p-3 text-right font-medium text-gray-600">
          Total Qty
        </div>
        <div className="p-3 text-right font-medium text-gray-600">
          Delivered
        </div>
        <div className="p-3 text-right font-medium text-gray-600">
          Remaining
        </div>

        {/* Grid Items */}
        {Array.from(dishTotals.entries()).map(([dishKey, details]) => (
          <TitleButton
            key={`grid-${details.dishId}-${dishKey}`}
            dishKey={dishKey}
            details={details}
            deliveryData={deliveryData}
            remainingData={remainingData}
          />
        ))}

        {/* Footer Totals */}
        <div className="p-3 font-medium text-gray-500 border-t">Total</div>
        <div className="p-3 text-right font-medium text-gray-500 border-t">
          {Array.from(dishTotals.values()).reduce(
            (sum, item) => sum + item.quantity,
            0
          )}
          <br />
          <span className="text-primary">
            {Array.from(dishTotals.values()).reduce(
              (sum, item) => sum + item.totalPrice,
              0
            )}{" "}
          </span>
        </div>
        <div className="p-3 text-right font-medium text-gray-500 border-t">
          {Object.values(deliveryData).reduce((sum, val) => sum + val, 0)}
          <br />
          <span className="text-primary">
            {Array.from(dishTotals.entries()).reduce(
              (sum, [dishKey, details]) => {
                const title = dishKey.split("-")[0];
                const delivered = deliveryData[title] || 0;
                return (
                  sum + (delivered * details.totalPrice) / details.quantity
                );
              },
              0
            )}
          </span>
        </div>
        <div className="p-3 text-right font-medium text-gray-500 border-t">
          {Object.values(remainingData).reduce((sum, val) => sum + val, 0)}
          <br />
          <span className="text-gray-500">
            {Array.from(dishTotals.entries()).reduce(
              (sum, [dishKey, details]) => {
                const title = dishKey.split("-")[0];
                const remaining =
                  remainingData[title] ||
                  details.quantity - (deliveryData[title] || 0);
                return (
                  sum + (remaining * details.totalPrice) / details.quantity
                );
              },
              0
            )}
          </span>
        </div>
      </div>
    </div>
  );
};

export default ItemsBreakdown;
