package models

import "time"

type ProjectStand struct {
	PosX         *int    `json:"posX,omitempty" bson:"posX,omitempty"`
	PosY         *int    `json:"posY,omitempty" bson:"posY,omitempty"`
	Title        *string `json:"title,omitempty" bson:"title,omitempty"`
	Description  *string `json:"description,omitempty" bson:"description,omitempty"`
	ThumbnailUrl *string `json:"thumbnailUrl,omitempty" bson:"thumbnailUrl,omitempty"`
	ProjectUrl   *string `json:"projectUrl,omitempty" bson:"projectUrl,omitempty"`
	ObjectType   *string `json:"objectType,omitempty" bson:"objectType,omitempty"`
}

type Scene struct {
	Type          *string `json:"type,omitempty" bson:"type,omitempty"`
	BackgroundUrl *string `json:"backgroundUrl,omitempty" bson:"backgroundUrl,omitempty"`
	Size          *string `json:"size,omitempty" bson:"size,omitempty"`
}

type User struct {
	UserId      *string `json:"userId,omitempty" bson:"userId,omitempty"`
	FullName    *string `json:"fullName,omitempty" bson:"fullName,omitempty"`
	Email       *string `json:"email,omitempty" bson:"email,omitempty"`
	Role        *string `json:"role,omitempty" bson:"role,omitempty"`
	DateCreated *string `json:"dateCreated,omitempty" bson:"dateCreated,omitempty"`
}

type ShowcaseRoom struct {
	Id            string         `json:"id,omitempty" bson:"_id,omitempty"`
	RoomName      *string        `json:"roomName,omitempty" bson:"roomName,omitempty"`
	Scene         *Scene         `json:"scene,omitempty" bson:"scene,omitempty"`
	DateCreated   *string        `json:"dateCreated,omitempty" bson:"dateCreated,omitempty"`
	CreatedBy     *User          `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	IsListed      *bool          `json:"isListed,omitempty" bson:"isListed,omitempty"`
	ProjectStands []ProjectStand `json:"projectStands,omitempty" bson:"projectStands,omitempty"`
}

func NewShowcaseRoom() ShowcaseRoom {
	instance := ShowcaseRoom{}
	roomName := "Showcase Room"
	dateCreated := time.Now().Format("2006-01-02")
	isListed := true
	sceneType := "Custom"
	sceneSize := "md"
	instance.RoomName = &roomName
	instance.ProjectStands = []ProjectStand{}
	instance.DateCreated = &dateCreated
	instance.IsListed = &isListed
	scene := Scene{}
	scene.Type = &sceneType
	scene.Size = &sceneSize
	instance.Scene = &scene
	return instance
}
