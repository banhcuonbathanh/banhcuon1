import { Dish } from "@/zusstand/dished/domain/dish.schema";

const DishItem = ({
  dish,
  handleAddToOrder
}: {
  dish: Dish;
  handleAddToOrder: (dish: Dish) => void;
}) => {
  return (
    <div className="dish-item">
      <img src={dish.image} alt={dish.name} width={100} height={100} />
      <div>
        <h3>{dish.name}</h3>
        <p>{dish.description}</p>
        <p>Price: ${dish.price}</p>
        <button onClick={() => handleAddToOrder(dish)}>Add to Order</button>
      </div>
    </div>
  );
};

export default DishItem;
