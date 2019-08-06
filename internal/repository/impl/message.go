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

type MessageRepo struct {
	col *mongo.Collection
}

func NewMessageRepo(db *mongo.Database) repository.MessageRepo {
	return &MessageRepo{col: db.Collection("message")}
}

func (m MessageRepo) Add(msg repository.Message) (id primitive.ObjectID, err error) {
	msg.CreatedAt = time.Now()
	if msg.Status == "" {
		msg.Status = repository.MessageStatusPending
	}

	rs, err := m.col.InsertOne(context.TODO(), msg)
	if err != nil {
		return id, err
	}

	return rs.InsertedID.(primitive.ObjectID), err
}

func (m MessageRepo) Get(id primitive.ObjectID) (msg repository.Message, err error) {
	err = m.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&msg)
	if err == mongo.ErrNoDocuments {
		return msg, repository.ErrNotFound
	}

	return msg, err
}

func (m MessageRepo) Find(filter bson.M) (messages []repository.Message, err error) {
	cur, err := m.col.Find(context.TODO(), filter)
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var msg repository.Message
		if err = cur.Decode(&msg); err != nil {
			return
		}

		messages = append(messages, msg)
	}

	return
}

func (m MessageRepo) Paginate(filter bson.M, offset, limit int64) (messages []repository.Message, next int64, err error) {
	cur, err := m.col.Find(context.TODO(), filter, options.Find().SetLimit(limit).SetSkip(offset))
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var msg repository.Message
		if err = cur.Decode(&msg); err != nil {
			return
		}

		messages = append(messages, msg)
	}

	if int64(len(messages)) == limit {
		next = offset + limit
	}

	return messages, next, err
}

func (m MessageRepo) Delete(filter bson.M) error {
	_, err := m.col.DeleteMany(context.TODO(), filter)
	return err
}

func (m MessageRepo) DeleteID(id primitive.ObjectID) error {
	return m.Delete(bson.M{"_id": id})
}

func (m MessageRepo) Traverse(filter bson.M, cb func(msg repository.Message) error) error {
	cur, err := m.col.Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	for cur.Next(context.TODO()) {
		var msg repository.Message
		if err = cur.Decode(&msg); err != nil {
			return err
		}

		if err = cb(msg); err != nil {
			return err
		}
	}

	return nil
}

func (m MessageRepo) UpdateID(id primitive.ObjectID, update repository.Message) error {
	_, err := m.col.ReplaceOne(context.TODO(), bson.M{"_id": id}, update)
	return err
}

func (m MessageRepo) Count(filter bson.M) (int64, error) {
	return m.col.CountDocuments(context.TODO(), filter)
}