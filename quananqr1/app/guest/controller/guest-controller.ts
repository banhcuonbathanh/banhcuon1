import envConfig from "@/config";
import {
  DishListRes,
  DishListResType
} from "@/zusstand/dished/domain/dish.schema";

const get_dishes = async (): Promise<DishListResType> => {
  // Change return type
  console.log("Fetching dishes from controller...");

  try {
    const baseUrl =
      envConfig.NEXT_PUBLIC_URL + envConfig.NEXT_PUBLIC_Get_Dished_intenal;
    console.log("Fetching dishes from controller... 222 baseUrl", baseUrl);

    console.log("Fetching from URL:", baseUrl);

    const response = await fetch(baseUrl, {
      method: "GET",
      cache: "no-store"
    });

    // Check for response.ok before parsing JSON
    if (!response.ok) {
      const errorData = await response.json(); // Log the error response body
      console.log("Error response data:", errorData);
      throw new Error(
        `HTTP error! status: ${response.status}, message: ${errorData.message}`
      );
    }

    const data = await response.json();
    console.log("Fetched dishes data:", data);

    // Validate with your schema

    // Return the whole response object
    return {
      data: data, // Array of dishes
      message: "data.message" // Message from the API
    };
  } catch (error) {
    console.error("Error fetching dishes:", error);
    throw error; // This will propagate the error to the calling function
  }
};

export { get_dishes };
