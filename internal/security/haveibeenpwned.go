package security

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HaveIBeenPwnedClient struct {
	httpClient *http.Client
	userAgent  string
}

type PwnedPassword struct {
	Hash    string
	Count   int
	IsPwned bool
}

func NewHaveIBeenPwnedClient() *HaveIBeenPwnedClient {
	return &HaveIBeenPwnedClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: "whoami-auth-service/1.0",
	}
}

// CheckPassword checks if a password has been compromised using HaveIBeenPwned API
func (c *HaveIBeenPwnedClient) CheckPassword(password string) (*PwnedPassword, error) {
	// Hash the password with SHA-1
	hash := sha1.Sum([]byte(password))
	hashHex := strings.ToUpper(hex.EncodeToString(hash[:]))

	// Use k-anonymity: only send first 5 characters of hash
	prefix := hashHex[:5]
	suffix := hashHex[5:]

	// Make request to HaveIBeenPwned API
	url := fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", prefix)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Add-Padding", "true") // Add padding for k-anonymity

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		responseSuffix := parts[0]
		countStr := parts[1]

		if responseSuffix == suffix {
			count := 0
			fmt.Sscanf(countStr, "%d", &count)

			return &PwnedPassword{
				Hash:    hashHex,
				Count:   count,
				IsPwned: count > 0,
			}, nil
		}
	}

	// Password not found in breach database
	return &PwnedPassword{
		Hash:    hashHex,
		Count:   0,
		IsPwned: false,
	}, nil
}

// CheckPasswordWithRetry checks password with retry logic
func (c *HaveIBeenPwnedClient) CheckPasswordWithRetry(password string, maxRetries int) (*PwnedPassword, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		result, err := c.CheckPassword(password)
		if err == nil {
			return result, nil
		}

		lastErr = err
		time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
	}

	return nil, fmt.Errorf("failed after %d retries: %v", maxRetries, lastErr)
}
