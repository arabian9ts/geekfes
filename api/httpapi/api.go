package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arabian9ts/geekfes/domain/service"
)

func New(svc service.Service) http.Handler {
	r := gin.New()

	r.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	ctrl := &controller{svc: svc}
	//r.Use(gzip.Gzip(gzip.DefaultCompression))

	dataset := r.Group("/dataset")
	dataset.GET("/series/:seriesId", ctrl.GetDataset)

	series := r.Group("/series")
	series.GET("", ctrl.ListSeries)
	series.GET("/:seriesId", ctrl.GetSeries)
	series.GET("/:seriesId/seasons", ctrl.ListSeasons)

	seasons := r.Group("/seasons")
	seasons.GET("/:seasonId/episodes", ctrl.ListEpisodes)

	episodes := r.Group("/episodes")
	episodes.GET("/:episodeId/next", ctrl.GetNextEpisode)
	episodes.GET("/:episodeId/prev", ctrl.GetPrevEpisode)

	return r
}
