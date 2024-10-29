// stores/orderCreation.ts
import { create } from "zustand";
import { toast } from "@/components/ui/use-toast";
import envConfig from "@/config";
import { CreateOrderRequest } from "@/schemaValidations/interface/type_order";
import { useApiStore } from "@/zusstand/api/api-controller";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
import useOrderStore from "@/zusstand/order/order_zustand";

interface OrderCreationState {
  isLoading: boolean;
  createOrder: (
    bowlChili: number,
    bowlNoChili: number,
    Table_token: string
  ) => Promise<void>;
}

export const useOrderCreationStore = create<OrderCreationState>((set) => ({
  isLoading: false,

  createOrder: async (
    bowlChili: number,
    bowlNoChili: number,
    Table_token: string
  ) => {
    const { http } = useApiStore.getState();
    const { guest, user, isGuest, openLoginDialog } = useAuthStore.getState();
    const { tableNumber, getOrderSummary, clearOrder } =
      useOrderStore.getState();

    // Check authentication
    if (!user && !guest) {
      openLoginDialog();
      return;
    }

    const orderSummary = getOrderSummary();

    // Prepare order items
    const dish_items = orderSummary.dishes.map((dish) => ({
      dish_id: dish.id,
      quantity: dish.quantity
    }));

    const set_items = orderSummary.sets.map((set) => ({
      set_id: set.id,
      quantity: set.quantity
    }));

    // Determine IDs based on auth state
    const user_id = isGuest ? null : user?.id ?? null;
    const guest_id = isGuest ? guest?.id ?? null : null;

    const orderData: CreateOrderRequest = {
      guest_id,
      user_id,
      is_guest: isGuest,
      table_number: tableNumber,
      order_handler_id: 1,
      status: "pending",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      total_price: orderSummary.totalPrice,
      dish_items,
      set_items,
      bow_chili: bowlChili,
      bow_no_chili: bowlNoChili,
      takeAway: false,
      chiliNumber: 0,
      Table_token: Table_token
    };

    set({ isLoading: true });

    try {
      const link_order = `${envConfig.NEXT_PUBLIC_API_ENDPOINT}${envConfig.Order_External_End_Point}`;
      await http.post(link_order, orderData);

      toast({
        title: "Success",
        description: "Order has been created successfully"
      });

      clearOrder();
    } catch (error) {
      console.error("Order creation failed:", error);
      toast({
        variant: "destructive",
        title: "Error",
        description:
          error instanceof Error ? error.message : "Failed to create order"
      });
    } finally {
      set({ isLoading: false });
    }
  }
}));
