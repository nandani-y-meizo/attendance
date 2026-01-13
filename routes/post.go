package routes

import (
	"context"
	"net/http"
	"time"

	"attendance-service/requests"
	"attendance-service/services"

	"github.com/gin-gonic/gin"
)

// POST /api/v1/device/punch
// Called by biometric / QR / kiosk device
func DevicePunch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse & validate punch request
	req := requests.NewPunchRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	remoteIP := c.ClientIP()

	service := services.NewAttendanceService()
	record, err := service.HandlePunch(ctx, req, remoteIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Punch recorded successfully",
		"data":    record,
	})
}
