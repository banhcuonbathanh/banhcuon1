import { DishSchema } from "@/app/employee/data-guest/public-dish/dish.schema";
import envConfig from "@/config";
import {
  DishListRes,
  DishListResType,
  SetListResType,
  SetSchema
} from "@/schemaValidations/dish.schema";
import { z } from "zod";

const get_Sets = async (): Promise<SetListResType> => {
  try {
    const baseUrl =
      envConfig.NEXT_PUBLIC_URL + envConfig.NEXT_PUBLIC_Get_set_intenal;
    // console.log("quananqr1/zusstand/server/set-controller.ts baseUrl", baseUrl);
    const response = await fetch(baseUrl, {
      method: "GET",
      cache: "no-store"
    });

    // }

    const data = await response.json();

    // console.log("quananqr1/zusstand/server/set-controller.ts data", data);

    // console.log("Received data:", JSON.stringify(data, null, 2));

    // Create a more lenient schema for parsing
    const LenientSetSchema = SetSchema.extend({
      createdAt: z.string().or(z.date()).optional(),
      updatedAt: z.string().or(z.date()).optional()
    }).transform((set) => ({
      ...set,
      createdAt: set.createdAt ? new Date(set.createdAt) : new Date(),
      updatedAt: set.updatedAt ? new Date(set.updatedAt) : new Date()
    }));

    // Validate the response data against the lenient schema
    const validatedData = z.array(LenientSetSchema).parse(data.data || []);

    return validatedData;
  } catch (error) {
    console.error("Error fetching or parsing sets:", error);
    if (error instanceof z.ZodError) {
      console.error(
        "Zod validation errors:",
        JSON.stringify(error.errors, null, 2)
      );
    }
    throw error;
  }
};

export { get_Sets };
