package ginserver

import (
	"fmt"
	"net/http"
	"zd/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	zendeskService ports.ZendeskService
}

func New(service ports.ZendeskService) *Gin {
	return &Gin{
		zendeskService: service,
	}
}

func (s Gin) GetUserEvent(ctx *gin.Context) {
	userEvent, err := s.zendeskService.GetUserEvent()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Error getting new user event: %s", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, *userEvent)
}
