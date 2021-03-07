package service

import (
	"context"

	"github.com/arabian9ts/geekfes/domain/model"
)

func (s *service) GetSeries(ctx context.Context, id string) (model.Series, error) {
	sr, err := s.series.GetSeries(ctx, id)
	if err != nil {
		return model.Series{}, err
	}
	return sr, nil
}

func (s *service) ListSeries(ctx context.Context, limit int, since string) ([]model.Series, error) {
	srs, err := s.series.ListSeries(ctx, limit, since)
	if err != nil {
		return nil, err
	}
	return srs, nil
}
