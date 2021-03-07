package service

import (
	"context"

	"github.com/arabian9ts/geekfes/domain/model"
)

func (s *service) ListSeasons(ctx context.Context, seriesID string) ([]model.Season, error) {
	sr, err := s.series.GetSeries(ctx, seriesID)
	if err != nil {
		return nil, err
	}
	ses, err := s.season.ListSeasonsBySeries(ctx, sr.ID)
	if err != nil {
		return nil, err
	}
	return ses, nil
}
