import { get_dishes } from "@/zusstand/server/dish-controller";
import { DishSelection } from "./component/dish/dishh_list";
import { Dish } from "@/schemaValidations/interface/type_dish";

// This is a server component
export default async function GuestPage() {
  // Fetch dishes on the server side
  const dishesData: Dish[] = await get_dishes();
  // const dishes: Dish[] = dishesData;
  console.log("quananqr1/app/guest/page.tsx dishes.data asdf", dishesData);
  return (
    <div className="guest-page">
      <DishSelection dishes={dishesData} />
    </div>
  );
}
