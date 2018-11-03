package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/robot/pkg/api/operations"
	"github.com/supergiant/robot/pkg/models"
	"github.com/supergiant/robot/pkg/storage"
)

type recommendationPluginsHandler struct {
	storage storage.Interface
	log     logrus.FieldLogger
}

func NewRecommendationPluginsHandler(storage storage.Interface, logger logrus.FieldLogger) operations.GetRecommendationPluginsHandler {
	return &recommendationPluginsHandler{
		storage: storage,
		log:     logger,
	}
}

func (h *recommendationPluginsHandler) Handle(params operations.GetRecommendationPluginsParams) middleware.Responder {
	h.log.Infof("got request at: %v, request: %+v", time.Now(), params)
	pluginRaw, err := h.storage.GetAll(context.Background(), "/robot/plugins/")

	if err != nil {
		r := operations.NewGetRecommendationPluginsDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	result := &operations.GetRecommendationPluginsOKBody{
		InstalledRecommendationPlugins: []*models.RecommendationPlugin{},
	}

	for _, rawPlugin := range pluginRaw {
		p := &models.RecommendationPlugin{}
		err := p.UnmarshalBinary(rawPlugin)
		if err != nil {
			r := operations.NewGetRecommendationPluginsDefault(http.StatusInternalServerError)
			msg := err.Error()
			r.Payload = &models.Error{
				Code:    http.StatusInternalServerError,
				Message: &msg,
			}
			return r
		}
		result.InstalledRecommendationPlugins = append(result.InstalledRecommendationPlugins, p)
	}
	h.log.Infof("request processing finished at: %v, request: %+v", time.Now(), params)
	return operations.NewGetRecommendationPluginsOK().WithPayload(result)
}
