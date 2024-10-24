import envConfig from "@/config";

import {
  GetOrdersRequest,
  Order
} from "@/schemaValidations/interface/type_order";

export const get_Orders = async (
  params: GetOrdersRequest
): Promise<Order[]> => {
  try {
    const baseUrl = `${envConfig.NEXT_PUBLIC_URL}${envConfig.Order_Internal_End_Point}`;
    const queryParams = new URLSearchParams({
      page: params.page.toString(),
      page_size: params.page_size.toString()
    });

    console.log(
      "quananqr1/zusstand/server/order-controller.ts baseUrl",
      `${baseUrl}?${queryParams}`
    );

    const response = await fetch(`${baseUrl}?${queryParams}`, {
      method: "GET",
      cache: "no-store"
    });

    const data = await response.json();

    console.log(
      "quananqr1/zusstand/server/order-controller.ts data",
      data.data
    );

    const validatedData: Order[] = data.data.map((order: any) => ({
      id: order.id,
      guest_id: order.guest_id,
      user_id: order.user_id,
      is_guest: order.is_guest,
      table_number: order.table_number,
      order_handler_id: order.order_handler_id,
      status: order.status,
      created_at: order.created_at,
      updated_at: order.updated_at,
      total_price: order.total_price,
      dish_items: order.dish_items || [],
      set_items: order.set_items || [],
      bow_chili: order.bow_chili,
      bow_no_chili: order.bow_no_chili,
      takeAway: order.take_away,
      chiliNumber: order.chili_number
    }));

    return validatedData;
  } catch (error) {
    console.error("Error fetching or parsing orders:", error);
    throw error;
  }
};
