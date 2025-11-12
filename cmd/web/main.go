package main

import (
	"embed"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/cxt314/drvc-go/internal/config"
	"github.com/cxt314/drvc-go/internal/driver"
	"github.com/cxt314/drvc-go/internal/handlers"
	"github.com/cxt314/drvc-go/internal/helpers"
	"github.com/cxt314/drvc-go/internal/models"
	"github.com/cxt314/drvc-go/internal/render"
	"github.com/pressly/goose/v3"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

//go:embed migrations/*.sql
var embedMigrations embed.FS

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	// run migrations in goose
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
        panic(err)
    }

    if err := goose.Up(db.SQL, "migrations"); err != nil {
        panic(err)
    }

	// start application
	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() (*driver.DB, error) {
	// What is going to be stored in the session
	gob.Register(models.User{})
	gob.Register(models.Member{})
	gob.Register(models.MemberAlias{})
	gob.Register(models.Vehicle{})
	gob.Register(models.MileageLog{})
	gob.Register(models.Trip{})
	gob.Register(models.Rider{})

	// read environment variables
	inProduction := os.Getenv("IS_PRODUCTION")
	useCache := os.Getenv("USE_CACHE")
	dbURL := os.Getenv("DATABASE_URL")

	// set from production env variable
	if inProduction == "false" {
		app.InProduction = false
	} else {
		app.InProduction = true
	}

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// initialize sessions
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(dbURL)

	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	// set app.UseCache from env variables
	if useCache == "false" {
		app.UseCache = false
	} else {
		app.UseCache = true
	}

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
