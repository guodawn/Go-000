package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

var s = &http.Server{
	Addr: fmt.Sprintf(":%d", HTTP_PORT),
}
var groupFinish = make(chan int)

func main() {

	ctx, cancle := context.WithTimeout(context.Background(), time.Second*1)
	defer cancle()
	go startHttpServer()
	go dealBiz(ctx)
	gracefulShutdown()

}

const HTTP_PORT = 8088

func startHttpServer() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello")
	})
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("s.ListenAndServe err:%v", err)
	}
}

func dealBiz(ctx context.Context) {
	g := new(errgroup.Group)
	var urls = []string{
		"https://www.baidu.com",
		//		"https://www.google.com",
		//"https://www.unknown.com"
	}
	for _, url := range urls {
		tUrl := url
		g.Go(func() error {
			resp, err := http.Get(tUrl)
			time.Sleep(time.Second * 5)
			if err == nil {
				defer resp.Body.Close()
			}
			return err
		})
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				//		groupFinish <- 1
			default:
				time.Sleep(time.Second * 1)
				fmt.Println("watching")
			}
		}
	}()
	if err := g.Wait(); err == nil {
		groupFinish <- 1
		fmt.Println("Success fetch all urls")
	}
}
func gracefulShutdown() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down")
	<-groupFinish
	fmt.Println("finish all jobs")
	ctx := context.Background()
	err := s.Shutdown(ctx)
	if err != nil {
		log.Fatalf("s.Shutdown err:%v", err)
	}
}
