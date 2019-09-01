package healthchecker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Services which contains list of services
type Services struct {
	Services []Service `json:"Services"`
}

// Service with data which contains data who will be checked
type Service struct {
	Name        string `json:"Name"`
	Protocol    string `json:"Protocol"`
	URI         string `json:"URI"`
	Path        string `json:"Path"`
	Port        int    `json:"Port"`
	Method      string `json:"Method"`
	RequestTest string `json:"RequestTest"`
	StatusTime  time.Time
	StatusText  string
}

func doCheckService(service *Service) {

	service.StatusTime = time.Now()

	var completeURL string
	completeURL = fmt.Sprintf("%s://%s:%d/%s", service.Protocol, service.URI, service.Port, service.Path)

	var resp *http.Response
	var err error

	if service.Method == "POST" {
		var jsonStr = []byte(service.RequestTest)

		req, err := http.NewRequest(service.Method, completeURL, bytes.NewBuffer(jsonStr))
		if err != nil {
			service.StatusText = fmt.Sprint(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)

	} else {
		resp, err = http.Get(completeURL)
	}

	if resp != nil {
		defer resp.Body.Close() // Memory management
	}

	if err != nil {
		service.StatusText = fmt.Sprint(err)
		return
	}

	service.StatusText = strconv.Itoa(resp.StatusCode) + " " + http.StatusText(resp.StatusCode)

	//=================================================
	// READ THE BODY EVEN THE DATA IS NOT IMPORTANT
	// THIS MUST TO DO, TO AVOID MEMORY LEAK WHEN REUSING HTTP
	// CONNECTION
	//=================================================
	_, err = io.Copy(ioutil.Discard, resp.Body) // WE READ THE BODY
	if err != nil {
		return
	}
}

func doHealthCheck(services []Service) {
	for i := range services {
		service := &services[i]
		service.StatusTime = time.Now()
		service.StatusText = "checking..."
		go doCheckService(service)
	}
}

func startHTTPServer(services *Services) {
	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))

		services = &Services{}
		materializeServices(services)

		w.Write([]byte("\nReloaded.\n"))

		doHealthCheck(services.Services)
		w.Write([]byte("\nChecking..."))
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
		doHealthCheck(services.Services)
		w.Write([]byte("\nChecking..."))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String() + "\n"))
		listServices, _ := json.MarshalIndent(services, "", "  ")
		w.Write(listServices)
	})

	fmt.Println("Listening at port 3001, at " + time.Now().String())
	httpErr := http.ListenAndServe(":3001", nil)

	if httpErr != nil {
		log.Fatal(httpErr)
	}
}

func materializeServices(services *Services) {
	jsonFile, err := os.Open("services.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &services)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Services loaded.")
	return
}

func main() {
	var services Services
	materializeServices(&services)

	doHealthCheck(services.Services)

	startHTTPServer(&services)

}
