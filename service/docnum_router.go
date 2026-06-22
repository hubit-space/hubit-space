package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client membungkus koneksi ke database mongo milik service.
// Service hanya menyimpan log/activity timeline di sini, bukan data transaksional.
type Client struct {
	DB *mongo.Database
}

// Connect membuka koneksi ke MongoDB dan memverifikasi dengan ping.
func Connect(uri, dbName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Client{DB: client.Database(dbName)}, nil
}

// ActivityEvent adalah satu entri di timeline status sebuah order.
type ActivityEvent struct {
	Status    string    `bson:"status" json:"status"`
	Note      string    `bson:"note" json:"note"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

// LogActivity menambahkan satu event baru ke timeline order tertentu.
// Dokumen "order_activity_logs" otomatis dibuat kalau belum ada (upsert).
func (c *Client) LogActivity(ctx context.Context, orderID, status, note string) error {
	collection := c.DB.Collection("order_activity_logs")

	event := ActivityEvent{
		Status:    status,
		Note:      note,
		Timestamp: time.Now(),
	}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"order_id": orderID},
		bson.M{"$push": bson.M{"events": event}},
		options.Update().SetUpsert(true),
	)
	return err
}
