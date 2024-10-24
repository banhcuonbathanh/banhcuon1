"use client";





import { DataTable } from "@/components/ui/data-table";
import { Heading } from "@/components/ui/heading";

import { Separator } from "@/components/ui/separator";




import { Order } from "@/schemaValidations/interface/type_order";
import { columns } from "./order-columns";
//   const { data: sets, isLoading: setsLoading, error: setsError, refetch: refetchSets } = useSetListQuery();
interface OrderClientProps {
  data:  Order[];
}

export const OrderClient: React.FC<OrderClientProps> = ({ data }) => {


  

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
          <DataTable searchKey="name" columns={columns} data={data} />

      
        </div>
      )}
    </>
  );
};
