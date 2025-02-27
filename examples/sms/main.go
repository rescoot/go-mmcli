// examples/sms/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rescoot/go-mmcli"
)

func main() {
	// Get first available modem ID
	id, err := mmcli.GetFirstModemID()
	if err != nil {
		log.Fatal("No modem found:", err)
	}

	fmt.Printf("Using modem %s\n\n", id)

	// Get messaging status
	fmt.Println("=== Messaging Status ===")
	status, err := mmcli.GetMessagingStatus(id)
	if err != nil {
		log.Printf("Failed to get messaging status: %v", err)
	} else {
		fmt.Printf("Default Storages: %s\n", strings.Join(status.DefaultStorages, ", "))
		fmt.Printf("Supported Storages: %s\n", strings.Join(status.SupportedStorages, ", "))
	}

	// List SMS messages
	fmt.Println("\n=== SMS Messages ===")
	messages, err := mmcli.ListSMS(id)
	if err != nil {
		log.Printf("Failed to list SMS messages: %v", err)
	} else {
		if len(messages) == 0 {
			fmt.Println("No SMS messages found")
		} else {
			fmt.Printf("Found %d SMS messages:\n", len(messages))
			for i, smsPath := range messages {
				fmt.Printf("%d: %s\n", i+1, smsPath)

				// Get SMS details
				parts := strings.Split(smsPath, "/")
				if len(parts) > 0 {
					smsID := parts[len(parts)-1]
					smsInfo, err := mmcli.GetSMSInfo(smsID)
					if err != nil {
						log.Printf("Failed to get SMS info: %v", err)
					} else {
						fmt.Printf("   Number: %s\n", smsInfo.Properties.Number)
						fmt.Printf("   Text: %s\n", smsInfo.Properties.Text)
						fmt.Printf("   Timestamp: %s\n", smsInfo.Properties.Timestamp)
						fmt.Printf("   State: %s\n", smsInfo.Properties.State)
					}
				}
			}
		}
	}

	// Create and send SMS if arguments provided
	if len(os.Args) == 3 {
		number := os.Args[1]
		text := os.Args[2]

		fmt.Printf("\n=== Sending SMS ===\n")
		fmt.Printf("To: %s\n", number)
		fmt.Printf("Message: %s\n", text)

		// Create SMS settings
		settings := mmcli.SMSCreateSettings{
			Number: number,
			Text:   text,
		}

		// Create SMS
		fmt.Println("Creating SMS...")
		smsID, err := mmcli.CreateSMS(id, settings)
		if err != nil {
			log.Fatalf("Failed to create SMS: %v", err)
		}
		fmt.Printf("SMS created with ID: %s\n", smsID)

		// Send SMS
		fmt.Println("Sending SMS...")
		err = mmcli.SendSMS(smsID)
		if err != nil {
			log.Fatalf("Failed to send SMS: %v", err)
		}
		fmt.Println("SMS sent successfully!")
	} else if len(os.Args) > 1 {
		fmt.Println("\nUsage: go run main.go [phone_number] [message]")
		fmt.Println("Example: go run main.go +1234567890 \"Hello, world!\"")
	}
}
