package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/stretchr/testify/require"
)

type APIClient struct {
	baseURL string
	client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       100,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
		},
	}
}

func (c *APIClient) CreateAccount(t *testing.T, account models.AccountCreate) (string, error) {
	resp, err := c.client.Post(
		c.baseURL+"/account",
		"application/json",
		jsonBody(account),
	)
	if err != nil {
		return "", fmt.Errorf("account creation failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.ID, nil
}
func (c *APIClient) CreateTransaction(t *testing.T, tx models.Transaction) (string, error) {
	resp, err := c.client.Post(
		c.baseURL+"/transaction",
		"application/json",
		jsonBody(tx),
	)
	fmt.Println(tx)
	fmt.Println(c.baseURL + "/transaction")
	if err != nil {
		return "", fmt.Errorf("transaction request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode transaction response: %w", err)
	}

	return result.ID, nil
}

func TestHighConcurrencyWorkflow(t *testing.T) {

	client := NewAPIClient("http://localhost:8080")

	const numAccounts = 100
	accounts := createAccounts(t, client, numAccounts)

	fmt.Println(accounts)

	const transactionsPerAccount = 1
	performTransactions(t, client, accounts, transactionsPerAccount)

}

func createAccounts(t *testing.T, client *APIClient, count int) []string {
	var wg sync.WaitGroup
	accounts := make([]string, 0, count)
	errChan := make(chan error, count)
	accountChan := make(chan string, count)
	sem := make(chan struct{}, 2)

	var mu sync.Mutex

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			accountID, err := client.CreateAccount(t, models.AccountCreate{
				FirstName: fmt.Sprintf("LoadUser%d", id),
				LastName:  "LoadTest",
				Email:     fmt.Sprintf("loaduser%d-%d@test.com", time.Now().UnixNano(), id),
				Phone:     fmt.Sprintf("+1555LOAD%06d-%d", id, time.Now().UnixNano()),
			})
			if err != nil {
				errChan <- fmt.Errorf("account %d: %v", id, err)
				return
			}

			mu.Lock()
			accounts = append(accounts, accountID)
			mu.Unlock()
			accountChan <- accountID
		}(i)
	}

	wg.Wait()
	close(accountChan)
	close(errChan)

	require.Empty(t, readErrors(errChan), "Account creation errors")

	return accounts
}

func performTransactions(t *testing.T, client *APIClient, accounts []string, txPerAccount int) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(accounts)*txPerAccount)
	sem := make(chan struct{}, 2)

	var mu sync.Mutex

	for _, accountID := range accounts {
		for i := 0; i < txPerAccount; i++ {
			wg.Add(1)
			go func(accID string, txNum int) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				txType := models.DEPOSIT
				amount := 100.0

				_, err := client.CreateTransaction(t, models.Transaction{
					AccountID: accID,
					Type:      txType,
					Amount:    amount,
				})
				if err != nil {
					mu.Lock()
					errChan <- fmt.Errorf("account %s tx %d: %v", accID, txNum, err)
					mu.Unlock()
				}
			}(accountID, i)
		}
	}

	wg.Wait()
	close(errChan)

	require.Empty(t, readErrors(errChan), "Transaction errors")
}

func readErrors(errChan <-chan error) []error {
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}
	return errors
}

func jsonBody(v interface{}) io.Reader {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(v)
	return buf
}
