package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/MAXIIIMVS/lets_go/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	debug          bool
	formDecoder    *form.Decoder
	logger         *slog.Logger
	sessionManager *scs.SessionManager
	snippets       models.SnippetModelInterface
	templateCache  map[string]*template.Template
	users          models.UserModelInterface
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	dbUsername := os.Getenv("MYSQL_USERNAME")
	dbPassword := os.Getenv("MYSQL_PASSWORD")

	if dbUsername == "" || dbPassword == "" {
		log.Fatal(
			"Failed to get the username and password for the database." +
				"\nPlease set them in a .env file. (e.g. MYSQL_USERNAME=web)",
		)
	}

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String(
		"dsn",
		dbUsername+":"+dbPassword+"@/snippetbox?parseTime=true",
		"MySQL data source name",
	)
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("✅ Connected to database")
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := application{
		debug:          *debug,
		formDecoder:    formDecoder,
		logger:         logger,
		sessionManager: sessionManager,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		users:          &models.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: log.New(
			os.Stderr,
			"ERROR:\t",
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		TLSConfig:    tlsConfig,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info(fmt.Sprintf("✅ Starting server on https://%s%s\n", "localhost", *addr))
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			db.Close()
		}
	}()
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
