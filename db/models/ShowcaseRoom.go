package models

import "time"

type Stand struct {
	Id           string  `json:"id,omitempty" bson:"_id,omitempty"`
	PosX         *int    `json:"posX,omitempty" bson:"posX,omitempty"`
	PosY         *int    `json:"posY,omitempty" bson:"posY,omitempty"`
	Title        *string `json:"title,omitempty" bson:"title,omitempty"`
	Description  *string `json:"description,omitempty" bson:"description,omitempty"`
	ThumbnailUrl *string `json:"thumbnailUrl,omitempty" bson:"thumbnailUrl,omitempty"`
	ProjectUrl   *string `json:"projectUrl,omitempty" bson:"projectUrl,omitempty"`
	ObjectType   *string `json:"objectType,omitempty" bson:"objectType,omitempty"`
}

type ShowcaseRoom struct {
	Id          string  `json:"id,omitempty" bson:"_id,omitempty"`
	RoomName    *string `json:"roomName,omitempty" bson:"roomName,omitempty"`
	Scene       *string `json:"scene,omitempty" bson:"scene,omitempty"`
	DateCreated *string `json:"dateCreated,omitempty" bson:"dateCreated,omitempty"`
	CreatedBy   *string `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	Stands      []Stand `json:"stands,omitempty" bson:"stands,omitempty"`
}

// Constructor
func NewShowcaseRoom() ShowcaseRoom {
	instance := ShowcaseRoom{}
	*instance.RoomName = "Showcase Room"
	*instance.Scene = "Default"
	instance.Stands = []Stand{}
	*instance.DateCreated = time.Now().Format("2006-01-02")
	return instance
}
