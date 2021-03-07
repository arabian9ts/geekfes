package repository

import (
	"context"

	"github.com/arabian9ts/geekfes/domain/model"
)

type SeriesRepository interface {
	GetSeries(ctx context.Context, id string) (model.Series, error)
	ListSeries(ctx context.Context, limit int, since string) ([]model.Series, error)
}

type SeasonRepository interface {
	GetSeason(ctx context.Context, seasonID string) (model.Season, error)
	ListSeasonsBySeries(ctx context.Context, seriesID string) ([]model.Season, error)
}

type EpisodeRepository interface {
	GetEpisode(ctx context.Context, id string) (model.Episode, error)
	GetEpisodeWithSeq(ctx context.Context, seriesID, seasonID string, seq int32) (model.Episode, error)
	ListEpisodesBySeason(ctx context.Context, seasonID string) ([]model.Episode, error)
}
