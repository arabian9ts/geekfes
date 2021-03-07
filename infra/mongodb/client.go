package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/arabian9ts/geekfes/code"
	"github.com/arabian9ts/geekfes/domain/model"
	"github.com/arabian9ts/geekfes/domain/repository"
)

var (
	_ repository.SeriesRepository  = &Client{}
	_ repository.SeasonRepository  = &Client{}
	_ repository.EpisodeRepository = &Client{}
)

type Client struct {
	*mongo.Client
	db string
}

func (c *Client) GetSeries(ctx context.Context, id string) (model.Series, error) {
	q := bson.M{"_id": id}
	sr := model.Series{}
	if err := c.collection(seriesCol).FindOne(ctx, q).Decode(&sr); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Series{}, code.Errorf(code.NotFound, "series not found: id=%s", id)
		}
		return model.Series{}, code.Errorf(code.Internal, "find series error: %w", err)
	}
	return sr, nil
}

func (c *Client) ListSeries(ctx context.Context, limit int, since string) ([]model.Series, error) {
	q := bson.M{"_id": bson.M{"$gt": since}}
	cur, err := c.collection(seriesCol).Find(ctx, q, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, code.Errorf(code.Internal, "list series error: %w", err)
	}
	srs := make([]model.Series, 0, limit)
	if err := cur.All(ctx, &srs); err != nil {
		return nil, code.Errorf(code.Internal, "list series error: %w", err)
	}
	return srs, nil
}

func (c *Client) GetSeason(ctx context.Context, seasonID string) (model.Season, error) {
	q := bson.M{"_id": seasonID}
	se := model.Season{}
	if err := c.collection(seasonCol).FindOne(ctx, q).Decode(&se); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Season{}, code.Errorf(code.NotFound, "season not found: id=%s", seasonID)
		}
		return model.Season{}, code.Errorf(code.Internal, "find season error: %w", err)
	}
	return se, nil
}

func (c *Client) ListSeasonsBySeries(ctx context.Context, seriesID string) ([]model.Season, error) {
	q := bson.M{"seriesId": seriesID}
	cur, err := c.collection(seasonCol).Find(ctx, q)
	if err != nil {
		return nil, code.Errorf(code.Internal, "list season error: %w", err)
	}
	ses := make([]model.Season, 0, 10)
	if err := cur.All(ctx, &ses); err != nil {
		return nil, code.Errorf(code.Internal, "list season error: %w", err)
	}
	return ses, nil
}

func (c *Client) GetEpisode(ctx context.Context, id string) (model.Episode, error) {
	q := bson.M{"_id": id}
	ep := model.Episode{}
	if err := c.collection(episodeCol).FindOne(ctx, q).Decode(&ep); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Episode{}, code.Errorf(code.NotFound, "episode not found: id=%s", id)
		}
		return model.Episode{}, code.Errorf(code.Internal, "find episode error: %w", err)
	}
	return ep, nil
}

func (c *Client) GetEpisodeWithSeq(ctx context.Context, seriesID, seasonID string, seq int32) (model.Episode, error) {
	q := bson.M{"seriesId": seriesID, "seasonId": seasonID, "sequence": seq}
	ep := model.Episode{}
	if err := c.collection(episodeCol).FindOne(ctx, q).Decode(&ep); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Episode{}, code.Errorf(code.NotFound, "episode not found: seriesId=%s, seasonId=%s, sequence=%d", seriesID, seasonID, seq)
		}
		return model.Episode{}, code.Errorf(code.Internal, "find episode error: %w", err)
	}
	return ep, nil
}

func (c *Client) ListEpisodesBySeason(ctx context.Context, seasonID string) ([]model.Episode, error) {
	q := bson.M{"seasonId": seasonID}
	cur, err := c.collection(episodeCol).Find(ctx, q)
	if err != nil {
		return nil, code.Errorf(code.Internal, "list episode error: %w", err)
	}
	eps := make([]model.Episode, 0, 100)
	if err := cur.All(ctx, &eps); err != nil {
		return nil, code.Errorf(code.Internal, "list episode error: %w", err)
	}
	return eps, nil
}

func NewSeries(cli *mongo.Client, db string) *Client {
	return &Client{cli, db}
}

func NewSeason(cli *mongo.Client, db string) *Client {
	return &Client{cli, db}
}

func NewEpisode(cli *mongo.Client, db string) *Client {
	return &Client{cli, db}
}

func (c *Client) database(name string) *mongo.Database {
	return c.Database(name)
}

func (c *Client) collection(name string) *mongo.Collection {
	return c.database(c.db).Collection(name)
}
