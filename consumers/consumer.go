package consumers

import (
	"encoding/json"
	"fmt"
	"log"

	"shared/constants"
	"shared/events"

	"github.com/IBM/sarama"
)

func GymPunchHandler(msg *sarama.ConsumerMessage) error {
	var evt events.EventMessage
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		log.Printf("❌ Failed to parse Attendance event: %v", err)
		return err
	}

	fmt.Println("📨 Attendance EVENT RECEIVED:", evt.EventType)

	// ctx := context.Background()

	if evt.EventType == constants.TopicGymPunchPush {
		log.Println("✅ Attendance: DevicePunch event received")

		// Debug: Log the entire payload
		payloadJSON, _ := json.Marshal(evt.Payload)
		log.Printf("📨 INCOMING PAYLOAD: %s", string(payloadJSON))

		// Extract deviceId from payload
		deviceId, ok := evt.Payload["deviceId"].(string)
		if !ok || deviceId == "" {
			log.Printf("❌ Missing deviceId in payload")
			return fmt.Errorf("missing deviceId")
		}

		// Call license RPC to get company_code from device serial number
		log.Printf("🔍 Fetching device info for deviceId: %s", deviceId)
		// deviceInfo, err := getDeviceInfo(ctx, deviceId)
		// if err != nil {
		// 	log.Printf("❌ Failed to get device info: %v", err)
		// 	return fmt.Errorf("failed to get device info: %w", err)
		// }

		// if !deviceInfo.IsActive {
		// 	log.Printf("⚠️ Device %s is not active, skipping", deviceId)
		// 	return nil
		// }

		// companyCode := deviceInfo.CompanyCode
		// log.Printf("✅ Device belongs to company: %s (device: %s, location: %s)", companyCode, deviceInfo.DeviceName, deviceInfo.Location)

		// // Add device metadata to payload for processing
		// evt.Payload["device_in_out"] = deviceInfo.InOut
		// evt.Payload["device_location"] = deviceInfo.Location
		// evt.Payload["device_name"] = deviceInfo.DeviceName
		// evt.Payload["device_code"] = deviceInfo.DeviceCode
		// // evt.Payload["product_code"] = deviceInfo

		// // Call service to process the punch records
		// // attendanceService := services.NewAttendanceService()
		// // err = attendanceService.ProcessDevicePunches(ctx, companyCode, evt.Payload)
		// if err != nil {
		// 	log.Printf("❌ Failed to process device punches: %v", err)
		// 	return err
		// }

		// log.Printf("✅ Successfully processed device punches for company %s", companyCode)
		return nil
	} else {
		log.Printf("⚠️ Attendance: Unhandled event type %s", evt.EventType)
		return nil
	}
}

// getDeviceInfo calls the license RPC to get device information by serial number
// func getDeviceInfo(ctx context.Context, serialNo string) (*license.GetDeviceBySerialResponse, error) {
// 	client := clients.GetLicenseRPCClient()
// 	if client == nil {
// 		return nil, fmt.Errorf("license RPC client not initialized")
// 	}

// 	req := &license.GetDeviceBySerialRequest{
// 		SerialNo: serialNo,
// 	}

// 	reply, err := client.GetDeviceBySerial(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("RPC call failed: %w", err)
// 	}

// 	if reply.Error != "" {
// 		return nil, fmt.Errorf("device lookup error: %s", reply.Error)
// 	}

// 	return reply, nil
// }
