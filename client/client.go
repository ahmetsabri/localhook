package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey, err := getOrCreateAPIKey()
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(os.Getenv("CLIENT_CONNECT_URL"), apiKey)

	forwardingUrl := flag.String("f", os.Getenv("CLIENT_DEFAULT_FORWARDING_URL"), "The local URL to forward the request to ")
	flag.Parse()

	color.Cyan("incoming responses will be forwarded to : %s", *forwardingUrl)
	color.Yellow("webhook url : "+os.Getenv("CLIENT_WEBHOOK_URL"), apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return
	}

	reader := bufio.NewReader(resp.Body)

	var eventBuffer bytes.Buffer

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by server")
				break
			}
			fmt.Println("Error reading from stream:", err)
			break
		}

		if len(line) == 1 && line[0] == '\n' {
			color.Cyan("Response received!")
			body := eventBuffer.String()
			color.Yellow(body)
			forwardResponse(*forwardingUrl, body)
			eventBuffer.Reset()
			continue
		}

		eventBuffer.Write(line)
	}

	fmt.Println("Disconnected")
}

func forwardResponse(url string, body string) error {
	color.Cyan("Forwarding Request to: %s\n", url)

	client := &http.Client{}

	jsonBody := bytes.NewBuffer([]byte(body))

	req, err := http.NewRequest("POST", url, jsonBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		color.Green("Response Status: %s\n", resp.Status)
	} else {
		color.Red("Response Status: %s\n", resp.Status)
	}
	fmt.Printf("Response Body: %s\n", string(responseBody))

	return nil
}

func getOrCreateAPIKey() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	localhookDir := filepath.Join(homeDir, ".localhook")
	apiKeyFile := filepath.Join(localhookDir, ".apikey")

	if _, err := os.Stat(apiKeyFile); err == nil {
		content, err := os.ReadFile(apiKeyFile)
		if err != nil {
			return "", fmt.Errorf("failed to read .apikey file: %v", err)
		}
		return string(content), nil
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(localhookDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create .localhook directory: %v", err)
		}

		randomString, err := generateRandomString(8)
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %v", err)
		}

		if err := os.WriteFile(apiKeyFile, []byte(randomString), 0600); err != nil {
			return "", fmt.Errorf("failed to write .apikey file: %v", err)
		}

		return randomString, nil
	} else {
		return "", fmt.Errorf("failed to check .apikey file: %v", err)
	}
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
