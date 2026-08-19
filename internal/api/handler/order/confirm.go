package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (oh *OrderHandler) Confirm(c *gin.Context) {
	reservationID, ok := reservationIDFromPath(c)
	if !ok {
		return
	}

	reservation, err := oh.orderService.Confirm(c.Request.Context(), reservationID)
	if err != nil {
		respondLifecycleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservation_id": reservation.ID,
		"status":         reservation.Status,
	})
}
