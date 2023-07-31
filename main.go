package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jnka9755/go-05USER/internal/user"
	"github.com/jnka9755/go-05USER/package/boostrap"
	"github.com/joho/godotenv"
)

func main() {

	router := mux.NewRouter()

	_ = godotenv.Load()

	log := boostrap.InitLooger()

	db, err := boostrap.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		log.Fatal("paginator limit defauly is required")
	}

	userRepository := user.NewRepository(log, db)
	userBusiness := user.NewBusiness(log, userRepository)
	userController := user.MakeEndpoints(userBusiness, user.Config{LimPageDef: pagLimDef})

	router.HandleFunc("/users", userController.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userController.Get).Methods("GET")
	router.HandleFunc("/users", userController.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userController.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userController.Delete).Methods("DELETE")

	port := os.Getenv("PORT")

	address := fmt.Sprintf("127.0.0.1:%s", port)

	server := http.Server{
		Handler:      http.TimeoutHandler(router, time.Second*5, "Timeout!"),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	error := server.ListenAndServe()

	if err != nil {
		log.Fatal(error)
	}
}
