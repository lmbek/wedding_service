package main

import (
	"log"
	"wedding_service/env"
	"wedding_service/flags"
	"wedding_service/webserver"
	"wedding_service/webserver/website/frontend"
)

var mainWebserver webserver.Webserver

// Counter for the id field
var logCounter int

func main() {
	// Start the log routine in a separate goroutine
	//go logRoutine()

	// Load flags and initialize the environment
	flags.LoadFrontendFlag()

	// Initialize the environment
	if err := initEnv(); err != nil {
		log.Fatalf("Error initializing environment: %v", err)
	}

	// Create and start the webserver
	if err := startWebserver(); err != nil {
		log.Fatalf("Error starting webserver: %v", err)
	}
}

func initEnv() error {
	if err := env.Init(); err != nil {
		log.Printf("Error initializing environment: %v", err)
		return err
	}
	return nil
}

func startWebserver() error {
	newFrontend := frontend.NewFrontend()
	var err error
	// Directly assign to the global variable
	mainWebserver, err = createWebserver(newFrontend)
	if err != nil {
		log.Printf("Error creating webserver: %v", err)
		return err
	}

	// Start the webserver
	if err := mainWebserver.ListenAndServe(); err != nil {
		log.Printf("Error starting webserver: %v", err)
		return err
	}

	return nil
}

func createWebserver(newFrontend frontend.Frontend) (webserver.Webserver, error) {
	w, err := webserver.NewWebserver(newFrontend)
	if err != nil {
		log.Printf("Error creating webserver: %v", err)
		return nil, err
	}

	return w, nil
}

//
//// logRoutine simulates logging of response time in JSON format every 2 seconds
//func logRoutine() {
//	for {
//		// Increment the logCounter for the id field
//		logCounter++
//
//		// Current timestamp
//		timestamp := time.Now()
//
//		// Convert timestamp to nanoseconds
//		tsNs := timestamp.UnixNano()
//
//		// Construct the log entry
//		logEntry := map[string]interface{}{
//			"labels": map[string]string{
//				"detected_level": "unknown", // Adjust based on your log level detection logic
//				"filename":       "/var/lib/docker/containers/296f5ac2a5e1fc7692edccf3a962b079c2d36a42ab195f0f727d1ff27b4b450f/296f5ac2a5e1fc7692edccf3a962b079c2d36a42ab195f0f727d1ff27b4b450f-json.log",
//				"job":            "go-logs",
//				"service_name":   "go-logs",
//			},
//			"Time": timestamp.Format("2006-01-02 15:04:05.000"),                                                                                                                                  // Time formatted to match the required format
//			"Line": `{"log":"2025-07-08 12:15:46+00:00 [Note] [Entrypoint]: Entrypoint script for MySQL Server 8.4.5-1.el9 started.","stream":"stdout","time":"2025-07-08T12:15:46.127048495Z"}`, // The log content
//			"tsNs": tsNs,                                                                                                                                                                         // The nanosecond timestamp
//			"id":   fmt.Sprintf("%d_%d", tsNs, logCounter),                                                                                                                                       // Combine tsNs and logCounter to create a unique ID
//		}
//
//		// Log the entry as JSON
//		logData, err := json.Marshal(logEntry)
//		if err != nil {
//			log.Fatalf("Error marshalling log entry: %v", err)
//		}
//		log.Println(string(logData))
//
//		// Sleep for 10 seconds before the next log entry
//		time.Sleep(10 * time.Second)
//	}
//}
