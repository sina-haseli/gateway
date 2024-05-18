package main

import (
	"fmt"
	"gateway/pkg/mocks"
	"github.com/stretchr/testify/mock"
	"log"
	"net/http"
)

func main() {
	service1Mock1 := new(mocks.ServiceOneMock)
	//service1Mock2 := new(mocks.ServiceOneMock)
	service1Mock3 := new(mocks.ServiceOneMock)

	setupMockServiceOne(service1Mock1, 8081)
	//setupMockServiceOne(service1Mock2, 8082)
	setupMockServiceOne(service1Mock3, 8083)

	select {} // Keep the program running
}

func setupMockServiceOne(src *mocks.ServiceOneMock, port int) {
	src.On("ServeHTTP", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service One Response" + " " + fmt.Sprint(port)))
	})

	mux := http.NewServeMux()
	mux.Handle("/", src)

	go func() {
		log.Printf("Service One Mock running on port %d", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
			log.Fatalf("Failed to start Service One Mock on port %d: %v", port, err)
		}
	}()
}
