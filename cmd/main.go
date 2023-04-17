package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jonathannavas/gocourse_enrollment/internal/enrollment"
	"github.com/jonathannavas/gocourse_enrollment/pkg/bootstrap"
	"github.com/jonathannavas/gocourse_enrollment/pkg/handler"

	courseSdk "github.com/jonathannavas/go_course_sdk/course"
	userSdk "github.com/jonathannavas/go_course_sdk/user"
)

func main() {

	_ = godotenv.Load()

	logs := bootstrap.InitLogger()
	db, err := bootstrap.DBConnection()

	if err != nil {
		logs.Fatal(err)
	}

	pagLimitDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimitDef == "" {
		logs.Fatal("paginator limit default is required")
	}

	courseTrans := courseSdk.NewHttpClient(os.Getenv("API_COURSE_URL"), "")
	userTrans := userSdk.NewHttpClient(os.Getenv("API_USER_URL"), "")

	ctx := context.Background()
	enrollmentRepository := enrollment.NewRepo(logs, db)
	enrollmentService := enrollment.NewService(logs, enrollmentRepository, courseTrans, userTrans)
	h := handler.NewUserHTTPServer(ctx, enrollment.MakeEndpoints(enrollmentService, enrollment.Config{LimitPageDef: pagLimitDef}))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		// Handler:      http.TimeoutHandler(router, time.Second*3, "Timeout!!"),
		Handler:      accessControl(h),
		Addr:         address,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	errCh := make(chan error)

	go func() {
		log.Println("Listen in", address)
		errCh <- srv.ListenAndServe()
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
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type, DNT")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
