package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SendOTPEmailAPI(targetEmail string, otpCode string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderName := os.Getenv("SENDER_NAME")

	htmlContent := fmt.Sprintf(`
		<html>
			<body>
				<h2>Verifikasi Akun</h2>
				<p>Berikut adalah kode OTP Anda:</p>
				<h1 style="color: #4CAF50; letter-spacing: 5px;">%s</h1>
				<p>Kode ini hanya berlaku selama 10 menit.</p>
			</body>
		</html>
	`, otpCode)

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  senderName,
			"email": senderEmail,
		},
		"to": []map[string]string{
			{
				"email": targetEmail,
			},
		},
		"subject":     "Kode Verifikasi OTP Anda",
		"htmlContent": htmlContent,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal mem-parsing JSON: %v", err)
	}

	url := "https://api.brevo.com/v3/smtp/email"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("gagal membuat request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("api-key", apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal mengirim request HTTP ke Brevo: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("API error (status %d): %s", res.StatusCode, string(bodyBytes))
	}

	return nil
}