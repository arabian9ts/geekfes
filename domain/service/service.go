package service

import (
	"context"

	"github.com/arabian9ts/geekfes/domain/model"
	"github.com/arabian9ts/geekfes/domain/repository"
)

type Service interface {
	DatasetService
	SeriesService
	SeasonService
	EpisodeService
}

type DatasetService interface {
	GetDataset(ctx context.Context, seriesID string) (model.Dataset, error)
}

type SeriesService interface {
	GetSeries(ctx context.Context, id string) (model.Series, error)
	ListSeries(ctx context.Context, limit int, since string) ([]model.Series, error)
}

type SeasonService interface {
	ListSeasons(ctx context.Context, seriesID string) ([]model.Season, error)
}

type EpisodeService interface {
	GetNextEpisode(ctx context.Context, epID string) (model.Episode, error)
	GetPrevEpisode(ctx context.Context, epID string) (model.Episode, error)
	ListEpisodes(ctx context.Context, seasonID string) ([]model.Episode, error)
}

var _ Service = &service{}

type service struct {
	series  repository.SeriesRepository
	season  repository.SeasonRepository
	episode repository.EpisodeRepository
}

func New(
	series repository.SeriesRepository,
	season repository.SeasonRepository,
	episode repository.EpisodeRepository,
) *service {
	return &service{
		series:  series,
		season:  season,
		episode: episode,
	}
}
