package initializer

import (
	"time"
)

func Init() {
	// This is the entry point of the application.
	// It initializes the application and starts the server.
	// It is the equivalent of the main function in Go.

	// Set the timezone to JST.
	jst, _ := time.LoadLocation("Asia/Tokyo")
	time.Local = jst

}
