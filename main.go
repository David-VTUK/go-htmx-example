package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

var tpl *template.Template

// DogAPIResponse represents the JSON response from the dog.ceo API
type DogAPIResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func init() {
	tpl = template.Must(template.ParseFiles("templates/index.html"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/getImage", imageHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	//query := r.URL.Query().Get("query")

	url := "https://dog.ceo/api/breeds/image/random"
	client := &http.Client{Timeout: 10 * time.Second}

	// Make the API request
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	// Decode the JSON response into the DogAPIResponse struct
	var apiResponse DogAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Fatalf("Failed to decode JSON response: %v", err)
	}

	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "image", apiResponse)
}
