package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kavix/orders-api-golang/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

var ErrNotExist = errors.New("Order Does not EXIST")


func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to Encode Error: %w", err) // Fix 2
	}

	key := orderIDKey(order.OrderID)
	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0) // Fix 1: was r.Client.SetNX
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("Failed to Set %w", err)
	}

	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("Failed to add order to orders set %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		txn.Discard()
		return fmt.Errorf("Fail to Execute %w", err)
	}

	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)
	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("Get Order: %w", err)
	}

	var order model.Order
	if err = json.Unmarshal([]byte(value), &order); err != nil {
		return model.Order{}, fmt.Errorf("Failed to Decode: %w", err)
	}

	return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)
	err := r.Client.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("Delete Order: %w", err)
	}
	return nil
}

// ✅ Fix 3: return the marshal error in Update
func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to Encode Order: %w", err) // Fix 3: was missing return
	}

	key := orderIDKey(order.OrderID)
	res := r.Client.SetXX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		return fmt.Errorf("Failed to Update: %w", err)
	}

	return nil
}

type FindAllPage struct {
	Size   uint
	Offset uint
}

type findResult struct {
	Orders []model.Order
	Cursor uint64
}

// ✅ Fix 4: typo "cursoe" → "cursor"
// ✅ Fix 5: fetch actual order data from keys and return the result
func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (findResult, error) {
	res := r.Client.SScan(ctx, "orders", uint64(page.Offset), "*", int64(page.Size))

	keys, cursor, err := res.Result() // Fix 4: cursoe → cursor
	if err != nil {
		return findResult{}, fmt.Errorf("Fail to Get Order Ids: %w", err)
	}

	// Fix 5: function was incomplete — fetch each order by key
	if len(keys) == 0 {
		return findResult{Orders: []model.Order{}, Cursor: cursor}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return findResult{}, fmt.Errorf("Failed to Get Orders: %w", err)
	}

	orders := make([]model.Order, len(xs))
	for i, x := range xs {
		x := x.(string)
		var order model.Order
		if err = json.Unmarshal([]byte(x), &order); err != nil {
			return findResult{}, fmt.Errorf("Failed to Decode Order: %w", err)
		}
		orders[i] = order
	}

	return findResult{Orders: orders, Cursor: cursor}, nil
}