package main

import (
	"log"

	"final-project-backend/db"
	"final-project-backend/internal/handlers"
	"final-project-backend/internal/middlewares"
	"final-project-backend/internal/repositories"
	"final-project-backend/internal/routes"
	"final-project-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	db := db.Get()
	rp := repositories.New(db)
	s := services.New(rp)
	h := handlers.New(s)

	r := gin.Default()
	r.Use(middlewares.Cors)
	routes.InitRoutes(r, h)

	err = r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
