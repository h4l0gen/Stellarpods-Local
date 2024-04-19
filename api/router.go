package api

import (
	"fmt"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/j4ck4l-24/StellarPods/constants"
	log "github.com/sirupsen/logrus"
)

func StartRouter(){
	port := os.Getenv(constants.ENV_STELLARPODS_PORT)
	app := fiber.New()
	if err := app.Listen(fmt.Sprintf(":%s",port));err != nil{
		log.Fatal(err)
	}
}