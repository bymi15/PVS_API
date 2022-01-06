package services

import (
	"context"
	"errors"
	"time"

	"github.com/bymi15/PVS_API/db/models"
	"github.com/bymi15/PVS_API/db/permissions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShowcaseRoomService struct {
	Collection *mongo.Collection
}

func NewShowcaseRoomService(db *mongo.Database) ShowcaseRoomService {
	return ShowcaseRoomService{
		Collection: db.Collection("showcaseRoom"),
	}

}

func (service ShowcaseRoomService) GetShowcaseRooms(showOnlyListed bool) ([]models.ShowcaseRoom, error) {
	rooms := []models.ShowcaseRoom{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	qry := bson.D{}
	if showOnlyListed {
		qry = bson.D{{"isListed", true}}
	}

	cursor, err := service.Collection.Find(ctx, qry)
	if err != nil {
		defer cursor.Close(ctx)
		return rooms, err
	}

	for cursor.Next(ctx) {
		room := models.NewShowcaseRoom()
		err := cursor.Decode(&room)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (service ShowcaseRoomService) GetShowcaseRoomsByUser(userId string) ([]models.ShowcaseRoom, error) {
	rooms := []models.ShowcaseRoom{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	qry := []bson.M{
		{
			"$match": bson.M{
				"userId": userId,
			},
		},
	}

	cursor, err := service.Collection.Find(ctx, qry)
	if err != nil {
		defer cursor.Close(ctx)
		return nil, err
	}

	for cursor.Next(ctx) {
		room := models.NewShowcaseRoom()
		err := cursor.Decode(&room)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (service ShowcaseRoomService) GetShowcaseRoomById(id, authUserId string) (*models.ShowcaseRoom, error) {
	room := models.NewShowcaseRoom()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = service.Collection.FindOne(ctx, bson.D{{"_id", objectId}}).Decode(&room)
	if err != nil {
		return nil, err
	}
	if !permissions.IsRoomPublicOrOwner(authUserId, room) {
		return nil, errors.New("forbidden access")
	}
	return &room, nil

}

func (service ShowcaseRoomService) CreateShowcaseRoom(showcaseRoom models.ShowcaseRoom) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := service.Collection.InsertOne(ctx, showcaseRoom)
	if err != nil {
		return err
	}
	return nil
}

func (service ShowcaseRoomService) UpdateShowcaseRoom(id, authUserId string, showcaseRoom models.ShowcaseRoom) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data bson.M
	bytes, err := bson.Marshal(showcaseRoom)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	room := models.NewShowcaseRoom()
	err = service.Collection.FindOne(ctx, bson.D{{"_id", objectId}}).Decode(&room)
	if err != nil {
		return err
	}
	if !permissions.IsRoomOwner(authUserId, room) {
		return errors.New("forbidden access")
	}

	_, err = service.Collection.UpdateOne(
		ctx,
		bson.D{{"_id", objectId}},
		bson.D{{"$set", data}},
	)
	return err
}

func (service ShowcaseRoomService) DeleteShowcaseRoom(id, authUserId string) error {
	room := models.NewShowcaseRoom()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = service.Collection.FindOne(ctx, bson.D{{"_id", objectId}}).Decode(&room)
	if err != nil {
		return err
	}
	if !permissions.IsRoomOwner(authUserId, room) {
		return errors.New("forbidden access")
	}
	_, err = service.Collection.DeleteOne(ctx, bson.D{{"_id", objectId}})
	if err != nil {
		return err
	}
	return nil
}
