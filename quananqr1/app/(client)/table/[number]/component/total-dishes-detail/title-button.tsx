import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Plus, Minus } from "lucide-react";

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
  const remaining = quantity - delivered;
  const pricePerUnit = details.totalPrice / details.quantity;

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

  const QuantityControl = ({ value }: { value: number }) => (
    <div className="flex justify-evenly">
      <Button
        variant="outline"
        size="sm"
        onClick={handleDecrement}
        disabled={quantity <= delivered}
        className="h-6 w-6 p-0"
      >
        <Minus className="h-3 w-3" />
      </Button>
      <span className="">{value}</span>
      <Button
        variant="outline"
        size="sm"
        onClick={handleIncrement}
        className="h-6 w-6 p-0"
      >
        <Plus className="h-3 w-3" />
      </Button>
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
        <QuantityControl value={quantity} />
      </div>
      <div className="p-3 text-right text-gray-300 border-t">{delivered}</div>
      <div className="p-3 text-right text-gray-300 border-t">{remaining}</div>
    </React.Fragment>
  );
};

export default TitleButton;
