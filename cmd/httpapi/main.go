package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/xerrors"

	"github.com/arabian9ts/geekfes/api/httpapi"
	"github.com/arabian9ts/geekfes/domain/service"
	"github.com/arabian9ts/geekfes/infra/mongodb"
)

type config struct {
	MongoDB struct {
		URI      string        `envconfig:"URI" default:"mongodb://localhost:27017/geekfes"`
		Timeout  time.Duration `envconfig:"TIMEOUT" default:"15s"`
		Database string        `envconfig:"DATABASE" default:"geekfes"`
	} `envconfig:"MONGODB"`

	Port int32 `envconfig:"PORT" default:"8080"`
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	conf := config{}
	envconfig.MustProcess("FES", &conf)
	h, err := register(conf)
	if err != nil {
		return err
	}

	server := http.Server{
		Handler: h,
		Addr:    fmt.Sprintf(":%d", conf.Port),
	}
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}

	return nil
}

func register(conf config) (http.Handler, error) {
	mongoCli, err := newMongoDB(conf)
	if err != nil {
		return nil, err
	}
	srRepo := mongodb.NewSeries(mongoCli, conf.MongoDB.Database)
	seRepo := mongodb.NewSeason(mongoCli, conf.MongoDB.Database)
	epRepo := mongodb.NewEpisode(mongoCli, conf.MongoDB.Database)
	svc := service.New(srRepo, seRepo, epRepo)
	h := httpapi.New(svc)
	return h, nil
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
