package services

import (
	"context"
	"time"

	"github.com/bymi15/PVS_API/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShowcaseRoomService struct {
	Collection *mongo.Collection
}

func NewShowcaseRoomService(db *mongo.Database, collectionName string) ShowcaseRoomService {
	return ShowcaseRoomService{
		Collection: db.Collection(collectionName),
	}

}

func (service ShowcaseRoomService) GetShowcaseRooms() ([]models.ShowcaseRoom, error) {
	ShowcaseRooms := []models.ShowcaseRoom{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := service.Collection.Find(ctx, bson.D{})
	if err != nil {
		defer cursor.Close(ctx)
		return ShowcaseRooms, err
	}

	for cursor.Next(ctx) {
		ShowcaseRoom := models.NewShowcaseRoom()
		err := cursor.Decode(&ShowcaseRoom)
		if err != nil {
			return ShowcaseRooms, err
		}
		ShowcaseRooms = append(ShowcaseRooms, ShowcaseRoom)
	}

	return ShowcaseRooms, nil
}

func (service ShowcaseRoomService) GetShowcaseRoomById(id string) (models.ShowcaseRoom, error) {
	ShowcaseRoom := models.NewShowcaseRoom()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ShowcaseRoom, err
	}

	err = service.Collection.FindOne(ctx, bson.D{{"_id", objectId}}).Decode(&ShowcaseRoom)
	if err != nil {
		return ShowcaseRoom, err
	}
	return ShowcaseRoom, nil

}

func (service ShowcaseRoomService) CreateShowcaseRoom(ShowcaseRoom models.ShowcaseRoom) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := service.Collection.InsertOne(ctx, ShowcaseRoom)
	if err != nil {
		return err
	}
	return nil
}

func (service ShowcaseRoomService) UpdateShowcaseRoom(id string, ShowcaseRoom models.ShowcaseRoom) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data bson.M
	bytes, err := bson.Marshal(ShowcaseRoom)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}
	_, err = service.Collection.UpdateOne(
		ctx,
		bson.D{{"_id", objectId}},
		bson.D{{"$set", data}},
	)
	return err
}

func (service ShowcaseRoomService) DeleteShowcaseRoom(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = service.Collection.DeleteOne(ctx, bson.D{{"_id", objectId}})
	if err != nil {
		return err
	}
	return nil
}
