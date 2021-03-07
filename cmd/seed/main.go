package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/teris-io/shortid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/xerrors"

	"github.com/arabian9ts/geekfes/domain/model"
)

// collection
const (
	seriesCol  = "series"
	seasonCol  = "seasons"
	episodeCol = "episodes"
)

// seed config
type config struct {
	MongoDB struct {
		URI      string        `envconfig:"URI" default:"mongodb://localhost:27017/geekfes"`
		Timeout  time.Duration `envconfig:"TIMEOUT" default:"15s"`
		Database string        `envconfig:"DATABASE" default:"geekfes"`
	} `envconfig:"MONGODB"`
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	conf := config{}
	envconfig.MustProcess("FES", &conf)
	mongoCli, err := newMongoDB(conf)
	if err != nil {
		return err
	}

	ctx := context.Background()
	defer mongoCli.Disconnect(ctx)

	db := mongoCli.Database(conf.MongoDB.Database)

	//
	// drop
	//
	db.Collection(seriesCol).Drop(ctx)
	db.Collection(seasonCol).Drop(ctx)
	db.Collection(episodeCol).Drop(ctx)

	//
	// fake series
	//
	srs, err := fakeSeries()
	if err != nil {
		return err
	}
	{
		ops := []mongo.WriteModel{}
		for i := range srs {
			op := mongo.NewInsertOneModel()
			op.SetDocument(srs[i])
			ops = append(ops, op)
		}
		srCol := db.Collection(seriesCol)
		if _, err := srCol.BulkWrite(ctx, ops); err != nil {
			return err
		}
	}

	//
	// fake seasons
	//
	ses := make([]model.Season, 0, 300)
	for i := range srs {
		ss, err := fakeSeasons(srs[i])
		if err != nil {
			return err
		}
		ses = append(ses, ss...)
	}
	{
		ops := []mongo.WriteModel{}
		for i := range ses {
			op := mongo.NewInsertOneModel()
			op.SetDocument(ses[i])
			ops = append(ops, op)
		}
		srCol := db.Collection(seasonCol)
		if _, err := srCol.BulkWrite(ctx, ops); err != nil {
			return err
		}
	}

	//
	// fake episodes
	//
	eps := make([]model.Episode, 0, 3000)
	for i := range ses {
		ep, err := fakeEpisodes(ses[i])
		if err != nil {
			return err
		}
		eps = append(eps, ep...)
	}
	{
		ops := []mongo.WriteModel{}
		for i := range eps {
			op := mongo.NewInsertOneModel()
			op.SetDocument(eps[i])
			ops = append(ops, op)
		}
		srCol := db.Collection(episodeCol)
		if _, err := srCol.BulkWrite(ctx, ops); err != nil {
			return err
		}
	}

	return nil
}

func newMongoDB(conf config) (*mongo.Client, error) {
	cli, err := mongo.NewClient(
		options.Client().ApplyURI(conf.MongoDB.URI),
		options.Client().SetConnectTimeout(conf.MongoDB.Timeout),
		options.Client().SetReadPreference(readpref.SecondaryPreferred()),
	)
	if err != nil {
		return nil, xerrors.Errorf("new mongodb client error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := cli.Connect(ctx); err != nil {
		return nil, fmt.Errorf("mongodb connection error: %w", err)
	}

	if err := cli.Ping(ctx, readpref.Primary()); err != nil {
		return nil, xerrors.Errorf("ping mongodb error: %w", err)
	}

	return cli, nil
}

func fakeSeries() ([]model.Series, error) {
	name := "シリーズ %03d"
	srs := make([]model.Series, 0, 100)
	for i := 1; i <= 100; i++ {
		id, err := shortid.Generate()
		if err != nil {
			return nil, err
		}
		srs = append(srs, model.Series{
			ID:        id,
			Name:      fmt.Sprintf(name, i),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		})
	}
	return srs, nil
}

func fakeSeasons(series model.Series) ([]model.Season, error) {
	name := "シリーズ %03d > シーズン %03d"
	ses := make([]model.Season, 0, 3)
	for i := 1; i <= 3; i++ {
		id, err := shortid.Generate()
		if err != nil {
			return nil, err
		}
		ses = append(ses, model.Season{
			ID:        id,
			SeriesID:  series.ID,
			Name:      fmt.Sprintf(name, series.ID, i),
			Sequence:  int32(i),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		})
	}
	return ses, nil
}

func fakeEpisodes(season model.Season) ([]model.Episode, error) {
	name := "シリーズ %03d > シーズン %03d > エピソード %03d"
	eps := make([]model.Episode, 0, 10)
	for i := 1; i <= 10; i++ {
		id, err := shortid.Generate()
		if err != nil {
			return nil, err
		}
		eps = append(eps, model.Episode{
			ID:        id,
			SeriesID:  season.SeriesID,
			SeasonID:  season.ID,
			Name:      fmt.Sprintf(name, season.SeriesID, season.ID, i),
			Sequence:  int32(i),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		})
	}
	return eps, nil
}
