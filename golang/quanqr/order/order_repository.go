package order_grpc

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	"english-ai-full/quanqr/proto_qr/order"
)


type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (or *OrderRepository) CreateOrders(ctx context.Context, req *order.CreateOrderRequest) ([]*order.Order, error) {
	// Start a transaction

	log.Print("golang/quanqr/order/order_repository.go CreateOrders")
	tx, err := or.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert the main order - removed dish_snapshot_id from the query
	query := `
		INSERT INTO orders (guest_id, table_number, order_handler_id, status, created_at, updated_at, total_price)
		VALUES ($1, $2, $3, $4, $5, $5, $6)
		RETURNING id, created_at, updated_at
	`
	var o order.Order
	var createdAt, updatedAt time.Time
	err = tx.QueryRow(ctx, query,
		req.GuestId,
		req.TableNumber,
		req.OrderHandlerId,
		req.Status,
		time.Now(),
		req.TotalPrice,
	).Scan(&o.Id, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	// Insert dish items - Updated to use dish_snapshot_id instead of dish_id
	for _, item := range req.DishItems {
		_, err := tx.Exec(ctx, `
			INSERT INTO dish_order_items (order_id, dish_snapshot_id, quantity)
			VALUES ($1, $2, $3)
		`, o.Id, item.Dish.Id, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error inserting dish order item: %w", err)
		}
	}

	// Insert set items
	for _, item := range req.SetItems {
		_, err := tx.Exec(ctx, `
			INSERT INTO set_order_items (order_id, set_snapshot_id, quantity)
			VALUES ($1, $2, $3)
		`, o.Id, item.Set.Id, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error inserting set order item: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	// Populate the order struct with the inserted data
	o.GuestId = req.GuestId
	o.TableNumber = req.TableNumber
	o.OrderHandlerId = req.OrderHandlerId
	o.Status = req.Status
	o.TotalPrice = req.TotalPrice
	o.CreatedAt = timestamppb.New(createdAt)
	o.UpdatedAt = timestamppb.New(updatedAt)
	o.DishItems = req.DishItems
	o.SetItems = req.SetItems

	return []*order.Order{&o}, nil
}
func (or *OrderRepository) GetOrders(ctx context.Context, req *order.GetOrdersRequest) ([]*order.Order, error) {
	query := `
		SELECT id, guest_id, table_number, dish_snapshot_id, order_handler_id, status, created_at, updated_at, total_price 
		FROM orders
		WHERE created_at BETWEEN $1 AND $2
	`
	rows, err := or.db.Query(ctx, query, req.FromDate.AsTime(), req.ToDate.AsTime())
	if err != nil {
		return nil, fmt.Errorf("error fetching orders: %w", err)
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&o.Id,
			&o.GuestId,
			&o.TableNumber,
			&o.DishSnapshotId,
			&o.OrderHandlerId,
			&o.Status,
			&createdAt,
			&updatedAt,
			&o.TotalPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		o.CreatedAt = timestamppb.New(createdAt)
		o.UpdatedAt = timestamppb.New(updatedAt)

		// Fetch dish items
		o.DishItems, err = or.getDishItems(ctx, o.Id)
		if err != nil {
			return nil, err
		}

		// Fetch set items
		o.SetItems, err = or.getSetItems(ctx, o.Id)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over orders: %w", err)
	}

	return orders, nil
}

func (or *OrderRepository) GetOrderDetail(ctx context.Context, id int64) (*order.Order, error) {
	query := `
		SELECT id, guest_id, table_number, dish_snapshot_id, order_handler_id, status, created_at, updated_at, total_price 
		FROM orders 
		WHERE id = $1
	`
	var o order.Order
	var createdAt, updatedAt time.Time
	err := or.db.QueryRow(ctx, query, id).Scan(
		&o.Id,
		&o.GuestId,
		&o.TableNumber,
		&o.DishSnapshotId,
		&o.OrderHandlerId,
		&o.Status,
		&createdAt,
		&updatedAt,
		&o.TotalPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching order detail: %w", err)
	}
	o.CreatedAt = timestamppb.New(createdAt)
	o.UpdatedAt = timestamppb.New(updatedAt)

	// Fetch dish items
	o.DishItems, err = or.getDishItems(ctx, o.Id)
	if err != nil {
		return nil, err
	}

	// Fetch set items
	o.SetItems, err = or.getSetItems(ctx, o.Id)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (or *OrderRepository) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.Order, error) {
	// Start a transaction
	tx, err := or.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update the main order
	query := `
		UPDATE orders
		SET guest_id = $2, table_number = $3, dish_snapshot_id = $4, order_handler_id = $5, status = $6, updated_at = $7, total_price = $8
		WHERE id = $1
		RETURNING created_at, updated_at
	`
	var createdAt, updatedAt time.Time
	err = tx.QueryRow(ctx, query,
		req.Id,
		req.GuestId,
		req.TableNumber,
		req.DishSnapshotId,
		req.OrderHandlerId,
		req.Status,
		time.Now(),
		req.TotalPrice,
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	// Delete existing dish items and set items
	_, err = tx.Exec(ctx, "DELETE FROM dish_order_items WHERE order_id = $1", req.Id)
	if err != nil {
		return nil, fmt.Errorf("error deleting existing dish order items: %w", err)
	}
	_, err = tx.Exec(ctx, "DELETE FROM set_order_items WHERE order_id = $1", req.Id)
	if err != nil {
		return nil, fmt.Errorf("error deleting existing set order items: %w", err)
	}

	// Insert updated dish items
	for _, item := range req.DishItems {
		_, err := tx.Exec(ctx, `
			INSERT INTO dish_order_items (order_id, dish_id, quantity)
			VALUES ($1, $2, $3)
		`, req.Id, item.Dish.Id, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error inserting updated dish order item: %w", err)
		}
	}

	// Insert updated set items
	for _, item := range req.SetItems {
		setItemId := int64(0)
		err := tx.QueryRow(ctx, `
			INSERT INTO set_order_items (order_id, set_id, quantity)
			VALUES ($1, $2, $3)
			RETURNING id
		`, req.Id, item.Set.Id, item.Quantity).Scan(&setItemId)
		if err != nil {
			return nil, fmt.Errorf("error inserting updated set order item: %w", err)
		}

		// Insert modified dishes for set items
		for _, modifiedDish := range item.ModifiedDishes {
			_, err := tx.Exec(ctx, `
				INSERT INTO set_order_item_modified_dishes (set_order_item_id, dish_id, name, price)
				VALUES ($1, $2, $3, $4)
			`, setItemId, modifiedDish.Id, modifiedDish.Name, modifiedDish.Price)
			if err != nil {
				return nil, fmt.Errorf("error inserting modified dish for updated set order item: %w", err)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &order.Order{
		Id:              req.Id,
		GuestId:         req.GuestId,
		TableNumber:     req.TableNumber,
		DishSnapshotId:  req.DishSnapshotId,
		OrderHandlerId:  req.OrderHandlerId,
		Status:          req.Status,
		CreatedAt:       timestamppb.New(createdAt),
		UpdatedAt:       timestamppb.New(updatedAt),
		TotalPrice:      req.TotalPrice,
		DishItems:       req.DishItems,
		SetItems:        req.SetItems,
	}, nil
}

func (or *OrderRepository) PayGuestOrders(ctx context.Context, guestId int64) ([]*order.Order, error) {
	// Start a transaction
	tx, err := or.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update the status of all unpaid orders for the guest
	query := `
		UPDATE orders
		SET status = 'PAID', updated_at = $2
		WHERE guest_id = $1 AND status != 'PAID'
		RETURNING id, guest_id, table_number, dish_snapshot_id, order_handler_id, status, created_at, updated_at, total_price
	`
	rows, err := tx.Query(ctx, query, guestId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error updating guest orders: %w", err)
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&o.Id,
			&o.GuestId,
			&o.TableNumber,
			&o.DishSnapshotId,
			&o.OrderHandlerId,
			&o.Status,
			&createdAt,
			&updatedAt,
			&o.TotalPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning paid order: %w", err)
		}
		o.CreatedAt = timestamppb.New(createdAt)
		o.UpdatedAt = timestamppb.New(updatedAt)

		// Fetch dish items and set items (these methods should be implemented)
		o.DishItems, err = or.getDishItems(ctx, o.Id)
		if err != nil {
			return nil, err
		}
		o.SetItems, err = or.getSetItems(ctx, o.Id)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over paid orders: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return orders, nil
}


func (or *OrderRepository) getDishItems(ctx context.Context, orderId int64) ([]*order.DishOrderItem, error) {
	query := `
		SELECT doi.id, doi.quantity, d.id, d.name, d.price, d.description, d.image, d.status, d.created_at, d.updated_at
		FROM dish_order_items doi
		JOIN dishes d ON doi.dish_id = d.id
		WHERE doi.order_id = $1
	`
	rows, err := or.db.Query(ctx, query, orderId)
	if err != nil {
		return nil, fmt.Errorf("error fetching dish items: %w", err)
	}
	defer rows.Close()

	var items []*order.DishOrderItem
	for rows.Next() {
		var item order.DishOrderItem
		var dish order.DishOrder
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&item.Id,
			&item.Quantity,
			&dish.Id,
			&dish.Name,
			&dish.Price,
			&dish.Description,
			&dish.Image,
			&dish.Status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning dish item: %w", err)
		}
		dish.CreatedAt = timestamppb.New(createdAt)
		dish.UpdatedAt = timestamppb.New(updatedAt)
		item.Dish = &dish
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over dish items: %w", err)
	}

	return items, nil
}


func (or *OrderRepository) getSetItems(ctx context.Context, orderId int64) ([]*order.SetOrderItem, error) {
	query := `
		SELECT soi.id, soi.quantity, 
			s.id, s.name, s.description, s.created_at, s.updated_at, s.is_favourite, s.is_public, s.image,
			spd.id, spd.name, spd.price,
			soimd.id, soimd.name, soimd.price
		FROM set_order_items soi
		JOIN sets s ON soi.set_id = s.id
		LEFT JOIN set_proto_dishes spd ON s.id = spd.set_id
		LEFT JOIN set_order_item_modified_dishes soimd ON soi.id = soimd.set_order_item_id
		WHERE soi.order_id = $1
	`
	rows, err := or.db.Query(ctx, query, orderId)
	if err != nil {
		return nil, fmt.Errorf("error fetching set items: %w", err)
	}
	defer rows.Close()

	setItemMap := make(map[int64]*order.SetOrderItem)
	for rows.Next() {
		var setItem order.SetOrderItem
		var set order.SetProto
		// var setDish order.SetProtoDish
		// var modifiedDish order.SetProtoDish
		var setCreatedAt, setUpdatedAt time.Time
		var setDishId, modifiedDishId sql.NullInt64
		var setDishName, modifiedDishName sql.NullString
		var setDishPrice, modifiedDishPrice sql.NullInt32

		err := rows.Scan(
			&setItem.Id, &setItem.Quantity,
			&set.Id, &set.Name, &set.Description, &setCreatedAt, &setUpdatedAt, &set.IsFavourite, &set.IsPublic, &set.Image,
			&setDishId, &setDishName, &setDishPrice,
			&modifiedDishId, &modifiedDishName, &modifiedDishPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning set item: %w", err)
		}

		set.CreatedAt = timestamppb.New(setCreatedAt)
		set.UpdatedAt = timestamppb.New(setUpdatedAt)

		if existingSetItem, ok := setItemMap[setItem.Id]; ok {
			// If the set item already exists, we just need to add the dishes
			if setDishId.Valid {
				existingSetItem.Set.Dishes = append(existingSetItem.Set.Dishes, &order.SetProtoDish{
					Id:    setDishId.Int64,
					Name:  setDishName.String,
					Price: setDishPrice.Int32,
				})
			}
			if modifiedDishId.Valid {
				existingSetItem.ModifiedDishes = append(existingSetItem.ModifiedDishes, &order.SetProtoDish{
					Id:    modifiedDishId.Int64,
					Name:  modifiedDishName.String,
					Price: modifiedDishPrice.Int32,
				})
			}
		} else {
			// If it's a new set item, we create it with the current set and dishes
			setItem.Set = &set
			if setDishId.Valid {
				setItem.Set.Dishes = []*order.SetProtoDish{{
					Id:    setDishId.Int64,
					Name:  setDishName.String,
					Price: setDishPrice.Int32,
				}}
			}
			if modifiedDishId.Valid {
				setItem.ModifiedDishes = []*order.SetProtoDish{{
					Id:    modifiedDishId.Int64,
					Name:  modifiedDishName.String,
					Price: modifiedDishPrice.Int32,
				}}
			}
			setItemMap[setItem.Id] = &setItem
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over set items: %w", err)
	}

	// Convert the map to a slice
	setItems := make([]*order.SetOrderItem, 0, len(setItemMap))
	for _, item := range setItemMap {
		setItems = append(setItems, item)
	}

	return setItems, nil
}