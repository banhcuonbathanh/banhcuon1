import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent
} from "@/components/ui/card";
import React from "react";

import { OrderClient } from "./component/components-data-table-set/order-client";
import { get_Orders } from "@/zusstand/server/order-controller";

export default async function SetPage() {
  const order = await get_Orders({
    page: 1,
    page_size: 10
  });
  console.log("quananqr1/app/admin/orders/page.tsx order", order);
  // console.log("quananqr1/app/admin/set/page.tsx set", set[0].dishes);
  return (
    <main className="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
      <div className="space-y-2">
        <Card x-chunk="dashboard-06-chunk-0">
          <CardHeader>
            <CardTitle>Set Món ăn</CardTitle>
            <CardDescription>Quản lý set món ăn</CardDescription>
          </CardHeader>
          <CardContent>
            <OrderClient data={order} />
          </CardContent>
        </Card>
      </div>
    </main>
  );
}
