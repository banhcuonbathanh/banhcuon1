import DishList from "./component/dishlish";
import { get_dishes } from "./controller/guest-controller";
import { DishListResType, Dish } from "@/zusstand/dished/domain/dish.schema";


// This is a server component
export default async function GuestPage() {
  // Fetch dishes on the server side
  const dishesData: DishListResType = await get_dishes();
  const dishes: Dish[] = dishesData.data;

  return (
    <div className="guest-page">
      <h1>Menu</h1>
      {/* Pass the fetched dishes to the client-side DishList component */}
      <DishList dishes={dishes} />
    </div>
  );
}
