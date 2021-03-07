package model

type Episode struct {
	ID        string `bson:"_id"`
	SeriesID  string `bson:"seriesId"`
	SeasonID  string `bson:"seasonId"`
	Name      string `bson:"name"`
	Sequence  int32  `bson:"sequence"`
	CreatedAt int64  `bson:"createdAt"`
	UpdatedAt int64  `bson:"updatedAt"`
}
