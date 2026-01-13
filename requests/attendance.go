package requests

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

//
// =======================
// Device Punch Request
// =======================
// This request comes from gym device or device gateway
//

type PunchRequest struct {
	AttendanceCode string    `json:"attendance_code" validate:"required"`
	
	DeviceID       string    `json:"device_id" validate:"required"`
	Seq            string    `json:"seq" validate:"required"`
	PunchTime      time.Time `json:"punch_time" validate:"required"`
	RawLine        string    `json:"raw_line,omitempty"`
	VerifyCode     string    `json:"verify_code,omitempty"`
}

//
// =======================
// Constructor
// =======================
//

func NewPunchRequest() *PunchRequest {
	return &PunchRequest{}
}

//
// =======================
// Validation
// =======================
//

func (r *PunchRequest) Validate(c *gin.Context) error {
	// Bind JSON
	if err := c.ShouldBindJSON(r); err != nil {
		return err
	}

	// Basic required validation
	if r.AttendanceCode == "" {
		return fmt.Errorf("attendance_code is required")
	}
	if r.DeviceID == "" {
		return fmt.Errorf("device_id is required")
	}
	if r.Seq == "" {
		return fmt.Errorf("seq (device serial) is required")
	}

	// Punch time should not be too far in the future
	if r.PunchTime.After(time.Now().Add(24 * time.Hour)) {
		return fmt.Errorf("punch_time cannot be more than 24 hours in the future")
	}

	return nil
}
