package main

import (
	"log"
	"net/http"
	emp "test_app_employees/Employees"

	"github.com/gorilla/mux"
	"github.com/kardianos/service"
	"github.com/rs/cors"
)

const ApplicationName = "Employees"
const ApplicationReleaseNumber = "0.1.0"
const ApplicationReleaseDate = "11/05/2026" //"dd/MM/YYYY"

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	log.Println("-------------------------------------------------------------------")
	log.Println("Application name: ", ApplicationName)
	log.Println("Application release number: ", ApplicationReleaseNumber)
	log.Println("Application release date: ", ApplicationReleaseDate)
	log.Println("-------------------------------------------------------------------")

	//read and parse configuration
	err := emp.GetDefaultConfiguration()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Establishing connections...")
	UserControl := emp.NewUserControl()
	// UserControl.InitializeDAO()
	// UserControl.InitializeCache()
	// UserControl.IndexesMaintenanceProcess()

	err = UserControl.InitializeMongoxRepositories()
	if err != nil {
		log.Fatal("failed to create mongodb reporitories: ", err)
	}
	UserControl.RedisDataLoader()
	//Add user routers to the web service
	log.Println("Add routers to the web service")
	router := mux.NewRouter().StrictSlash(true)
	UserControl.AddToRouter(router, UserControl)
	HttpServicePort := emp.Configuration.HttpServicePort
	log.Println("HTTP listen and serve on port " + HttpServicePort) //auc.Configuration.HttpServicePort
	//log.Fatal(http.ListenAndServe(":"+HttpServicePort, router))

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, //you service is available and allowed for this base url
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application

		},
	})
	log.Fatal(http.ListenAndServe(":"+HttpServicePort, corsOpts.Handler(router)))

}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        ApplicationName,
		DisplayName: ApplicationName,
		Description: ApplicationName + " service",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
