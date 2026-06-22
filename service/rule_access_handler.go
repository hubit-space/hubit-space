package handler

import (
	"hubit-space/service/repository"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocnumHandler struct {
	repo repository.DocnumRepository
}

func NewDocnumHandler(repo repository.DocnumRepository) *DocnumHandler {
	return &DocnumHandler{
		repo: repo,
	}
}

func (h *DocnumHandler) GetLastDocNumber(ctx *gin.Context) {
	log.Println("GetLastDocNumber called")
	var requestBody map[string]any

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	docname, ok := requestBody["docname"].(string)
	if !ok || docname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Docname is required and must be a string"})
		return
	}
	log.Println(docname)

	year := 0
	if y, ok := requestBody["year"].(float64); ok {
		year = int(y)
	}

	month := 0
	if m, ok := requestBody["month"].(float64); ok {
		month = int(m)
	}

	day := 0
	if d, ok := requestBody["day"].(float64); ok {
		day = int(d)
	}

	result := h.repo.GetLastDocNumber(docname, year, month, day)
	if result["message"] != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result["message"]})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
