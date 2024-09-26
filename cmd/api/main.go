package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"backend/internal/repository/mailrepo"
	"backend/internal/repository/storagerepo"
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const port = 8080

type application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	Storage      repository.StorageRepo
	EM           repository.MailRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

func main() {
	// set application config
	var app application

	// load environment variables
	if os.Getenv("RUNNING_IN_DOCKER") != "true" {
		// Load .env file if not running in Docker
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	// After either loading from the .env file or relying on passed environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	awsRegion := os.Getenv("AWS_REGION")
	fromAddress := os.Getenv("FROM_ADDRESS")

	// read from command line
	flag.StringVar(&app.DSN, "dsn", mongoURI, "MongoDB connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.Parse()

	// connect to the database
	conn, err := app.connectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	//app.DB = &dbrepo.MongoDBRepo{
	//	Client:   conn,
	//	Database: "Share2Teach",
	//}
	app.DB = dbrepo.NewMongoDBRepo(conn, "Share2Teach")

	defer func() {
		if err := conn.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "__Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	// initialize s3
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	sesClient := ses.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	// Assign both clients to your storage repository struct
	app.Storage = &storagerepo.StorageRepo{
		S3Client:      s3Client,
		PresignClient: presignClient,
	}

	app.EM = &mailrepo.MailRepo{
		SESClient:   sesClient,
		FromAddress: fromAddress,
	}

	log.Println("Starting application on port", port)

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
