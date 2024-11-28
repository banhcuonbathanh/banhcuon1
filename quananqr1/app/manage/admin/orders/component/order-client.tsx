"use client";

import { Heading } from "@/components/ui/heading";

import { Separator } from "@/components/ui/separator";

import {
  OrderDetailedResponse,
  PaginationInfo
} from "@/schemaValidations/interface/type_order";

import { useCallback, useEffect, useState } from "react";
import { get_Orders } from "@/zusstand/server/order-controller";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent
} from "@/components/ui/card";

import { DeliveryInterface } from "@/schemaValidations/interface/type_delivery";
import { YourComponent1 } from "./admin-table";
import { useApiStore } from "@/zusstand/api/api-controller";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
import { useWebSocketStore } from "@/zusstand/web-socket/websocketStore";
import { WebSocketMessage } from "@/zusstand/web-socket/websoket-service";

//   const { data: sets, isLoading: setsLoading, error: setsError, refetch: refetchSets } = useSetListQuery();
interface OrderClientProps {
  initialData: OrderDetailedResponse[];
  initialPagination: PaginationInfo;
  deliveryData: DeliveryInterface[];
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
  // ---------------------
  const { http } = useApiStore();
  const { guest, user, isGuest, openLoginDialog } = useAuthStore();

  const { connect, disconnect, addMessageHandler } = useWebSocketStore();
  // connect(isGuest ? guest : user, isGuest);

  useEffect(() => {
    connect(isGuest ? guest : user, isGuest);
    console.log(
      "quananqr1/app/admin/orders/component/order-client.tsx connect"
    );
    return () => {
      disconnect();
    };
  }, [connect, disconnect, guest, user, isGuest]);

  const handleWebSocketMessage = useCallback(
    (message: WebSocketMessage) => {
      console.log(
        "quananqr1/app/admin/orders/component/order-client.tsx message: 111111",
        message
      );
      if (message.type === "NEW_ORDER") {
        console.log(
          "quananqr1/app/admin/orders/component/order-client.tsx message: 111111",
          message
        );

        setData((prevData) => {
          if (currentPage === 1) {
            const newOrder: OrderDetailedResponse = {
              id: 12,
              table_number: 1,
              status: "message.content.status",
              created_at: " message.content.timestamp",
              data_set: [],
              data_dish: [],
              guest_id: 0,
              user_id: 0,
              is_guest: false,
              order_handler_id: 0,
              updated_at: "",
              total_price: 0,
              bow_chili: 0,
              bow_no_chili: 0,
              takeAway: false,
              chiliNumber: 0,
              table_token: "",
              order_name: "",
              deliveryData: undefined
            };
            return [newOrder, ...prevData.slice(0, -1)];
          }
          return prevData;
        });

        setPagination((prev) => ({
          ...prev,
          total_items: prev.total_items + 1
        }));
      }
    },
    [currentPage]
  );
  useEffect(() => {
    // Register the handler on mount
    const unsubscribe = addMessageHandler(handleWebSocketMessage);

    // Cleanup the handler on unmount
    return unsubscribe;
  }, [addMessageHandler, handleWebSocketMessage]);
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

              <YourComponent1 initialData={initialData} />

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

                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => {
                      connect(isGuest ? guest : user, isGuest);
                      console.log(
                        "quananqr1/app/admin/orders/component/order-client.tsx Button "
                      );
                    }}
                    disabled={
                      currentPage === pagination.total_pages || isLoading
                    }
                  >
                    connect
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
