import { create } from "zustand";

import { useApiStore } from "@/zusstand/api/api-controller";

import {
  CreateDishBodyType,
  Dish,
  DishListResType,
  DishParamsType,
  DishRes,
  DishResType,
  DishSchema,
  UpdateDishBodyType
} from "../domain/dish.schema";
import envConfig from "@/config";

interface DishStore {
  dish: Dish | null;
  dishes: Dish[];
  isLoading: boolean;
  error: string | null;
  getDish: (id: number) => Promise<void>;
  updateDish: (body: UpdateDishBodyType & DishParamsType) => Promise<void>;
  getDishes: () => Promise<void>;
  addDish: (body: CreateDishBodyType) => Promise<void>;
  deleteDish: (id: number) => Promise<void>;
}

export const useDishStore = create<DishStore>((set) => ({
  dish: null,
  dishes: [],
  isLoading: false,
  error: null,
  getDish: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      const response = await useApiStore
        .getState()
        .http.get<DishResType>(`/api/dishes/${id}`);
      set({ dish: response.data.data, isLoading: false });
    } catch (error) {
      set({ isLoading: false, error: "Failed to fetch dish" });
      throw error;
    }
  },
  updateDish: async (body: UpdateDishBodyType & DishParamsType) => {
    set({ isLoading: true, error: null });
    try {
      const response = await useApiStore
        .getState()
        .http.put<DishResType>(`/api/dishes/${body.id}`, body);
      set({ dish: response.data.data, isLoading: false });
    } catch (error) {
      set({ isLoading: false, error: "Failed to update dish" });
      throw error;
    }
  },
  getDishes: async () => {
    set({ isLoading: true, error: null });
    try {
      const response = await useApiStore
        .getState()
        .http.get<DishListResType>("/api/dishes");
      set({ dishes: response.data.data, isLoading: false });
    } catch (error) {
      set({ isLoading: false, error: "Failed to fetch dishes" });
      throw error;
    }
  },
  addDish: async (body: CreateDishBodyType) => {

    const link = envConfig.NEXT_PUBLIC_API_ENDPOINT + envConfig.NEXT_PUBLIC_Add_Dished
    set({ isLoading: true, error: null });
    try {
      const response = await useApiStore
        .getState()
        .http.post<DishResType>(link, body);
      set((state) => ({
        dishes: [...state.dishes, response.data.data],
        isLoading: false
      }));
    } catch (error) {
      set({ isLoading: false, error: "Failed to add dish" });
      throw error;
    }
  },
  deleteDish: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      await useApiStore
        .getState()
        .http.delete<DishResType>(`/api/dishes/${id}`);
      set((state) => ({
        dishes: state.dishes.filter((dish) => dish.id !== id),
        isLoading: false
      }));
    } catch (error) {
      set({ isLoading: false, error: "Failed to delete dish" });
      throw error;
    }
  }
}));

// Custom hooks for each operation
export const useAddDishMutation = () => {
  const { addDish, isLoading, error } = useDishStore();
  return {
    mutateAsync: addDish,
    isPending: isLoading,
    error
  };
};

export const useDeleteDishMutation = () => {
  const { deleteDish, isLoading, error } = useDishStore();
  return {
    mutateAsync: deleteDish,
    isPending: isLoading,
    error
  };
};

export const useDishListQuery = () => {
  const { getDishes, dishes, isLoading, error } = useDishStore();
  return {
    refetch: getDishes,
    data: dishes,
    isLoading,
    error
  };
};

export const useGetDishQuery = () => {
  const { getDish, dish, isLoading, error } = useDishStore();
  return {
    refetch: getDish,
    data: dish,
    isLoading,
    error
  };
};

export const useUpdateDishMutation = () => {
  const { updateDish, isLoading, error } = useDishStore();
  return {
    mutateAsync: updateDish,
    isPending: isLoading,
    error
  };
};



// const { mutateAsync: addDish, isPending: isAdding, error: addError } = useAddDishMutation();
// const { mutateAsync: deleteDish, isPending: isDeleting, error: deleteError } = useDeleteDishMutation();
// const { refetch: fetchDishes, data: dishes, isLoading: isLoadingDishes, error: dishesError } = useDishListQuery();
// const { refetch: fetchDish, data: dish, isLoading: isLoadingDish, error: dishError } = useGetDishQuery();
// const { mutateAsync: updateDish, isPending: isUpdating, error: updateError } = useUpdateDishMutation();