package service

import (
	"context"

	"github.com/arabian9ts/geekfes/domain/model"
)

func (s *service) GetNextEpisode(ctx context.Context, epID string) (model.Episode, error) {
	ep, err := s.episode.GetEpisode(ctx, epID)
	if err != nil {
		return model.Episode{}, err
	}
	nextSeq := ep.Sequence + 1
	next, err := s.episode.GetEpisodeWithSeq(ctx, ep.SeriesID, ep.SeasonID, nextSeq)
	if err != nil {
		return model.Episode{}, err
	}
	return next, nil
}

func (s *service) GetPrevEpisode(ctx context.Context, epID string) (model.Episode, error) {
	ep, err := s.episode.GetEpisode(ctx, epID)
	if err != nil {
		return model.Episode{}, err
	}
	prevSeq := ep.Sequence - 1
	prev, err := s.episode.GetEpisodeWithSeq(ctx, ep.SeriesID, ep.SeasonID, prevSeq)
	if err != nil {
		return model.Episode{}, err
	}
	return prev, nil
}

func (s *service) ListEpisodes(ctx context.Context, seasonID string) ([]model.Episode, error) {
	se, err := s.season.GetSeason(ctx, seasonID)
	if err != nil {
		return nil, err
	}
	eps, err := s.episode.ListEpisodesBySeason(ctx, se.ID)
	if err != nil {
		return nil, err
	}
	return eps, nil
}
