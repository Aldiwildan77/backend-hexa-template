package http

import "github.com/Aldiwildan77/backend-hexa-template/pkg/pagination"

type Response struct {
	Data              interface{} `json:"data,omitempty"`
	Message           string      `json:"message"`
	*MetadataResponse `json:"metadata,omitempty"`
}

type MetadataResponse struct {
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
}
