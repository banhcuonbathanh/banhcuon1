// "use client";
// import React from "react";
// import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
// import { Order } from "@/schemaValidations/interface/type_order";
// import useOrderStore from "@/zusstand/order/order_zustand";

// interface OrderCardProps {
//   order: Order;
//   className?: string;
// }

// // Status styles defined outside component to prevent recreation
// const STATUS_STYLES = {
//   pending: "bg-yellow-100 text-yellow-800",
//   completed: "bg-green-100 text-green-800",
//   default: "bg-gray-100 text-gray-800"
// } as const;

// // Separate components for each section to prevent unnecessary rerenders
// const OrderDetails = React.memo(({ order }: { order: Order }) => (
//   <div className="space-y-1 text-sm">
//     <p>Order Name: {order.order_name}</p>
//     <p>Table Number: {order.table_number}</p>
//     <p>Take Away: {order.takeAway ? "Yes" : "No"}</p>
//   </div>
// ));

// const DishList = React.memo(({ id }: { id: number }) => {
//   // Use a stable selector function
//   const dishItems = useOrderStore(
//     React.useCallback(
//       (state) =>
//         state.listOfOrders.find((order) => order.id === id)?.dish_items ?? [],
//       [id]
//     )
//   );

//   return (
//     <div className="space-y-2">
//       {dishItems.map((dish) => (
//         <div
//           key={dish.dish_id}
//           className="flex justify-between items-center text-sm"
//         >
//           <span>
//             {dish.name || `Dish #${dish.dish_id}`} x{dish.quantity}
//           </span>
//         </div>
//       ))}
//     </div>
//   );
// });

// const SetList = React.memo(({ id }: { id: number }) => {
//   const sets = useOrderStore(
//     (state) =>
//       state.listOfOrders.find((order) => order.id === id)?.set_items ?? []
//   );

//   return (
//     <div className="space-y-2">
//       {sets.map((set) => (
//         <div key={set.set_id} className="text-sm">
//           {/* <div className="flex justify-between items-center">
//           <span>
//             {set.name} x{set.quantity}
//           </span>
//         </div> */}
//           <div className="ml-4 mt-1 text-gray-600">
//             {set.dishes.map((dish) => (
//               <div key={dish.dish_id} className="text-xs">
//                 - {dish.name} x{dish.quantity}
//               </div>
//             ))}
//           </div>
//         </div>
//       ))}
//     </div>
//   );
// });

// const OrderCard = ({ order, className = "" }: OrderCardProps) => {
//   // Get status style without computation
//   const statusStyle =
//     STATUS_STYLES[order.status as keyof typeof STATUS_STYLES] ||
//     STATUS_STYLES.default;

//   return (
//     <Card className={`w-full ${className}`}>
//       <CardHeader>
//         <CardTitle className="flex justify-between items-center">
//           <span>Order #{order.id}</span>
//           <span className={`px-3 py-1 rounded-full text-sm ${statusStyle}`}>
//             {order.status}
//           </span>
//         </CardTitle>
//       </CardHeader>

//       <CardContent>
//         <div className="grid md:grid-cols-2 gap-4">
//           <div>
//             <h3 className="font-semibold mb-2">Order Details</h3>
//             <OrderDetails order={order} />
//           </div>

//           <div>
//             <div className="mb-4">
//               <h3 className="font-semibold mb-2">Dishes</h3>
//               <DishList id={order.id} />
//             </div>

//             <div>
//               <h3 className="font-semibold mb-2">Sets</h3>
//               {/* <SetList sets={order.set_items} /> */}
//             </div>
//           </div>
//         </div>

//         {order.topping && (
//           <div className="mt-4 text-sm">
//             <h3 className="font-semibold mb-1">Additional Notes</h3>
//             {/* <p>{order.topping}</p> */}
//           </div>
//         )}
//       </CardContent>
//     </Card>
//   );
// };

// export default React.memo(OrderCard);
