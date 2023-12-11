package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AudioFile struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Genre      string             `bson:"genre"`
	Album      string             `bson:"album"`
	Artist     string             `bson:"artist"`
	TrackName  string             `bson:"track_name"`
	Path       string             `bson:"path"`
	UploadDate time.Time          `bson:"upload_date"`
}

type Selection struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	SelectionName string               `bson:"selection_name"`
	TrackIDs      []primitive.ObjectID `bson:"track_ids"`
}

type SelectionDto struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	SelectionName string             `bson:"selection_name"`
	Tracks        []AudioFile        `bson:"tracks"`
}
