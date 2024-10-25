"use client";

import { DataTable } from "@/components/ui/data-table";
import { Heading } from "@/components/ui/heading";

import { Separator } from "@/components/ui/separator";

import { OrderDetailedResponse } from "@/schemaValidations/interface/type_order";
import { columns } from "./order-columns";
import { OrderDataTable } from "@/components/ui/order-data-table";
//   const { data: sets, isLoading: setsLoading, error: setsError, refetch: refetchSets } = useSetListQuery();
interface OrderClientProps {
  data: OrderDetailedResponse[];
}

export const OrderClient: React.FC<OrderClientProps> = ({ data }) => {
  const handleStatusChange = (orderId: number, newStatus: string) => {
    // Implement your status update logic here
    console.log(`Updating order ${orderId} status to ${newStatus}`);
  };

  const handlePaymentMethodChange = (orderId: number, newMethod: string) => {
    // Implement your payment method update logic here
    console.log(`Updating order ${orderId} payment method to ${newMethod}`);
  };

  return (
    <>
      {data && (
        <div className="flex flex-col">
          <div className="flex flex-rol">
            <Heading
              title={`set (${data.length})`}
              description="Manage set for your store"
            />
            {/* <AddSet /> */}
          </div>

          <Separator />
          <OrderDataTable
            columns={columns}
            data={data}
            searchKey="id"
            onStatusChange={handleStatusChange}
            onPaymentMethodChange={handlePaymentMethodChange}
          />
        </div>
      )}
    </>
  );
};
