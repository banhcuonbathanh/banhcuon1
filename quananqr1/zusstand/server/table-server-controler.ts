import envConfig from "@/config";

import { z } from "zod";
import { TableListResType, TableSchema } from "../table/table.schema";
import { TableStatusValues } from "@/constants/type";

const get_tables = async (): Promise<TableListResType> => {
  try {
    const baseUrl =
      envConfig.NEXT_PUBLIC_URL + envConfig.NEXT_PUBLIC_intern_table_end_point;

    const response = await fetch(baseUrl, {
      method: "GET",
      cache: "no-store"
    });

    if (!response.ok) {
      const errorData = await response.json();
      console.log("Error response data:", errorData);
      throw new Error(
        `HTTP error! status: ${response.status}, message: ${errorData.message}`
      );
    }

    const data = await response.json();

   // Create a more lenient schema for parsing
   const LenientTableSchema = TableSchema.extend({
    createdAt: z.string().or(z.date()).optional(),
    updatedAt: z.string().or(z.date()).optional(),
    status: z.enum(TableStatusValues).optional()
  }).transform((table) => ({
    ...table,
    number: z.coerce.number().parse(table.number),
    capacity: z.coerce.number().parse(table.capacity),
    status: table.status || TableStatusValues[0], // Use first status as default if not provided
    token: table.token || "", // Provide a default empty string if token is missing
    createdAt: table.createdAt ? new Date(table.createdAt) : new Date(),
    updatedAt: table.updatedAt ? new Date(table.updatedAt) : new Date()
  }));

  // Validate the response data against the lenient schema
  const validatedData = z.array(LenientTableSchema).parse(data.data || data);

  return {
    data: validatedData,
    message: data.message || "Tables fetched successfully"
  };
  } catch (error) {
    console.error("Error fetching or parsing dishes:", error);
    if (error instanceof z.ZodError) {
      console.error(
        "Zod validation errors:",
        JSON.stringify(error.errors, null, 2)
      );
    }
    throw error;
  }
};

export { get_tables };
