import { DishInterface } from "@/schemaValidations/interface/type_dish";
import {
  DishOrderItem,
  SetOrderItem
} from "@/schemaValidations/interface/type_order";
import {
  SetInterface,
  SetProtoDish
} from "@/schemaValidations/interface/types_set";
import { create } from "zustand";

interface OrderState {
  dishItems: DishOrderItem[];
  setItems: SetOrderItem[];
  isLoading: boolean;
  error: string | null;

  // Functions remain the same, but with updated types
  addDishItem: (dish: DishInterface, quantity: number) => void;
  removeDishItem: (id: number) => void;
  updateDishQuantity: (id: number, quantity: number) => void;
  addSetItem: (
    set: SetInterface,
    quantity: number,
    modifiedDishes?: SetProtoDish[]
  ) => void;
  removeSetItem: (id: number) => void;
  updateSetQuantity: (id: number, quantity: number) => void;
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
  findDishOrderItem: (id: number) => DishOrderItem | undefined;
  findSetOrderItem: (id: number) => SetOrderItem | undefined;
}

const useOrderStore = create<OrderState>((set, get) => ({
  dishItems: [],
  setItems: [],
  isLoading: false,
  error: null,

  addDishItem: (dish, quantity) =>
    set((state) => {
      const existingItem = state.dishItems.find((i) => i.id === dish.id);
      if (existingItem) {
        return {
          dishItems: state.dishItems.map((i) =>
            i.id === dish.id ? { ...i, quantity: i.quantity + quantity } : i
          )
        };
      } else {
        const newItem: DishOrderItem = {
          id: dish.id,
          quantity,
          dish: dish
        };
        return { dishItems: [...state.dishItems, newItem] };
      }
    }),

  removeDishItem: (id) =>
    set((state) => ({
      dishItems: state.dishItems.filter((i) => i.id !== id)
    })),

  updateDishQuantity: (id, quantity) =>
    set((state) => ({
      dishItems: state.dishItems.map((i) =>
        i.id === id ? { ...i, quantity } : i
      )
    })),

  addSetItem: (
    setItem: SetInterface,
    quantity: number,
    modifiedDishes?: SetProtoDish[]
  ) => {
    set((state: OrderState) => {
      const existingItem = state.setItems.find(
        (item) => item.set.id === setItem.id
      );
      if (existingItem) {
        return {
          setItems: state.setItems.map((item) =>
            item.set.id === setItem.id
              ? { ...item, quantity: item.quantity + quantity }
              : item
          )
        };
      } else {
        const newItem: SetOrderItem = {
          id: setItem.id,
          quantity,
          set: setItem,
          modifiedDishes: modifiedDishes || setItem.dishes
        };
        return { setItems: [...state.setItems, newItem] };
      }
    });
  },

  removeSetItem: (id) =>
    set((state) => ({
      setItems: state.setItems.filter((i) => i.id !== id)
    })),

  updateSetQuantity: (id, quantity) =>
    set((state) => ({
      setItems: state.setItems.map((i) =>
        i.id === id ? { ...i, quantity } : i
      )
    })),

  updateSetDishes: (setId, modifiedDishes) =>
    set((state) => ({
      setItems: state.setItems.map((i) =>
        i.id === setId ? { ...i, modifiedDishes } : i
      )
    })),

  clearOrder: () => set({ dishItems: [], setItems: [] }),

  setLoading: (isLoading) => set({ isLoading }),

  setError: (error) => set({ error }),

  getOrderSummary: () => {
    const { dishItems, setItems } = get();
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

    dishItems.forEach((item) => {
      totalItems += item.quantity;
      totalPrice += item.dish.price * item.quantity;
      dishes.push({
        id: item.dish.id,
        name: item.dish.name,
        quantity: item.quantity,
        price: item.dish.price
      });
    });

    setItems.forEach((item) => {
      totalItems += item.quantity;
      const setPrice = calculateSetPrice(item.modifiedDishes);
      totalPrice += setPrice * item.quantity;
      sets.push({
        id: item.set.id,
        name: item.set.name,
        quantity: item.quantity,
        price: setPrice,
        dishes: item.modifiedDishes
      });
    });

    return { totalItems, totalPrice, dishes, sets };
  },

  findDishOrderItem: (id) => get().dishItems.find((item) => item.id === id),
  findSetOrderItem: (id) => get().setItems.find((item) => item.id === id)
}));

function calculateSetPrice(dishes: SetProtoDish[] | undefined): number {
  if (!dishes || dishes.length === 0) {
    return 0;
  }
  return dishes.reduce((acc, d) => {
    if (typeof d.price === "number" && typeof d.quantity === "number") {
      return acc + d.price * d.quantity;
    }
    return acc;
  }, 0);
}

export default useOrderStore;
