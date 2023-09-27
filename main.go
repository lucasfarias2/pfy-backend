package main

import (
	"log"
	"github.com/joho/godotenv"
	"os"

	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/gofiber/fiber/v2"
)

type DatabaseRequest struct {
	DBName string `json:"dbname"`
}

func createDatabase(c *fiber.Ctx) error {
	req := new(DatabaseRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	// Initialize a session for AWS
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create RDS service client
	svc := rds.New(sess)

	// Create the RDS instance
	input := &rds.CreateDBInstanceInput{
		AllocatedStorage:     aws.Int64(20),
		DBInstanceIdentifier: aws.String(req.DBName),
		DBInstanceClass:      aws.String("db.t2.micro"),
		Engine:               aws.String("postgres"),
		MasterUsername:       aws.String("username"),
		MasterUserPassword:   aws.String("password"),
	}

	result, err := svc.CreateDBInstance(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	
	app := fiber.New()

	region := os.Getenv("AWS_REGION")
    fmt.Println("AWS_REGION:", region)

	app.Get("/api/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/api/database", createDatabase)

	log.Println("Server started at :8080")
	app.Listen(":8080")
}

