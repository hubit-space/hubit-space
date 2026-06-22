package handler

import (
	"hubit-space/service/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OptionHandler struct {
	repo repository.OptionRepository
}

func NewOption(repo repository.OptionRepository) *OptionHandler {
	return &OptionHandler{
		repo: repo,
	}
}

func (h *OptionHandler) GetOptions(c *gin.Context) {
	result, err := h.repo.GetOptions(c)
	if err != nil {
		log.Println("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *OptionHandler) GetParamater(c *gin.Context) {

	parameter := c.DefaultQuery("parameter_code", "")

	result, err := h.repo.GetParamater(parameter)
	if err != nil {
		log.Println("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
