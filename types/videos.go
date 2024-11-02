package types

import (
	"fampay-assignment/config"
	"fampay-assignment/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type GetLatestVideosRequest struct {
	SortOrder      string `json:"sort_order"`
	PaginationSize int    `json:"pagination_size"`
	PaginationPage int    `json:"pagination_page"`
	PublishedAfter string `json:"published_after"`
}

func (req GetLatestVideosRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.SortOrder, validation.Required, validation.In("asc", "desc")),
		validation.Field(&req.PaginationSize, validation.Required, validation.Min(1), validation.Max(config.MAX_PAGINATION_SIZE)),
		validation.Field(&req.PaginationPage, validation.Required, validation.Min(1)),
		validation.Field(&req.PublishedAfter, validation.Date(config.DATE_FORMAT)),
	)
}

type GetLatestVideosResponse struct {
	Videos []models.Video `json:"videos"`
}
