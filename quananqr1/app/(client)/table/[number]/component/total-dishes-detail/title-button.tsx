import React from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from "@/components/ui/dialog";

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
}

const TitleButton = ({
  dishKey,
  details,
  deliveryData,
  remainingData
}: GridItemProps) => {
  const title = dishKey.split("-")[0];
  const delivered = deliveryData[title] || 0;
  const remaining = remainingData[title] || details.quantity - delivered;

  return (
    <React.Fragment key={`grid-${details.dishId}-${dishKey}`}>
      <div className="p-3 border-t">
        <Dialog>
          <DialogTrigger className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-gray-100 text-gray-900 hover:bg-gray-200 h-8 px-4 py-2">
            {title}
          </DialogTrigger>
          <DialogContent className="bg-white shadow-lg border rounded-lg">
            <DialogHeader className="bg-white">
              <DialogTitle className="text-lg font-semibold text-gray-900">
                {title} Details
              </DialogTitle>
            </DialogHeader>
            <div className="space-y-4 bg-white p-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="text-sm font-medium text-gray-500">
                  Order ID:
                </div>
                <div className="text-sm text-gray-900">{details.dishId}</div>
                <div className="text-sm font-medium text-gray-500">
                  Quantity:
                </div>
                <div className="text-sm text-gray-900">{details.quantity}</div>
                <div className="text-sm font-medium text-gray-500">
                  Delivered:
                </div>
                <div className="text-sm text-gray-900">{delivered}</div>
                <div className="text-sm font-medium text-gray-500">
                  Remaining:
                </div>
                <div className="text-sm text-gray-900">{remaining}</div>
                <div className="text-sm font-medium text-gray-500">
                  Total Price:
                </div>
                <div className="text-sm text-gray-900">
                  ${details.totalPrice.toFixed(2)}
                </div>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      </div>
      <div className="p-3 text-right text-gray-300 border-t">
        {details.quantity}
      </div>
      <div className="p-3 text-right text-gray-300 border-t">{delivered}</div>
      <div className="p-3 text-right text-gray-300 border-t">{remaining}</div>
    </React.Fragment>
  );
};

export default TitleButton;
