import { DishInterface } from "@/schemaValidations/interface/type_dish";
import {
  SetInterface,
  SetProtoDish
} from "@/schemaValidations/interface/types_set";
import { create } from "zustand";

interface OrderItemBase {
  id: number;
  quantity: number;
}

interface DishOrderItem extends OrderItemBase {
  type: "dish";
  dish: DishInterface;
}

interface SetOrderItem extends OrderItemBase {
  type: "set";
  set: SetInterface;
  modifiedDishes: SetProtoDish[];
}

type OrderItem = DishOrderItem | SetOrderItem;

interface OrderState {
  items: OrderItem[];
  isLoading: boolean;
  error: string | null;
  addItem: (
    item: DishInterface | SetInterface,
    quantity: number,
    modifiedDishes?: SetProtoDish[]
  ) => void;
  removeItem: (type: "dish" | "set", id: number) => void;
  updateQuantity: (type: "dish" | "set", id: number, quantity: number) => void;
  updateSetDishes: (setId: number, modifiedDishes: SetProtoDish[]) => void;
  clearOrder: () => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  getOrderSummary: () => {
    totalItems: number;
    totalPrice: number;
    dishes: { id: number; name: string; quantity: number; price: number }[];
    sets: {
      id: number;
      name: string;
      quantity: number;
      price: number;
      dishes: SetProtoDish[];
    }[];
  };
  findOrderItem: (type: "dish" | "set", id: number) => OrderItem | undefined;
}

const useOrderStore = create<OrderState>((set, get) => ({
  items: [],
  isLoading: false,
  error: null,

  addItem: (item, quantity, modifiedDishes) =>
    set((state) => {
      const type = "dishes" in item ? "dish" : "set";
      const existingItem = state.items.find(
        (i) => i.type === type && i.id === item.id
      );
      if (existingItem) {
        return {
          items: state.items.map((i) =>
            i.type === type && i.id === item.id
              ? { ...i, quantity: i.quantity + quantity }
              : i
          )
        };
      } else {
        const newItem: OrderItem =
          type === "dish"
            ? { type, id: item.id, quantity, dish: item as DishInterface }
            : {
                type,
                id: item.id,
                quantity,
                set: item as SetInterface,
                modifiedDishes: modifiedDishes || (item as SetInterface).dishes
              };
        return { items: [...state.items, newItem] };
      }
    }),

  removeItem: (type, id) =>
    set((state) => ({
      items: state.items.filter((i) => !(i.type === type && i.id === id))
    })),

  updateQuantity: (type, id, quantity) =>
    set((state) => ({
      items: state.items.map((i) =>
        i.type === type && i.id === id ? { ...i, quantity } : i
      )
    })),

  updateSetDishes: (setId, modifiedDishes) =>
    set((state) => ({
      items: state.items.map((i) =>
        i.type === "set" && i.id === setId ? { ...i, modifiedDishes } : i
      )
    })),

  clearOrder: () => set({ items: [] }),

  setLoading: (isLoading) => set({ isLoading }),

  setError: (error) => set({ error }),

  getOrderSummary: () => {
    const { items } = get();
    let totalItems = 0;
    let totalPrice = 0;
    const dishes: {
      id: number;
      name: string;
      quantity: number;
      price: number;
    }[] = [];
    const sets: {
      id: number;
      name: string;
      quantity: number;
      price: number;
      dishes: SetProtoDish[];
    }[] = [];

    items.forEach((item) => {
      totalItems += item.quantity;
      if (item.type === "dish") {
        totalPrice += item.dish.price * item.quantity;
        dishes.push({
          id: item.dish.id,
          name: item.dish.name,
          quantity: item.quantity,
          price: item.dish.price
        });
      } else {
        const setPrice = calculateSetPrice(item.modifiedDishes);
        totalPrice += setPrice * item.quantity;
        sets.push({
          id: item.set.id,
          name: item.set.name,
          quantity: item.quantity,
          price: setPrice,
          dishes: item.modifiedDishes
        });
      }
    });

    return { totalItems, totalPrice, dishes, sets };
  },

  findOrderItem: (type, id) => {
    return get().items.find((item) => item.type === type && item.id === id);
  }
}));

function calculateSetPrice(dishes: SetProtoDish[] | undefined): number {
  if (!dishes || dishes.length === 0) {
    return 0;
  }
  return dishes.reduce((acc, d) => {
    if (
      d &&
      d.dish &&
      typeof d.dish.price === "number" &&
      typeof d.quantity === "number"
    ) {
      return acc + d.dish.price * d.quantity;
    }
    return acc;
  }, 0);
}

export default useOrderStore;
