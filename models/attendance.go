package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//
// =====================
// Attendance Model
// =====================
// Represents a single punch record coming from a gym device
//

type Attendance struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EntityID       string             `bson:"entity_id" json:"entity_id"`
	// MemberEntityID string             `bson:"member_entity_id" json:"member_entity_id"`
	// FullName      string             `bson:"full_name" json:"full_name"`
	AttendanceCode string             `bson:"attendance_code" json:"attendance_code"`
	DeviceID       string             `bson:"device_id" json:"device_id"`
	PunchTime      time.Time          `bson:"punch_time" json:"punch_time"`
	Seq            string             `bson:"seq" json:"seq"`
	RawLine        string             `bson:"raw_line" json:"raw_line"`
	ReceivedAt     time.Time          `bson:"received_at" json:"received_at"`
	RemoteIP       string             `bson:"remote_ip" json:"remote_ip"`
	Status         PunchStatus        `bson:"status" json:"status"`
	VerifyCode     string             `bson:"verify_code" json:"verify_code"`
}

//
// =====================
// Punch Status Enum
// =====================
//

type PunchStatus string

const (
	PunchIn  PunchStatus = "IN"
	PunchOut PunchStatus = "OUT"
)

//
// =====================
// Daily Attendance Status
// =====================
// Used when checking if member is present or absent
//

type AttendanceDayStatus string

const (
	DayPresent AttendanceDayStatus = "PRESENT"
	DayAbsent  AttendanceDayStatus = "ABSENT"
)

//
// =====================
// Helper Methods
// =====================
//

func (a *Attendance) BindPunch(
	entityID string,
	memberEntityID string,
	attendanceCode string,
	// full_name string,
	deviceID string,
	seq string,
	rawLine string,
	remoteIP string,
	verifyCode string,
	status PunchStatus,
	punchTime time.Time,
) {
	a.EntityID = entityID
	// a.MemberEntityID = memberEntityID
	a.AttendanceCode = attendanceCode
	// a.FullName = full_name
	a.DeviceID = deviceID
	a.Seq = seq
	a.RawLine = rawLine
	a.RemoteIP = remoteIP
	a.VerifyCode = verifyCode
	a.Status = status
	a.PunchTime = punchTime
	a.ReceivedAt = time.Now()
}
