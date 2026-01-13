package routes

import (
	"shared/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(api *gin.RouterGroup) {
	// ================= DEVICE PUNCH =================
	// No auth (device based or internal)
	device := api.Group("/device")
	{
		device.POST("/punch", DevicePunch)
	}

	// ================= ADMIN / DASHBOARD =================
	// Auth protected
	attendance := api.Group("/attendance")
	attendance.Use(middleware.AuthcMiddleware())
	{
		// Get all attendance for a gym (by seq)
		attendance.GET("/gym/:seq", GetGymAttendances)

		// Get attendance of a member
		attendance.GET("/member/:attendance_code", GetMemberAttendances)

		// Check if member is present today
		attendance.GET("/present/:attendance_code", IsMemberPresent)

		// Get attendance by EntityID
		attendance.GET("/entity/:entity_id", GetAttendancesByEntityID)
	}
}
