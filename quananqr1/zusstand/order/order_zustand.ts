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
// Keep existing interfaces
interface DishOrderItemustand extends DishInterface {
  quantity: number;
}

interface SetOrderItemustand extends SetInterface {
  quantity: number;
}
interface FormattedSetItem {
  id: number;
  name: string;
  displayString: string; // e.g. "My Set12 - $27.00 x 4"
  itemsString: string; // e.g. "banh x 1, trung x 2"
  totalPrice: number;
  formattedTotalPrice: string; // e.g. "$108.00"
}

interface FormattedDishItem {
  id: number;
  name: string;
  displayString: string; // e.g. "Dish do - $9.00 x 5"
  totalPrice: number;
  formattedTotalPrice: string; // e.g. "$45.00"
}
interface OrderState {
  dishItems: DishOrderItemustand[];
  setItems: SetOrderItemustand[];
  isLoading: boolean;
  error: string | null;
  tableNumber: number;
  tabletoken: string;

  addTableToken: (tableToken: string) => void;
  addTableNumber: (tableNumber: number) => void;
  addDishItem: (dish: DishInterface, quantity: number) => void;
  removeDishItem: (id: number) => void;
  updateDishQuantity: (id: number, quantity: number) => void;
  addSetItem: (set: SetInterface, quantity: number) => void;
  removeSetItem: (id: number) => void;
  updateSetQuantity: (id: number, quantity: number) => void;
  updateSetDishes: (setId: number, modifiedDishes: SetProtoDish[]) => void;
  clearOrder: () => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;

  getFormattedSets: () => FormattedSetItem[];
  getFormattedDishes: () => FormattedDishItem[];
  getFormattedTotals: () => {
    totalItems: number;
    formattedTotalPrice: string;
  };
  getOrderSummary: () => {
    totalItems: number;
    totalPrice: number;
    dishes: DishOrderItemustand[]; // Updated to match implementation
    sets: SetOrderItemustand[]; // Updated to match implementation
  };
  findDishOrderItem: (id: number) => DishOrderItemustand | undefined; // Updated return type
  findSetOrderItem: (id: number) => SetOrderItemustand | undefined; // Updated return type
}

const useOrderStore = create<OrderState>((set, get) => ({
  dishItems: [],
  setItems: [],
  isLoading: false,
  error: null,
  tableNumber: 0,

  tabletoken: "",
  addTableToken: (token: string) =>
    set(() => {
      console.log(
        "quananqr1/zusstand/order/order_zustand.ts addTableToken",
        token
      );
      return { tabletoken: token };
    }),
  addTableNumber: (tableNumber) =>
    set(() => {
      return { tableNumber: tableNumber };
    }),

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
        const newItem: DishOrderItemustand = {
          ...dish,
          quantity
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

  addSetItem: (setItem: SetInterface, quantity: number) => {
    set((state) => {
      const existingItem = state.setItems.find(
        (item) => item.id === setItem.id
      );
      if (existingItem) {
        return {
          setItems: state.setItems.map((item) =>
            item.id === setItem.id
              ? { ...item, quantity: item.quantity + quantity }
              : item
          )
        };
      } else {
        const newItem: SetOrderItemustand = {
          ...setItem,
          quantity
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
        i.id === setId ? { ...i, dishes: modifiedDishes } : i
      )
    })),

  clearOrder: () => set({ dishItems: [], setItems: [] }),

  setLoading: (isLoading) => set({ isLoading }),

  setError: (error) => set({ error }),

  getOrderSummary: () => {
    const { dishItems, setItems } = get();

    // Calculate total items counting quantities
    const totalItems =
      dishItems.reduce((acc, item) => acc + item.quantity, 0) +
      setItems.reduce((acc, item) => acc + item.quantity, 0);

    // Calculate total price including quantities
    const dishesPrice = dishItems.reduce(
      (acc, item) => acc + item.price * item.quantity,
      0
    );
    const setsPrice = setItems.reduce((acc, item) => {
      const setPrice = calculateSetPrice(item.dishes);
      return acc + setPrice * item.quantity;
    }, 0);

    const totalPrice = dishesPrice + setsPrice;

    return {
      totalItems,
      totalPrice,
      dishes: dishItems,
      sets: setItems
    };
  },

  findDishOrderItem: (id) => get().dishItems.find((item) => item.id === id),
  findSetOrderItem: (id) => get().setItems.find((item) => item.id === id),

  getFormattedSets: () => {
    const { setItems } = get();

    return setItems.map((set) => {
      const basePrice = calculateSetPrice(set.dishes);
      const totalPrice = basePrice * set.quantity;

      // Format display strings
      const displayString = `${set.name} - ${formatCurrency(basePrice)} x ${
        set.quantity
      }`;
      const itemsString = set.dishes
        .map((dish) => `${dish.name} x ${dish.quantity}`)
        .join(", ");

      return {
        id: set.id,
        name: set.name,
        displayString,
        itemsString,
        totalPrice,
        formattedTotalPrice: formatCurrency(totalPrice)
      };
    });
  },

  getFormattedDishes: () => {
    const { dishItems } = get();

    return dishItems.map((dish) => {
      const totalPrice = dish.price * dish.quantity;

      return {
        id: dish.id,
        name: dish.name,
        displayString: `${dish.name} - ${formatCurrency(dish.price)} x ${
          dish.quantity
        }`,
        totalPrice,
        formattedTotalPrice: formatCurrency(totalPrice)
      };
    });
  },

  getFormattedTotals: () => {
    const { dishItems, setItems } = get();

    const totalItems =
      dishItems.reduce((acc, item) => acc + item.quantity, 0) +
      setItems.reduce((acc, item) => acc + item.quantity, 0);

    const totalPrice =
      dishItems.reduce((acc, item) => acc + item.price * item.quantity, 0) +
      setItems.reduce(
        (acc, item) => acc + calculateSetPrice(item.dishes) * item.quantity,
        0
      );

    return {
      totalItems,
      formattedTotalPrice: formatCurrency(totalPrice)
    };
  }
}));

function calculateSetPrice(dishes: SetProtoDish[]): number {
  if (!dishes || dishes.length === 0) {
    return 0;
  }
  return dishes.reduce((acc, dish) => acc + dish.price * dish.quantity, 0);
}

function formatCurrency(amount: number): string {
  return `${amount}`;
}

export default useOrderStore;
