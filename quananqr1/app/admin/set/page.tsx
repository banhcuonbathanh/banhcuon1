import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent
} from "@/components/ui/card";
import React from "react";
import { DishClient } from "../dish/components-data-table/dish-client";
import { SetClient } from "./component/components-data-table-set/set-client";
import { get_Sets } from "@/zusstand/server/set-controller";
import { get_dishes } from "@/zusstand/server/dish-controller";

export default async function SetPage() {
  const set = await get_Sets();

  //   console.log("quananqr1/app/admin/set/page.tsx set", set);
  return (
    <main className="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
      <div className="space-y-2">
        <Card x-chunk="dashboard-06-chunk-0">
          <CardHeader>
            <CardTitle>Set Món ăn</CardTitle>
            <CardDescription>Quản lý set món ăn</CardDescription>
          </CardHeader>
          <CardContent>
            <SetClient data={set} />
          </CardContent>
        </Card>
      </div>
    </main>
  );
}
