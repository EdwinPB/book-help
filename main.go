package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"image/png"

	"github.com/fogleman/gg"
)

type Question struct {
	ID       int
	Name     string
	Question string
}

var (
	questions     []Question
	questionMutex sync.RWMutex
)

func main() {
	fmt.Println("------>")

	// Serve static files from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// Handle root URL
	http.HandleFunc("/", indexHandler)
	// Handle download URL
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/quote", quoteHandler)

	// Graceful server shutdown
	server := &http.Server{Addr: ":1010"}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %s\n", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))

	questionMutex.RLock()
	defer questionMutex.RUnlock()

	questions = append(questions, Question{Name: "Edwin"})

	if err := tmpl.Execute(w, questions); err != nil {
		fmt.Printf("Error rendering template: %s\n", err)
		// Optionally handle the error, e.g., return an internal server error response
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Generate or read your SVG content
	svgContent := `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100"><circle cx="50" cy="50" r="40" fill="red" /></svg>`

	// Convert SVG to PNG
	img, err := convertSVGToPNG(svgContent)
	if err != nil {
		http.Error(w, err.Error()+" Internal Server Error vi", http.StatusInternalServerError)
		return
	}

	// Set Content-Disposition header to trigger a download
	w.Header().Set("Content-Disposition", "attachment; filename=downloaded.png")
	w.Header().Set("Content-Type", "image/png")

	// Write the PNG image to the response
	err = png.Encode(w, img)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	quote, err := makeOpenAIRequest("")
	fmt.Println("-------->", quote)
	fmt.Println("-------->", err)
}

func convertSVGToPNG(svgContent string) (image.Image, error) {
	// Create a new drawing context
	const (
		width  = 100
		height = 100
	)
	dc := gg.NewContext(width, height)

	// Draw the SVG onto the drawing context
	dc.SetRGB(1, 1, 1) // Set background color to white
	dc.Clear()
	dc.SetRGB(0, 0, 0)                                       // Set drawing color to black
	dc.LoadFontFace("static/Poppins-Bold.ttf", 12)           // Adjust the font path
	dc.DrawStringAnchored("Danier", width, height, 0.5, 0.5) // Draw text

	// Resize the image if needed
	resizedImage := dc.Image()

	return resizedImage, nil
}

func makeOpenAIRequest(phrase string) (string, error) {
	phrase = "Now is the time. Because now is the only time you have."
	if phrase == "" {
		println("No prhase include")
		return "", nil
	}
	openaiURL := "https://api.openai.com/v1/chat/completions"
	apiKey := "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"

	libro := "Discipline is destiny"
	author := "Ryan Holiday"
	prompt := "Realiza una cita profesional para la siguiente frase entre llaves, {" + phrase + "} del libro, " + libro + ", " + author + ". Antes revisa que si sea una frase real del libro y en caso no lo sea, solo respondeme entre llaves {no_found} en caso si lo sea encierra la cita con el formato profesional en llaves {}"
	requestData := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You work as an API that will return a phrase that the user takes from a book to get the correct quote for that phrase.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"model": "gpt-3.5-turbo",
		// Add any other parameters based on your requirements
	}

	jsonBody, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}
