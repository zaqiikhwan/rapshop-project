package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type CoreAPI struct {
	ca coreapi.Client
}

func NewMidtransDriver() CoreAPI {
	return CoreAPI{ca: coreapi.Client{}}
}

func (c *CoreAPI) HandleNotification(id string) (*coreapi.TransactionStatusResponse, error) {
	env := midtrans.Sandbox
	if strings.EqualFold(os.Getenv("MIDTRANS_ENV"), "production") {
		env = midtrans.Production
	}
	c.ca.New(os.Getenv("AUTHORIZATION_VALUE"), env)

	midtransReport, err := c.ca.CheckTransaction(id)
	if err != nil {
		return midtransReport, err
	}

	return midtransReport, nil
}

func basicAuth() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("AUTHORIZATION_VALUE")))
}

// Charge posts a charge payload to Midtrans. It returns the parsed response, the
// Midtrans status_code (as int), and any transport error.
func (c *CoreAPI) Charge(payload map[string]any) (map[string]any, int, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequest(http.MethodPost, os.Getenv("MIDTRANS"), strings.NewReader(string(data)))
	if err != nil {
		return nil, 0, err
	}
	return doMidtrans(req)
}

// CheckStatus fetches a transaction's live status from Midtrans.
func (c *CoreAPI) CheckStatus(id string) (map[string]any, error) {
	base := os.Getenv("MIDTRANS_STATUS_URL")
	if base == "" {
		base = "https://api.midtrans.com/v2"
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/status", base, id), nil)
	if err != nil {
		return nil, err
	}
	body, _, err := doMidtrans(req)
	return body, err
}

func doMidtrans(req *http.Request) (map[string]any, int, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", basicAuth())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, 0, err
	}

	status := 0
	if sc, ok := parsed["status_code"].(string); ok {
		status, _ = strconv.Atoi(sc)
	}
	return parsed, status, nil
}
