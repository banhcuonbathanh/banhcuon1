{
"id": 5,
"guest_id": 0,
"user_id": 1,
"table_number": 1,
"order_handler_id": 1,
"status": "pending",
"total_price": 5000,
"is_guest": false,
"topping": "xzdfgdfg",
"tracking_order": "1st Client - Order #1",
"take_away": true,
"chili_number": 3,
"table_token": "MTo0OkF2YWlsYWJsZTo0ODg5ODMyNTYy.ZiOccA4-JoM",
"order_name": "asdfasdf-344c1368-beaa-4bb2-8649-",
"current_version": 8,
"parent_order_id": 0,
"data_set": [
{
"id": 1,
"name": "day du chin ",
"description": "A delicious set of dishes.",
"dishes": [
{
"dish_id": 1,
"quantity": 3,
"name": "banh",
"price": 9,
"description": "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper",
"image": "https://example.com/spaghetti-carbonara.jpg",
"status": ""
},
{
"dish_id": 2,
"quantity": 1,
"name": "trung",
"price": 9,
"description": "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper",
"image": "https://example.com/spaghetti-carbonara.jpg",
"status": ""
},
{
"dish_id": 3,
"quantity": 1,
"name": "gio",
"price": 9,
"description": "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper",
"image": "https://example.com/spaghetti-carbonara.jpg",
"status": ""
}
],
"userId": 1,
"created_at": "2025-01-08T23:28:29.881134Z",
"updated_at": "2025-01-08T23:28:29.881134Z",
"is_favourite": false,
"like_by": null,
"is_public": true,
"image": "http://example.com/image.jpg",
"price": 45,
"quantity": 4
}
],
"data_dish": [
{
"dish_id": 3,
"quantity": 4,
"name": "gio",
"price": 9,
"description": "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper",
"image": "https://example.com/spaghetti-carbonara.jpg",
"status": "available"
},
{
"dish_id": 1,
"quantity": 2,
"name": "banh",
"price": 9,
"description": "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper",
"image": "https://example.com/spaghetti-carbonara.jpg",
"status": "available"
},
{
"dish_id": 2,
"quantity": 2,
"name": "trung",
"price": 9,
"description": "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper",
"image": "https://example.com/spaghetti-carbonara.jpg",
"status": "available"
}
],
"version_history": [
{
"version_number": 1,
"total_dishes_count": 3,
"total_sets_count": 2,
"version_total_price": 342,
"modification_type": "INITIAL",
"modified_at": "2025-01-09T07:29:56.963865Z",
"changes": [
{
"item_type": "DISH",
"item_id": 1,
"item_name": "banh",
"quantity_changed": 2,
"price": 9
}
]
}
],
"total_summary": {
"total_versions": 8,
"total_dishes_ordered": 36,
"total_sets_ordered": 55,
"cumulative_total_price": 2799,
"most_ordered_items": [
{
"item_type": "SET",
"item_id": 1,
"item_name": "day du chin ",
"total_quantity": 31
}
]
},
"delivery_history": [
{
"id": 1,
"order_id": 5,
"order_name": "asdfasdf-344c1368-beaa-4bb2-8649-",
"guest_id": 0,
"user_id": 1,
"table_number": 1,
"dish_items": [
{
"dish_id": 1,
"quantity": 2,
"name": "banh"
},
{
"dish_id": 3,
"quantity": 2,
"name": "gio"
}
],
"quantity_delivered": 4,
"delivery_status": "PARTIALLY_DELIVERED",
"delivered_at": "2025-01-26T03:00:00Z",
"delivered_by_user_id": 2,
"is_guest": false
},
{
"id": 2,
"order_id": 5,
"order_name": "asdfasdf-344c1368-beaa-4bb2-8649-",
"guest_id": 0,
"user_id": 1,
"table_number": 1,
"dish_items": [
{
"dish_id": 2,
"quantity": 2,
"name": "trung"
}
],
"quantity_delivered": 2,
"delivery_status": "FULLY_DELIVERED",
"delivered_at": "2025-01-26T03:30:00Z",
"delivered_by_user_id": 2,
"is_guest": false
}
],
"current_delivery_status": "PARTIALLY_DELIVERED",
"total_items_delivered": 6,
"last_delivery_at": "2025-01-26T03:30:00Z"
}
