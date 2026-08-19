package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (oh *OrderHandler) GetReservation(c *gin.Context) {
	reservationID, ok := reservationIDFromPath(c)
	if !ok {
		return
	}

	reservation, err := oh.orderService.GetReservation(c.Request.Context(), reservationID)
	if err != nil {
		respondLifecycleError(c, err)
		return
	}

	c.JSON(http.StatusOK, reservation)
}