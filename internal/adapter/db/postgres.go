package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"order-service/internal/domain"
)

type PgRepo struct {
	pool *pgxpool.Pool
}

func NewPgRepo(pgURL string) (*PgRepo, error) {
	pool, err := pgxpool.New(context.Background(), pgURL)
	if err != nil {
		return nil, err
	}
	return &PgRepo{pool: pool}, nil
}

func (p *PgRepo) Find(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order

	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid UUID %s: %v", id, err)
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	log.Printf("Looking for order %s in database", uuid)

	const orderSQL = `SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, 
                  delivery_service, shardkey, sm_id, date_created 
                  FROM orders WHERE order_uid=$1`

	err = p.pool.QueryRow(ctx, orderSQL, uuid).Scan(&order.OrderId, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.ShardKey, &order.SmId, &order.DateCreated)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Order %s not found in database", uuid)
			return nil, nil
		}
		log.Printf("Error querying order %s: %v", uuid, err)
		return nil, err
	}

	log.Printf("Found order %s, fetching delivery info", uuid)

	const delSQL = `SELECT del_name, phone, zip, city, address, region, email
	           FROM deliveries WHERE order_uid=$1`

	err = p.pool.QueryRow(ctx, delSQL, uuid).Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		log.Printf("Error querying delivery for order %s: %v", uuid, err)
		return nil, err
	}

	log.Printf("Found delivery info for %s, fetching payment info", uuid)

	const paySQL = `SELECT transaction_id, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	           FROM payments WHERE order_uid=$1`

	err = p.pool.QueryRow(ctx, paySQL, uuid).Scan(&order.Payment.TransactionId, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)

	if err != nil {
		log.Printf("Error querying payment for order %s: %v", uuid, err)
		return nil, err
	}

	log.Printf("Found payment info for %s, fetching items", uuid)

	const itemSQL = `SELECT chrt_id, track_number, price, rid, item_name, sale, item_size, total_price, nm_id, brand, status
	            FROM items WHERE order_uid=$1`

	rows, err := p.pool.Query(ctx, itemSQL, uuid)

	if err != nil {
		log.Printf("Error querying items for order %s: %v", uuid, err)
		return nil, err
	}
	defer rows.Close()
	
	order.Items = make([]domain.Item, 0)
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status); err != nil {
			log.Printf("Error scanning item for order %s: %v", uuid, err)
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	log.Printf("Successfully loaded order %s with %d items", uuid, len(order.Items))
	return &order, nil
}

func (p *PgRepo) Save(ctx context.Context, order *domain.Order) error {

	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	const orderSQL = `INSERT INTO orders( order_uid, track_number, entry, locale,internal_signature, 
                      customer_id,delivery_service, shardkey, sm_id, date_created)
                      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10)`

	_, err = tx.Exec(ctx, orderSQL, order.OrderId, order.TrackNumber,
		order.Entry, order.Locale, order.InternalSignature, order.CustomerId,
		order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated)

	if err != nil {
		return err
	}

	const delSQL = `INSERT INTO deliveries(order_uid, del_name, phone, zip, city, address, region, email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(ctx, delSQL, order.OrderId, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email)

	if err != nil {
		return err
	}

	const paySQL = `INSERT INTO payments(order_uid,transaction_id, request_id, currency, provider, 
                    amount,payment_dt, bank, delivery_cost, goods_total, custom_fee)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,$11)`

	_, err = tx.Exec(ctx, paySQL, order.OrderId, order.Payment.TransactionId, order.Payment.RequestId,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)

	if err != nil {
		return err
	}

	const itemSQL = `INSERT INTO items(order_uid, chrt_id, track_number, price, rid, 
                     item_name, sale, item_size, total_price,nm_id, brand, status) 
 					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	batch := &pgx.Batch{}

	for _, item := range order.Items {
		batch.Queue(itemSQL, order.OrderId, item.ChrtId, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
	}

	br := tx.SendBatch(ctx, batch)

	if err = br.Close(); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *PgRepo) CacheRestore(ctx context.Context) ([]*domain.Order, error) {
	result := make([]*domain.Order, 0, 10)
	const orderSQL = `SELECT order_uid FROM orders ORDER BY date_created DESC LIMIT 10 `

	rows, err := p.pool.Query(ctx, orderSQL)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderUID string
		if err = rows.Scan(&orderUID); err != nil {
			return nil, err
		}

		order, err := p.Find(ctx, orderUID)
		if err != nil {
			return nil, err
		}

		result = append(result, order)
	}

	return result, nil

	return nil, nil

}
