import { get_dishes } from "@/zusstand/server/dish-controller";
import { DishSelection } from "./component/dish/dishh_list";

import { DishListResType, Dish } from "@/zusstand/dished/domain/dish.schema";

// This is a server component
export default async function GuestPage() {
  // Fetch dishes on the server side
  const dishesData: DishListResType = await get_dishes();
  const dishes: Dish[] = dishesData;
  console.log("quananqr1/app/guest/page.tsx dishes.data asdf", dishes);
  return (
    <div className="guest-page">
      <DishSelection dishes={dishes} />
    </div>
  );
}
