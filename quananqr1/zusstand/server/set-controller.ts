import envConfig from "@/config";
import { DishInterface } from "@/schemaValidations/interface/type_dish";
import {
  SetInterface,
  SetListResType,
  SetProtoDish
} from "@/schemaValidations/interface/types_set";

const get_Sets = async (): Promise<SetInterface[]> => {
  try {
    const baseUrl =
      envConfig.NEXT_PUBLIC_URL + envConfig.NEXT_PUBLIC_Get_set_intenal;
    // console.log("quananqr1/zusstand/server/set-controller.ts baseUrl", baseUrl);
    const response = await fetch(baseUrl, {
      method: "GET",
      cache: "no-store"
    });

    const data = await response.json();

    // console.log(
    //   "quananqr1/zusstand/server/set-controller.ts data 2323232323232323232",
    //   data
    // );

    // Validate and transform the data
    const validatedData: SetInterface[] = data.data.map(
      (set: SetInterface) => ({
        id: set.id,
        name: set.name,
        description: set.description,
        dishes: set.dishes.map((setProtoDish: SetProtoDish) => ({
          dishId: setProtoDish.dishId,
          quantity: setProtoDish.quantity,
          dish: {
            ...setProtoDish.dish,
            price: Number(setProtoDish.dish.price) // Ensure price is a number
          }
        })),
        userId: set.userId,
        created_at: set.created_at,
        updated_at: set.updated_at,
        is_favourite: Boolean(set.is_favourite), // Ensure is_favourite is a boolean
        like_by: set.like_by || [], // Ensure like_by is an array, defaulting to empty if null
        is_public: Boolean(set.is_public) // Ensure is_public is a boolean
      })
    );

    return validatedData;
  } catch (error) {
    console.error("Error fetching or parsing sets:", error);
    throw error;
  }
};

export { get_Sets };
