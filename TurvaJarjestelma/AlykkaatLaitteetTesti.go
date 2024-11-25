package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/gorilla/mux"
)

// Configure Email
const (
	EmailSender    = "trustonkuos@hotmail.com"
	EmailPassword  = ""
	EmailRecipient = "roope.kuossari@centria.fi"
	SMTPServer     = "smtp-mail.outlook.com"
	SMTPPort       = 587
)

// Handle image upload
func handleImageUpload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	fileType := http.DetectContentType(buffer)
	if fileType != "image/jpeg" {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}
	file.Seek(0, 0) // Reset file pointer to the beginning after checking

	// Save the image locally
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	imagePath := fmt.Sprintf("received_images/%s.jpg", timestamp)
	os.MkdirAll("received_images", os.ModePerm)
	out, err := os.Create(imagePath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to write image", http.StatusInternalServerError)
		return
	}

	fmt.Println("Image received:", imagePath)

	// Analyze the image (mock AI processing)
	isIntruder := analyzeImage(imagePath)

	// Respond to the sender
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Image received"))

	// Handle intruder alert
	if isIntruder {
		sendEmail(imagePath)
	}
}

// Mock AI analysis
func analyzeImage(imagePath string) bool {
	fmt.Println("Analyzing image:", imagePath)
	// Replace with actual AI analysis (e.g., TensorFlow or external API)
	// For now, return true to simulate detecting an intruder
	return true
}

// Send email alert
func sendEmail(imagePath string) {
	m := gomail.NewMessage()
	m.SetHeader("From", EmailSender)
	m.SetHeader("To", EmailRecipient)
	m.SetHeader("Subject", "Intruder Alert!")
	m.SetBody("text/plain", "An intruder has been detected. See attached image.")
	m.Attach(imagePath)

	d := gomail.NewDialer(SMTPServer, SMTPPort, EmailSender, EmailPassword)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}

func main() {
	os.MkdirAll("received_images", os.ModePerm) // Ensure directory exists
	r := mux.NewRouter()
	r.HandleFunc("/upload", handleImageUpload).Methods("POST")

	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}
