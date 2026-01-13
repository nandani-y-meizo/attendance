package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"attendance-service/services"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/attendance/gym/:seq
func GetGymAttendances(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	seq := c.Param("seq")
	if seq == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "seq is required",
		})
		return
	}

	service := services.NewAttendanceService()
	data, err := service.GetAttendancesByGym(ctx, seq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
   fmt.Println("Retrieved attendances for gym seq", seq, ":", data)
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// GET /api/v1/attendance/member/:attendance_code
func GetMemberAttendances(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	code := c.Param("attendance_code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "attendance_code is required",
		})
		return
	}

	service := services.NewAttendanceService()
	data, err := service.GetAttendancesByMember(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// GET /api/v1/attendance/present/:attendance_code?seq=GYM001
func IsMemberPresent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	code := c.Param("attendance_code")
	seq := c.Query("seq")

	if code == "" || seq == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "attendance_code and seq are required",
		})
		return
	}

	service := services.NewAttendanceService()
	present, err := service.IsMemberPresentToday(ctx, code, seq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attendance_code": code,
		"seq":             seq,
		"present":         present,
	})
}

// GET /api/v1/attendance/entity/:entity_id
func GetAttendancesByEntityID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id := c.Param("entity_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "entity_id is required",
		})
		return
	}

	service := services.NewAttendanceService()
	data, err := service.GetAttendancesByEntityID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
