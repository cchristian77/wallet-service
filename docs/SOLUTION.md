# Solution 

### Assumptions 
1. The service only manages **wallet-to-wallet disbursement** with `/transfers/v1` endpoint.
2. Transfer `amount` and `balance` is an **integer**. 
3. Each user has exactly **one wallet**, API identifies wallets by `wallets.id` for the sender (`from`) & recipient (`to`) wallets. <br>
  `Idempotency-Key` header is required per transfer attempts. If the same `Idempotency-Key` is sent the API will return `409`error.
4. Authentication and Authorization are not implemented for this API. 

### Locking Strategy
In a concurrent environment, multiple clients may attempt to transfer from the same wallet simultaneously,
which may lead to **race condition issue**, for example, some funds may be lost during the transfer due to inaccurate and unpredictable data. </br>

To prevent this, I **serialize the transfer processes using PostgreSQL row-level locks** inside a single database transaction
to ensure an atomic (all-or-nothing) operation.
When a client attempts a disbursement, the system does the following: </br>

1. Insert a `transactions` row with the `Idempotency-Key` as `transaction_id` with `PENDING` status to claim the request.
2. Lock sender (`from`) and recipient (`to`) wallets with `SELECT … FOR UPDATE` in **ascending wallet ID order** so opposite-direction transfers (A→B vs B→A)
   by making the order of the locks in a uniform manner across all pairs.
3. Validate that the source wallet has sufficient balance.
4. Write entries to `transaction_ledgers`: DEBIT for the sender and CREDIT for the recipient.
5. Update both wallet balances.
6. Mark the transaction as `SUCCESS` and commit. </br>

Because the lock is row-level, concurrent transfers that share any wallet wait until that lock is released -
same pair or overlapping pairs (e.g. `1→2` behind `1→3` on wallet `1`). Fully disjoint pairs (e.g. `1→2` and `3→4`) proceed in parallel.
On failure before commit, the transaction is rolled back so the transfer is all-or-nothing. </br>

This approach eliminates race conditions on balances by enforcing that only one transfer mutates a given wallet pair at a time/

### Consideration 
- For money transfer, **correct balances and consistency matter more than latency** - 
  a slightly slower transfer is better than lost funds or double-spending.
- **Strong consistency**: transfer operations happen under the same locks and DB transaction to ensure atomic operations (all-or-nothing).
- **Deadlock prevention**: locking wallets in a deterministic ascending ID order avoids circular wait & deadlocks between concurrent transfers. </br>
- **Idempotency**: unique `transactions.transaction_id` (from `Idempotency-Key`) prevents duplicate disbursements on client retries.
- **Transaction ledgers**: DEBIT/CREDIT pairs provide an audit trail that reconciles with wallet balances.
- **Trade-off**: under very high requests on the same wallet, row locks increase wait time. Therefore, latency is bounded by DB lock wait. </br>

### Stress Test
I made 2 stress tests to verify the efficiency and correctness of my solution: </br>

1. Test case 1 (`stress_test/test_case1/main.go`): sending 100 requests concurrently to transfer from wallet 1 (seed balance 10000) with amount 1000 to a random different wallet. </br>
   The result: **10 success & 90 failed requests** (due to insufficient balance after the sender is drained) with the average latency of ~250ms, without corrupted balances or unbalanced ledgers.

2. Test case 2 (`stress_test/test_case2/main.go`): sending 100 requests concurrently that alternate opposite directions (wallet 1→2 and wallet 2→1) to verify deadlock prevention. </br>
   The result: **100 success & 0 failed requests** and **0 timeouts** (no deadlock occurred) with the average latency of ~400ms. </br>

### Alternative
Optimistic Locking </br>
**Pros**
- Avoids long `FOR UPDATE` locks - better throughput when wallet contention is low.

**Cons :**
- Under high contention, version conflicts cause many retries. 

Message Queue</br>
**Pros** 
- Fast API response - workers can process transfers sequentially per wallet and handle traffic spikes.

**Cons :** 
- Extra infrastructure and idempotency handling - more difficult to return immediate balances/status since the transfers run asynchronously.

Distributed Lock (e.g. Redis) per wallet pair </br>
**Pros** 
- Serializes hot-wallet transfers before the DB, reducing Postgres lock wait under spikes.

**Cons :** 
- Extra dependency (TTL/outage risks); balances still require an atomic DB update as source of truth.

### Further Improvement
1. Persist `FAILED` status on business failures (e.g., insufficient balance) instead of only rolling back the `PENDING` claim, so the transfer outcome is durable.
2. Add `external_wallet_id` (UUID) for data exchange between external systems instead of relying on internal auto-increment IDs.
3. Add wallet statement / ledger history endpoints for reconciliation and support debugging.
4. Add metrics (lock wait time, transfer success rate) and tighter DB timeouts under contention.

### Author
Chris Christian 
