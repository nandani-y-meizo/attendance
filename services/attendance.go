package services

import (
	"context"
	"fmt"
	"time"

	"attendance-service/models"
	"shared/pkgs/uuids"
	"attendance-service/requests"
	"attendance-service/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionAttendance = "attendance"
	DBName               = "gym_attendance"
)

type AttendanceService struct{}

func NewAttendanceService() *AttendanceService {
	return &AttendanceService{}
}

// HandlePunch handles automatic IN/OUT logic based on previous punches today
func (s *AttendanceService) HandlePunch(
	ctx context.Context,
	req *requests.PunchRequest,
	remoteIP string,
) (*models.Attendance, error) {

	db := storage.GetMongo()
	collection := db.Database(DBName).Collection(CollectionAttendance)

	// Normalize punch time
	punchTime := req.PunchTime.UTC()

	startOfDay := time.Date(
		punchTime.Year(),
		punchTime.Month(),
		punchTime.Day(),
		0, 0, 0, 0,
		time.UTC,
	)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Find last punch for this member today at this gym
	filter := bson.M{
		"attendance_code": req.AttendanceCode,
		"seq":             req.Seq,
		"punch_time": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	opts := options.FindOne().SetSort(bson.M{"punch_time": -1})

	var lastPunch models.Attendance
	err := collection.FindOne(ctx, filter, opts).Decode(&lastPunch)

	// Decide punch status: if last was IN, this is OUT
	status := models.PunchIn
	if err == nil && lastPunch.Status == models.PunchIn {
		status = models.PunchOut
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	// Find member to get MemberEntityID
	memberCollection := db.Database("fitpro").Collection("members")
	var member struct {
		EntityID string `bson:"entity_id"`
	}
	err = memberCollection.FindOne(ctx, bson.M{"attendance_code": req.AttendanceCode}).Decode(&member)
	memberEntityID := ""
	if err == nil {
		memberEntityID = member.EntityID
	}

	// Create attendance record
	record := &models.Attendance{}
	record.ID = primitive.NewObjectID()
	entityID, _ := uuids.NewUUID5(record.ID.Hex(), uuids.OidNamespace)

	record.BindPunch(
		entityID,
		memberEntityID,
		req.AttendanceCode,
		req.DeviceID,
		req.Seq,
		req.RawLine,
		remoteIP,
		req.VerifyCode,
		status,
		punchTime,
	)

	_, err = collection.InsertOne(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("failed to save attendance: %w", err)
	}
 fmt.Println("Inserted attendance record:", record)
	return record, nil
}

// IsMemberPresentToday checks if a member has any punch today
func (s *AttendanceService) IsMemberPresentToday(
	ctx context.Context,
	attendanceCode string,
	seq string,
) (bool, error) {

	db := storage.GetMongo()
	collection := db.Database(DBName).Collection(CollectionAttendance)

	now := time.Now().UTC()
	startOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.UTC,
	)
	endOfDay := startOfDay.Add(24 * time.Hour)

	filter := bson.M{
		"attendance_code": attendanceCode,
		"seq":             seq,
		"punch_time": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetAttendancesByGym retrieves all attendance for a gym
func (s *AttendanceService) GetAttendancesByGym(
	ctx context.Context,
	seq string,
) ([]models.Attendance, error) {

	db := storage.GetMongo()
	collection := db.Database(DBName).Collection(CollectionAttendance)

	filter := bson.M{"seq": seq}
	opts := options.Find().SetSort(bson.M{"punch_time": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []models.Attendance
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}

	return records, nil
}

// GetAttendancesByMember retrieves all attendance for a member
func (s *AttendanceService) GetAttendancesByMember(
	ctx context.Context,
	attendanceCode string,
) ([]models.Attendance, error) {

	db := storage.GetMongo()
	collection := db.Database(DBName).Collection(CollectionAttendance)

	filter := bson.M{"attendance_code": attendanceCode}
	opts := options.Find().SetSort(bson.M{"punch_time": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []models.Attendance
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}

	return records, nil
}

// GetAttendancesByEntityID retrieves all attendance for a member by their EntityID
func (s *AttendanceService) GetAttendancesByEntityID(
	ctx context.Context,
	entityID string,
) ([]models.Attendance, error) {

	db := storage.GetMongo()
	collection := db.Database(DBName).Collection(CollectionAttendance)

	filter := bson.M{"entity_id": entityID}
	opts := options.Find().SetSort(bson.M{"punch_time": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []models.Attendance
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}
    fmt.Println("Fetched attendance records:", records)
	return records, nil 	
}
