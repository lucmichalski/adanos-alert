package impl

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueueRepo struct {
	col *mongo.Collection
}

func NewQueueRepo(db *mongo.Database) repository.QueueRepo {
	return &QueueRepo{col: db.Collection("queue")}
}

func (q *QueueRepo) Enqueue(item repository.QueueItem) (id primitive.ObjectID, err error) {
	if item.ID.IsZero() {
		item.CreatedAt = time.Now()
		item.UpdatedAt = item.CreatedAt
		item.Status = repository.QueueItemStatusWait
		item.RequeueTimes = 0

		if item.NextExecuteAt.IsZero() {
			item.NextExecuteAt = item.CreatedAt
		}
		rs, err := q.col.InsertOne(context.TODO(), item)
		if err != nil {
			return id, err
		}

		return rs.InsertedID.(primitive.ObjectID), nil
	}

	item.UpdatedAt = time.Now()
	item.Status = repository.QueueItemStatusWait
	item.RequeueTimes = item.RequeueTimes + 1
	if item.NextExecuteAt.IsZero() {
		item.NextExecuteAt = item.UpdatedAt
	}

	if _, err := q.col.ReplaceOne(context.TODO(), bson.M{"_id": item.ID}, item); err != nil {
		return item.ID, err
	}

	return item.ID, nil
}

func (q *QueueRepo) Dequeue() (item repository.QueueItem, err error) {
	rs := q.col.FindOneAndUpdate(
		context.TODO(),
		bson.M{"status": repository.QueueItemStatusWait, "next_execute_at": bson.M{"$lt": time.Now()}},
		bson.M{"$set": bson.M{
			"status":     repository.QueueItemStatusRunning,
			"updated_at": time.Now(),
		}},
		options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After),
	)

	err = rs.Decode(&item)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}

	return
}

func (q *QueueRepo) Paginate(filter bson.M, offset, limit int64) (items []repository.QueueItem, next int64, err error) {
	cur, err := q.col.Find(context.TODO(), filter, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var item repository.QueueItem
		if err = cur.Decode(&item); err != nil {
			return
		}

		items = append(items, item)
	}

	if int64(len(items)) == limit {
		next = offset + limit
	}

	return
}

func (q *QueueRepo) Delete(filter bson.M) error {
	_, err := q.col.DeleteMany(context.TODO(), filter)
	return err
}

func (q *QueueRepo) DeleteID(id primitive.ObjectID) error {
	return q.Delete(bson.M{"_id": id})
}

func (q *QueueRepo) Get(id primitive.ObjectID) (repository.QueueItem, error) {
	rs := q.col.FindOne(context.TODO(), bson.M{"_id": id})
	var item repository.QueueItem
	if err := rs.Decode(&item); err != nil {
		if err == mongo.ErrNoDocuments {
			return item, repository.ErrNotFound
		}
		return item, err
	}

	return item, nil
}

func (q *QueueRepo) Count(filter bson.M) (int64, error) {
	return q.col.CountDocuments(context.TODO(), filter)
}

func (q *QueueRepo) Update(filter bson.M, item repository.QueueItem) error {
	item.UpdatedAt = time.Now()
	_, err := q.col.ReplaceOne(context.TODO(), filter, item)
	return err
}

func (q *QueueRepo) UpdateID(id primitive.ObjectID, item repository.QueueItem) error {
	return q.Update(bson.M{"_id": id}, item)
}