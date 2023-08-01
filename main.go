package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jnka9755/go-05USER/internal/user"
	"github.com/jnka9755/go-05USER/package/boostrap"
	"github.com/jnka9755/go-05USER/package/handler"
	"github.com/joho/godotenv"
)

func main() {

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

	ctx := context.Background()

	userRepository := user.NewRepository(log, db)
	userBusiness := user.NewBusiness(log, userRepository)

	handler := handler.NewUserHTTPServer(ctx, user.MakeEndpoints(userBusiness, user.Config{LimPageDef: pagLimDef}))

	// router.HandleFunc("/users", userController.Create).Methods("POST")
	// router.HandleFunc("/users/{id}", userController.Get).Methods("GET")
	// router.HandleFunc("/users", userController.GetAll).Methods("GET")
	// router.HandleFunc("/users/{id}", userController.Update).Methods("PATCH")
	// router.HandleFunc("/users/{id}", userController.Delete).Methods("DELETE")

	port := os.Getenv("PORT")

	address := fmt.Sprintf("127.0.0.1:%s", port)

	server := http.Server{
		Handler:      accessControl(handler),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errCh := make(chan error)

	go func() {
		log.Println("Listen in ", address)
		errCh <- server.ListenAndServe()
	}()

	err = <-errCh

	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
