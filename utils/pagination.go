package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

func GetPaginationParams(c *gin.Context) (*dto.PaginationDTO, error) {
	pagination := &dto.PaginationDTO{}

	// Bind query parameters ke struct
	if err := c.ShouldBindQuery(pagination); err != nil {
		return nil, err
	}

	// Set default values jika diperlukan
	if pagination.Limit == 0 {
		pagination.Limit = 10 // default limit
	}
	if pagination.Page == 0 {
		pagination.Page = 1 // default page
	}

	return pagination, nil
}
