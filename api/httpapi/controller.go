package httpapi

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/arabian9ts/geekfes/domain/service"
)

type controller struct {
	svc service.Service
}

// /dataset/series/:seriesId
func (ctrl *controller) GetDataset(c *gin.Context) {
	seriesID := c.Param("seriesId")
	ds, err := ctrl.svc.GetDataset(c, seriesID)
	if err != nil {
		apiError(c, err)
		return
	}
	write(c, ds)
}

// /series/:seriesId
func (ctrl *controller) GetSeries(c *gin.Context) {
	id := c.Param("seriesId")
	sr, err := ctrl.svc.GetSeries(c, id)
	if err != nil {
		apiError(c, err)
		return
	}
	write(c, sr)
}

// /series?limit=100
func (ctrl *controller) ListSeries(c *gin.Context) {
	l := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(l)
	if err != nil {
		badRequest(c, errors.New("invalid query: limit"))
		return
	}
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}
	since := c.Query("since")
	srs, err := ctrl.svc.ListSeries(c, limit, since)
	if err != nil {
		apiError(c, err)
		return
	}
	write(c, srs)
}

// /series/:seriesId/seasons
func (ctrl *controller) ListSeasons(c *gin.Context) {
	srID := c.Param("seriesId")
	ses, err := ctrl.svc.ListSeasons(c, srID)
	if err != nil {
		apiError(c, err)
		return
	}
	write(c, ses)
}

// /seasons/:seasonsId/episodes
func (ctrl *controller) ListEpisodes(c *gin.Context) {
	seID := c.Param("seasonId")
	eps, err := ctrl.svc.ListEpisodes(c, seID)
	if err != nil {
		apiError(c, err)
		return
	}
	write(c, eps)
}

// /episodes/:episodeId/next
func (ctrl *controller) GetNextEpisode(c *gin.Context) {
	epID := c.Param("episodeId")
	next, err := ctrl.svc.GetNextEpisode(c, epID)
	if err != nil {
		apiError(c, err)
		return
	}
	write(c, next)
}

// /episodes/:episodeId/next
func (ctrl *controller) GetPrevEpisode(c *gin.Context) {
	epID := c.Param("episodeId")
	next, err := ctrl.svc.GetPrevEpisode(c, epID)
	if err != nil {
		apiError(c, err)
		return
	}
	write(c, next)
}
