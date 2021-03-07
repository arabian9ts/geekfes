package service

import (
	"context"

	"github.com/arabian9ts/geekfes/domain/model"
)

func (s *service) GetDataset(ctx context.Context, seriesID string) (model.Dataset, error) {
	sr, err := s.series.GetSeries(ctx, seriesID)
	if err != nil {
		return model.Dataset{}, err
	}
	ses, err := s.season.ListSeasonsBySeries(ctx, sr.ID)
	if err != nil {
		return model.Dataset{}, err
	}
	ds := model.Dataset{
		Series:   sr,
		Seasons:  ses,
		Episodes: make([]model.Episode, 0, 100),
	}
	for i := range ses {
		eps, err := s.ListEpisodes(ctx, ses[i].ID)
		if err != nil {
			return model.Dataset{}, err
		}
		ds.Episodes = append(ds.Episodes, eps...)
	}
	return ds, nil
}
