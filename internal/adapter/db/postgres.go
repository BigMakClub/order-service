package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"order-service/order-service/internal/domain"
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

	const orderSQL = `SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, 
                      delivery_service, shardkey, sm_id, date_created 
                       FROM orders WHERE order_id=$1`

	err := p.pool.QueryRow(ctx, orderSQL, id).Scan(&order.OrderId, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.ShardKey, &order.SmId, &order.DateCreated)

	if err != nil {
		return nil, err
	}

	const delSQL = `SELECT del_name, phone, zip, city, address, region, email 
                    FROM deliveries WHERE order_id=$1`

	err = p.pool.QueryRow(ctx, delSQL, id).Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		return nil, err
	}

	const paySQL = `SELECT transaction_id, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee 
                    FROM payments WHERE order_id=$1`

	err = p.pool.QueryRow(ctx, paySQL, id).Scan(&order.Payment.TransactionId, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)

	if err != nil {
		return nil, err
	}

	const itemSQL = `SELECT item_id, chrt_id, track_number, price, rid, item_name, sale, item_size, nm_id, brand, status`

	rows, err := p.pool.Query(ctx, itemSQL, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.ChrtId, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.NmID, &item.Brand, &item.Status); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return &order, nil
}
