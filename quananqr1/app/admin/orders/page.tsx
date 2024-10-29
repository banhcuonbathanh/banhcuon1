import React from "react";

import { OrderClient } from "./component/components-data-table-set/order-client";
import { get_Orders } from "@/zusstand/server/order-controller";

export default async function OrdersPage() {
  const initialOrders = await get_Orders({
    page: 1,
    page_size: 10
  });

  return (
    <OrderClient
      initialData={initialOrders.data}
      initialPagination={initialOrders.pagination}
    />
  );
}
