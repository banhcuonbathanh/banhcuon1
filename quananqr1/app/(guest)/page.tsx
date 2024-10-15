import { get_dishes } from "@/zusstand/server/dish-controller";
import { DishSelection } from "./component/dish/dishh_list";

import { DishInterface } from "@/schemaValidations/interface/type_dish";
import { SetInterface } from "@/schemaValidations/interface/types_set";
import { get_Sets } from "@/zusstand/server/set-controller";
import { SetCardList } from "./component/set/sets_list";
import OrderSummary from "./component/order/order";

// This is a server component
export default async function GuestPage() {
  // Fetch dishes on the server side
  const dishesData: DishInterface[] = await get_dishes();

  const setsData: SetInterface[] = await get_Sets();
  // const dishes: Dish[] = dishesData;
  console.log("quananqr1/app/guest/page.tsx dishes.data asdf", dishesData);
  return (
    <div className="guest-page">
      <DishSelection dishes={dishesData} />

      <SetCardList sets={setsData} />

      <OrderSummary />
    </div>
  );
}
