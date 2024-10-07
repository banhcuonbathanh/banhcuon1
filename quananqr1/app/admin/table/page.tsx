import { get_tables } from "@/zusstand/server/table-server-controler";
import React from "react";

const HomeTable = async () => {
  const tables = await get_tables();

  console.log("quananqr1/app/admin/table/page.tsx tables", tables);
  return <div>HomeTable</div>;
};

export default HomeTable;
