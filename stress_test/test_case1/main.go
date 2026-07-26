package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

/*
This program is a lightweight concurrent stress test for the transfer (disbursement) API.

It fires 100 parallel POST /transfers/v1 requests that always debit wallet 1 (seed balance 10000)
with amount 1000 toward a random different wallet.

Under correct locking, the expected result should be:
1. wallet 1 with balance of 10000, 10 request with transfer amount 1000 should be successful to transfer to other wallets.
2. the remaining 90 requests should fail with insufficient balance, without corrupted balances or unbalanced ledgers.
*/

// TransferRequest the request payload for POST /transfers/v1.
type TransferRequest struct {
	IdempotencyKey string `json:"-"`
	From           uint64 `json:"from"`
	To             uint64 `json:"to"`
	Amount         int64  `json:"amount"`
}

// TransferResponse holds the raw HTTP status and body returned by the API.
type TransferResponse struct {
	StatusCode int
	Body       string
}

func main() {
	var (
		successCount    int64
		failureCount    int64
		totalDurationNs int64
		totalRequests   = 100
	)

	var wg sync.WaitGroup
	wg.Add(totalRequests)

	// Fire concurrent transfer requests.
	for i := 0; i < totalRequests; i++ {
		go func() {
			defer wg.Done()

			startTime := time.Now()

			// Always debit wallet 1 (initial balance 10000); credit a random different wallet (2 or 3).
			from := uint64(1)
			to := uint64(rand.Intn(3) + 1)
			for to == from {
				to = uint64(rand.Intn(3) + 1)
			}

			req := &TransferRequest{
				IdempotencyKey: fmt.Sprintf("TRX-STRESS-%d-%d", time.Now().UnixNano(), i),
				From:           from,
				To:             to,
				Amount:         1000,
			}

			fmt.Printf("[i=%d] Requesting transfer from wallet %d to wallet %d with amount %d \n", i, req.From, req.To, req.Amount)

			resp, err := doTransfer(req)
			duration := time.Since(startTime)

			if err != nil {
				atomic.AddInt64(&failureCount, 1)
				fmt.Printf("[KEY=%s] FAILED | err=%v\n", req.IdempotencyKey, err)
				fmt.Printf("[KEY=%s] Request duration : %s\n", req.IdempotencyKey, duration.String())
				atomic.AddInt64(&totalDurationNs, duration.Nanoseconds())
				return
			}

			if resp.StatusCode >= 400 {
				atomic.AddInt64(&failureCount, 1)
				fmt.Printf("[KEY=%s] FAILED | Status %d | Response : %s \n", req.IdempotencyKey, resp.StatusCode, resp.Body)
			} else if resp.StatusCode >= 200 {
				atomic.AddInt64(&successCount, 1)
				fmt.Printf("[KEY=%s] SUCCESS | Status %d | Response : %s\n", req.IdempotencyKey, resp.StatusCode, resp.Body)
			}

			fmt.Printf("[KEY=%s] Request duration : %s\n", req.IdempotencyKey, duration.String())
			atomic.AddInt64(&totalDurationNs, duration.Nanoseconds())
		}()
	}

	wg.Wait()

	// Aggregate results after all goroutines finish.
	fmt.Println("====== Results ======")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failureCount)
	fmt.Printf("Average Response Time: %s\n", time.Duration(totalDurationNs/int64(totalRequests)).String())
}

// doTransfer sends a single POST /transfers/v1 request to the local server.
func doTransfer(body *TransferRequest) (*TransferResponse, error) {
	url := "http://localhost:9000"

	client := &http.Client{Timeout: 5 * time.Second}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("JSON marshal error: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/transfers/v1", url)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", body.IdempotencyKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &TransferResponse{
		StatusCode: resp.StatusCode,
		Body:       string(respBodyBytes),
	}, nil
}
