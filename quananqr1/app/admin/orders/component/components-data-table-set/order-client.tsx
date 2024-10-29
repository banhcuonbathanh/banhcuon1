"use client";

import { Heading } from "@/components/ui/heading";

import { Separator } from "@/components/ui/separator";

import {
  OrderDetailedResponse,
  PaginationInfo
} from "@/schemaValidations/interface/type_order";
import { columns } from "./order-columns";
import { OrderDataTable } from "@/components/ui/order-data-table";
import { useState } from "react";
import { get_Orders } from "@/zusstand/server/order-controller";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent
} from "@/components/ui/card";
//   const { data: sets, isLoading: setsLoading, error: setsError, refetch: refetchSets } = useSetListQuery();
interface OrderClientProps {
  initialData: OrderDetailedResponse[];
  initialPagination: PaginationInfo;
}

export const OrderClient: React.FC<OrderClientProps> = ({
  initialData,
  initialPagination
}) => {
  const [currentPage, setCurrentPage] = useState(
    initialPagination.current_page
  );
  const [data, setData] = useState(initialData);
  const [pagination, setPagination] = useState(initialPagination);
  const [isLoading, setIsLoading] = useState(false);

  const handleStatusChange = (orderId: number, newStatus: string) => {
    console.log(`Updating order ${orderId} status to ${newStatus}`);
  };

  const handlePaymentMethodChange = (orderId: number, newMethod: string) => {
    console.log(`Updating order ${orderId} payment method to ${newMethod}`);
  };

  const handlePageChange = async (newPage: number) => {
    setIsLoading(true);
    try {
      const orders = await get_Orders({
        page: newPage,
        page_size: pagination.page_size
      });

      setData(orders.data);
      setPagination(orders.pagination);
      setCurrentPage(newPage);
    } catch (error) {
      console.error("Error fetching orders:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
      <div className="space-y-2">
        <Card>
          <CardHeader>
            <CardTitle>Orders Management</CardTitle>
            <CardDescription>Manage your restaurant orders</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col">
              <div className="flex flex-row justify-between items-center">
                <Heading
                  title={`Orders (${pagination.total_items})`}
                  description={`Page ${currentPage} of ${pagination.total_pages}`}
                />
              </div>

              <Separator className="my-4" />

              <OrderDataTable
                columns={columns}
                data={data}
                searchKey="id"
                onStatusChange={handleStatusChange}
                onPaymentMethodChange={handlePaymentMethodChange}
              />

              <div className="flex items-center justify-between space-x-2 py-4">
                <div className="flex-1 text-sm text-muted-foreground">
                  Showing {data.length} of {pagination.total_items} orders
                </div>
                <div className="flex space-x-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePageChange(currentPage - 1)}
                    disabled={currentPage === 1 || isLoading}
                  >
                    Previous
                  </Button>
                  <div className="flex items-center space-x-2">
                    {[...Array(Math.min(5, pagination.total_pages))].map(
                      (_, idx) => {
                        const pageNum = idx + 1;
                        return (
                          <Button
                            key={pageNum}
                            variant={
                              pageNum === currentPage ? "default" : "outline"
                            }
                            size="sm"
                            onClick={() => handlePageChange(pageNum)}
                            disabled={isLoading}
                          >
                            {pageNum}
                          </Button>
                        );
                      }
                    )}
                  </div>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePageChange(currentPage + 1)}
                    disabled={
                      currentPage === pagination.total_pages || isLoading
                    }
                  >
                    Next
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
