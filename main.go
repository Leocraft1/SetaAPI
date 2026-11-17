package main

import (
	"log"
	"net/http"
	"setaapi/config"
	"setaapi/internal/gtfs"
	"setaapi/internal/handler"
	"setaapi/internal/repository"
	"setaapi/internal/scheduler"
	"setaapi/internal/sse"
)


func main() {
	//Init config
	config.LoadConf()
	//DBs Init
	repository.InitMezzi()
	repository.InitStatic()
	err := gtfs.UpdateStatic()
	if err != nil {
		log.Fatal(gtfs.UpdateStatic())
	}
	// Initialize server-sent events for frontend synchronization.
	sse.InitHub()

	if config.IS_PRIMARY {
		if err := scheduler.Start(true); err != nil {
			log.Fatalf("Failed to start scheduler: %v", err)
		}
		defer scheduler.Stop()
	} else {
		log.Println("IS_PRIMARY set to false, not updating static GTFS data.")
		if err := scheduler.Start(false); err != nil {
			log.Fatalf("Failed to start scheduler: %v", err)
		}
		defer scheduler.Stop()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthCheckHandler)/*
	mux.HandleFunc("GET /arrivals/{id}", handler.ArrivalsHandler)
	mux.HandleFunc("GET /busesinservice", handler.BusesinserviceHandler)
	mux.HandleFunc("GET /vehicleinfo/{id}", handler.VehicleinfoHandler)
	mux.HandleFunc("GET /linelist", handler.LinelistHandler)
	mux.HandleFunc("GET /modelslist", handler.ModelslistHandler)
	mux.HandleFunc("GET /stops", handler.StoplistHandler)
	mux.HandleFunc("GET /updates", handler.UpdatesHandler)
	mux.HandleFunc("GET /routecodes", handler.RoutecodesHandler)
	mux.HandleFunc("GET /routestops/{id}", handler.RoutestopsHandler)
	mux.HandleFunc("GET /nextstops/{id}", handler.NextstopsHandler)
	mux.HandleFunc("GET /allnews", handler.AllnewsHandler)
	mux.HandleFunc("GET /news", handler.NewsHandler)
	mux.HandleFunc("GET /lineproblems", handler.LineproblemsHandler)
	mux.HandleFunc("GET /lineproblems/{id}", handler.LineproblemHandler)
	mux.HandleFunc("GET /timetable", handler.TimetableHandler)
	mux.HandleFunc("GET /routemap/{id}", handler.RoutemapHandler)
	mux.HandleFunc("GET /assignments", handler.AssignmentsHandler)

	//OAuth
	if config.ENABLE_AUTH {
		ctx := context.Background()

		provider, err := oidc.NewProvider(ctx, config.OIDC_ISSUER_URL)
		if err != nil {
			log.Fatal("impossibile contattare Authentik per il discovery OIDC:", err)
		}

		verifier := provider.Verifier(&oidc.Config{ClientID: config.OIDC_CLIENT_ID})
		authMiddleware := handler.RequireAuth(verifier, config.OIDC_ALLOWED_GROUPS)

		mux.Handle("POST /assignments/add", authMiddleware(http.HandlerFunc(handler.AddAssignmentHandler)))
		mux.Handle("PUT /assignments/changevehicle", authMiddleware(http.HandlerFunc(handler.UpdateAssignmentHandler)))
		mux.Handle("DELETE /assignments/remove", authMiddleware(http.HandlerFunc(handler.DeleteAssignmentHandler)))
	}*/

	finalHandler := handler.CorsMiddleware()(mux) //wraps mux

	//Listen on port and start API
	log.Println("Server started on port " + config.PORT)
	log.Print(http.ListenAndServe(config.PORT, finalHandler))
}
