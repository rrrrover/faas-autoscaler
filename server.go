package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/openfaas/faas-provider/auth"
	gtwHandlers "github.com/openfaas/faas/gateway/handlers"
	"github.com/rrrrover/faas-autoscaler/handlers"
	"github.com/rrrrover/faas-autoscaler/scaler"
	"github.com/rrrrover/faas-autoscaler/types"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

const (
	defaultFunctionsProviderUrl = "http://127.0.0.1:8080/"
	defaultAuthUser             = "admin"
	defaultAuthPassword         = "admin"
)

func main() {
	functionsProviderUrlPath := os.Getenv("functions_provider_url")
	if len(functionsProviderUrlPath) == 0 {
		functionsProviderUrlPath = defaultFunctionsProviderUrl
	}
	functionsProviderUrl, err := url.Parse(functionsProviderUrlPath)
	if err != nil {
		log.Fatalf("error occurred when parsing url: %v", err)
	}

	credentials := auth.BasicAuthCredentials{User: defaultAuthUser, Password: defaultAuthPassword}
	secretMountPath := "/var/secrets/"
	if val, ok := os.LookupEnv("secret_mount_path"); ok && len(val) > 0 {
		secretMountPath = val
	}
	if val, err := readFile(path.Join(secretMountPath, "basic-auth-user")); err == nil {
		credentials.User = val
	} else {
		log.Printf("Unable to read username: %s", err)
	}
	if val, err := readFile(path.Join(secretMountPath, "basic-auth-password")); err == nil {
		credentials.Password = val
	} else {
		log.Printf("Unable to read password: %s", err)
	}
	injector := gtwHandlers.BasicAuthInjector{Credentials: &credentials}
	scaler := scaler.NewAutoScaler(*functionsProviderUrl, injector)

	handlerSet := new(types.AutoScaleHandlerSet)
	handlerSet.ScaleUpHandlerFunc = handlers.NewScaleUpHandlerFunc(scaler)
	handlerSet.ScaleDownHandlerFunc = handlers.NewScaleDownHandlerFunc(scaler)

	r := mux.NewRouter()
	r.HandleFunc("/system/scale-up", handlerSet.ScaleUpHandlerFunc).Methods(http.MethodPost)
	r.HandleFunc("/system/scale-down", handlerSet.ScaleDownHandlerFunc).Methods(http.MethodPost)

	port := 8081
	readTimeout := 10 * time.Second
	writeTimeout := 10 * time.Second
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		Handler:        r,
	}
	log.Fatal(server.ListenAndServe())
}

func readFile(path string) (string, error) {
	_, err := os.Stat(path)
	if err == nil {
		data, readErr := ioutil.ReadFile(path)
		return strings.TrimSpace(string(data)), readErr
	}
	return "", err
}
