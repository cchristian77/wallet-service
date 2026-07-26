package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
This stress test verifies deadlock prevention for opposite-direction transfers.

It fires N concurrent POST /transfers/v1 requests that alternate:
  - even i → wallet 1 → wallet 2 (A→B)
  - odd  i → wallet 2 → wallet 1 (B→A)

With correct locking, requests should complete success without timeout occurring.
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
		timeoutCount    int64
		totalDurationNs int64
		totalRequests   = 100
	)

	var (
		walletA = uint64(1)
		walletB = uint64(2)
		amount  = int64(100)
	)

	var wg sync.WaitGroup
	wg.Add(totalRequests)

	fmt.Println("====== Deadlock Prevention Stress Test ======")

	for i := 0; i < totalRequests; i++ {
		i := i
		go func() {
			defer wg.Done()

			startTime := time.Now()

			// Alternate directions to create A→B and B→A lock contention.
			from, to := walletA, walletB
			direction := "A -> B"
			if i%2 == 1 {
				from, to = walletB, walletA
				direction = "B -> A"
			}

			req := &TransferRequest{
				IdempotencyKey: fmt.Sprintf("TRX-DEADLOCK-%d-%d", time.Now().UnixNano(), i),
				From:           from,
				To:             to,
				Amount:         amount,
			}

			fmt.Printf("[i=%d][%s] Requesting transfer from wallet %d to wallet %d with amount %d\n",
				i, direction, req.From, req.To, req.Amount)

			resp, err := doTransfer(req)
			duration := time.Since(startTime)
			atomic.AddInt64(&totalDurationNs, duration.Nanoseconds())

			if err != nil {
				atomic.AddInt64(&failureCount, 1)
				if isTimeout(err) {
					atomic.AddInt64(&timeoutCount, 1)
				}
				fmt.Printf("[KEY=%s][%s] FAILED | err=%v\n", req.IdempotencyKey, direction, err)
				fmt.Printf("[KEY=%s] Request duration : %s\n", req.IdempotencyKey, duration.String())
				return
			}

			if resp.StatusCode >= 400 {
				atomic.AddInt64(&failureCount, 1)
				fmt.Printf("[KEY=%s][%s] FAILED | Status %d | Response : %s\n",
					req.IdempotencyKey, direction, resp.StatusCode, resp.Body)
			} else if resp.StatusCode >= 200 {
				atomic.AddInt64(&successCount, 1)
				fmt.Printf("[KEY=%s][%s] SUCCESS | Status %d | Response : %s\n",
					req.IdempotencyKey, direction, resp.StatusCode, resp.Body)
			}

			fmt.Printf("[KEY=%s] Request duration : %s\n", req.IdempotencyKey, duration.String())
		}()
	}

	wg.Wait()

	fmt.Println("====== Results ======")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failureCount)
	fmt.Printf("Timeouts: %d\n", timeoutCount)
	fmt.Printf("Average Response Time: %s\n", time.Duration(totalDurationNs/int64(totalRequests)).String())

	if timeoutCount == 0 {
		fmt.Println("Deadlock check: PASS. No requests timed out.")
	} else {
		fmt.Println("Deadlock check: FAIL. One or more requests timed out.")
	}
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	type timeout interface{ Timeout() bool }
	if t, ok := err.(timeout); ok && t.Timeout() {
		return true
	}
	// net/http wraps timeouts; match common message.
	msg := err.Error()
	return strings.Contains(msg, "Client.Timeout") || strings.Contains(msg, "deadline exceeded")
}

// doTransfer sends a single POST /transfers/v1 request to the local server.
func doTransfer(body *TransferRequest) (*TransferResponse, error) {
	url := "http://localhost:9000"

	// Longer timeout than the default stress test: lock waits under A↔B contention are expected,
	// but a true deadlock would hang until this deadline.
	client := &http.Client{Timeout: 30 * time.Second}

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
