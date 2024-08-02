// this is the go server. Which serves a static file called index.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
)

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func main() {
	fmt.Println("Starting server on port 80")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/send-email", handleSendEmail)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("listenAndServe failed: %v", err)
	}
}

func handleSendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var form ContactForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = sendEmail(form)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Email sent successfully")
}

func sendEmail(form ContactForm) error {
	from := "your-email@example.com"
	password := "your-email-password"
	to := "your-email@example.com"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "New Contact Form Submission"
	body := fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s", form.Name, form.Email, form.Message)
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}
