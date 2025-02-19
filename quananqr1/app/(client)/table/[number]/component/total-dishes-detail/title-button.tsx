import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Plus, Minus, Save } from "lucide-react";

interface GridItemProps {
  dishKey: string;
  details: {
    dishId: number;
    quantity: number;
    totalPrice: number;
  };
  deliveryData: {
    [key: string]: number;
  };
  remainingData: {
    [key: string]: number;
  };
  onUpdateQuantity?: (dishId: number, newQuantity: number) => void;
}

const TitleButton = ({
  dishKey,
  details,
  deliveryData,
  remainingData,
  onUpdateQuantity
}: GridItemProps) => {
  const title = dishKey.split("-")[0];
  const delivered = deliveryData[title] || 0;
  const [quantity, setQuantity] = useState(details.quantity);
  const originalRemaining = details.quantity - delivered;
  const newRemaining = quantity - delivered;
  const pricePerUnit = details.totalPrice / details.quantity;

  const originalQuantity = details.quantity;
  const difference = quantity - originalQuantity;

  const handleIncrement = () => {
    const newQuantity = quantity + 1;
    setQuantity(newQuantity);
    onUpdateQuantity?.(details.dishId, newQuantity);
  };

  const handleDecrement = () => {
    if (quantity > delivered) {
      const newQuantity = quantity - 1;
      setQuantity(newQuantity);
      onUpdateQuantity?.(details.dishId, newQuantity);
    }
  };

  const QuantityControl = () => (
    <div className="flex flex-col items-center space-y-2">
      <div className="flex flex-col items-center">
        <div className="text-gray-500">Original: {originalQuantity}</div>
        <div className="text-gray-600">
          (${(originalQuantity * pricePerUnit).toFixed(2)})
        </div>
      </div>

      <div className="flex items-center space-x-2">
        <Button
          variant="outline"
          size="sm"
          onClick={handleDecrement}
          disabled={quantity <= delivered}
          className="h-6 w-6 p-0"
        >
          <Minus className="h-3 w-3" />
        </Button>
        <span
          className={
            difference > 0
              ? "text-green-500"
              : difference < 0
              ? "text-red-500"
              : "text-gray-500"
          }
        >
          {difference > 0 ? `+${difference}` : difference}
          {difference !== 0 && (
            <span className="ml-1">
              ({difference > 0 ? "+" : ""}$
              {(difference * pricePerUnit).toFixed(2)})
            </span>
          )}
        </span>
        <Button
          variant="outline"
          size="sm"
          onClick={handleIncrement}
          className="h-6 w-6 p-0"
        >
          <Plus className="h-3 w-3" />
        </Button>
      </div>

      {difference !== 0 && (
        <>
          <div className="flex flex-col items-center">
            <div className="text-gray-500">New Quantity: {quantity}</div>
            <div className="text-gray-600">
              (${(quantity * pricePerUnit).toFixed(2)})
            </div>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={handleIncrement}
            className="flex items-center gap-2 px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors duration-200 shadow-sm hover:shadow-md"
          >
            <Save className="h-4 w-4" />
            <span>Update</span>
          </Button>
        </>
      )}
    </div>
  );

  const RemainingDisplay = () => (
    <div className="flex flex-col items-center space-y-2">
      <div className="text-gray-500">Remaining: {originalRemaining}</div>
      {difference !== 0 && (
        <>
          <div className="flex items-center space-x-2">
            <span
              className={
                difference > 0
                  ? "text-green-500"
                  : difference < 0
                  ? "text-red-500"
                  : "text-gray-500"
              }
            >
              {difference > 0 ? `+${difference}` : difference}
            </span>
          </div>
          <div className="text-gray-500">New Remaining: {newRemaining}</div>
        </>
      )}
    </div>
  );

  return (
    <React.Fragment key={`grid-${details.dishId}-${dishKey}`}>
      <div className="p-3 border-t">
        <div className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-gray-100 text-gray-900 hover:bg-gray-200 h-8 px-4 py-2">
          {title}
        </div>
      </div>
      <div className="p-3 border-t">
        <QuantityControl />
      </div>
      <div className="p-3 border-t">
        <div className="p-3 text-center text-gray-300">{delivered}</div>
      </div>
      <div className="p-3 border-t">
        <RemainingDisplay />
      </div>
    </React.Fragment>
  );
};

export default TitleButton;
