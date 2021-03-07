package model

type Series struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	CreatedAt int64  `bson:"createdAt"`
	UpdatedAt int64  `bson:"updatedAt"`
}
