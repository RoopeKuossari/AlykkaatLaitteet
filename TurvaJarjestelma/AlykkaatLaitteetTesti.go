package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	// Telegram bot token
	BotToken      = "7870042087:AAEUFhDb_VnsnPNE0aW3d15E7qmEnDnP1b0"
	AdminChatID   = 5827699039 // Replace with your Telegram Chat ID
	ImageSavePath = "received_images"
)

func main() {
	// Initialize the Telegram bot
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Set up the HTTP server for image uploads
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		handleImageUpload(bot, w, r)
	})

	// Create the directory to save images
	os.MkdirAll(ImageSavePath, os.ModePerm)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handle image upload via HTTP
func handleImageUpload(bot *tgbotapi.BotAPI, w http.ResponseWriter, r *http.Request) {
	// Parse uploaded image
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file type
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	fileType := http.DetectContentType(buffer)
	if fileType != "image/jpeg" {
		http.Error(w, "Invalid file type. Only JPEG is allowed.", http.StatusBadRequest)
		return
	}
	file.Seek(0, 0) // Reset file pointer to the beginning

	// Save the image locally
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	imagePath := fmt.Sprintf("%s/%s.jpg", ImageSavePath, timestamp)

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

	log.Printf("Image received and saved: %s", imagePath)

	// Analyze the image for intruders
	isIntruder := analyzeImage(imagePath)

	// Respond to the sender
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Image received and processed."))

	// Alert via Telegram if intruder detected
	if isIntruder {
		sendIntruderAlert(bot, AdminChatID, imagePath)
	}
}

// Analyze the image (mock AI processing)
func analyzeImage(imagePath string) bool {
	fmt.Println("Analyzing image:", imagePath)
	// Replace this with real AI processing logic
	// Simulate detecting an intruder (true)
	return true
}

// Send intruder alert via Telegram
func sendIntruderAlert(bot *tgbotapi.BotAPI, chatID int64, imagePath string) {
	// Send alert message
	alertMessage := tgbotapi.NewMessage(chatID, "Intruder detected! Here is the photo:")
	if _, err := bot.Send(alertMessage); err != nil {
		log.Printf("Failed to send alert message: %v", err)
		return
	}

	// Send the intruder photo
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(imagePath))
	if _, err := bot.Send(photo); err != nil {
		log.Printf("Failed to send photo: %v", err)
		return
	}

	log.Printf("Intruder alert sent successfully with photo: %s", imagePath)
}
