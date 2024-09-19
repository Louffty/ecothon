package main

import (
	"github.com/Louffty/green-code-moscow/cmd/app"
	"github.com/Louffty/green-code-moscow/internal/adapters/config"
	"github.com/Louffty/green-code-moscow/internal/adapters/controller/api/setup"
	"github.com/Louffty/green-code-moscow/internal/domain/scheduler"
)

// main is the entry point of the application.
func main() {
	appConfig := config.GetConfig()
	bizkitEduApp := app.NewBizkitEduApp(appConfig)

	conferenceScheduler := scheduler.NewConferenceScheduler(bizkitEduApp)
	conferenceScheduler.Start()

	setup.Setup(bizkitEduApp)
	bizkitEduApp.Start()
}
