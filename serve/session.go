package serve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type SessionResponse struct {
	SessionID string `json:"session_id"`
}

func GetSessionS() (string, error) {
	// Generate a UUID for the device ID
	deviceID := uuid.New().String()

	// Create the request payload
	sessionProperties := map[string]string{"deviceId": deviceID}
	payload := map[string]interface{}{
		"session_properties": sessionProperties,
		"session_type":       1,
	}
	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make the HTTP POST request
	url := "https://studio-api.prod.suno.com/api/user/create_session_id/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)

	log.Println("response body:", string(body))
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response code: %d, body: %s", resp.StatusCode, string(body))
	}

	var sessionResp SessionResponse
	err = json.Unmarshal(body, &sessionResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Log and return the session ID
	log.Println("Session ID:", sessionResp.SessionID)
	return sessionResp.SessionID, nil
}
