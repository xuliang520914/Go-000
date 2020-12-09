package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main()  {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		log.Println("global defer")
		cancel()
	}()
	g, _ := errgroup.WithContext(ctx)
	// server
	g.Go(func() error {
		return runServer(ctx,8001,"正式server")
	})
	// debug server
	g.Go(func() error {
		return runServer(ctx,8111,"debug server")
	})
	c := make(chan os.Signal)
	signal.Notify(c,os.Interrupt)
	g.Go(func() error {
		return listenStopSignal(ctx,c,cancel)
	})
	log.Printf("Actual pid is %d", syscall.Getpid())
	if err := g.Wait(); err != nil {
		log.Println("g.wait finish")
		log.Printf("%+v",err)
	}
	log.Println("sever finish")
	time.Sleep(2 * time.Second)

}

func listenStopSignal(ctx context.Context, sig chan os.Signal, cancel context.CancelFunc) error {
	select {
	case sign := <- sig:
		log.Println("接受信号",sign)
		err := fmt.Errorf("handle signal: %d", sign)
		cancel()
		return err
	case <-ctx.Done():
		log.Println("listenStopSignal ctx.Done")
		return nil
	}

}
func runServer(ctx context.Context, serverPort int, serverName string) error {
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", serverPort),
		MaxHeaderBytes: 1 << 20,
	}

	defer func() {
		if err := recover(); err != nil {
			log.Printf("%+v",err)
			log.Println(serverName,"into defer recover()")
			cancelServer(s,serverName)
		}
	}()
	go func() {
		<-ctx.Done()
		cancelServer(s,serverName)
	}()
	log.Println(serverPort,serverName," with run")
	err := s.ListenAndServe()
	log.Println(serverName,"with error",err)
	return err
}

func cancelServer(s *http.Server, serverName string) {
	log.Println(serverName,"stop")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(timeoutCtx)
}