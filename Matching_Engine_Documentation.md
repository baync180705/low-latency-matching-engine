# Technical Documentation: Matching Engine Algorithm

## Overview

This matching engine is designed to execute trades with **low latency** by maintaining an in-memory order book structured using **custom heap-based priority queues**.  
It efficiently matches buy and sell orders based on price and time priority, following the standard **price-time priority matching principle** used in financial exchanges.

---

## Core Architecture

### 1. Order Book Structure
Each symbol (e.g., BTCUSD) maintains an independent order book consisting of:
- **Buy Heap** — A max-heap (highest price on top)
- **Sell Heap** — A min-heap (lowest price on top)
- **Price Map (TimeQueue)** — Maps each price level to a FIFO list of orders (`OrderList`), preserving time priority

This combination enables **O(log N)** insertion and **O(1)** best-price retrieval.

---

### 2. Matching Algorithm

#### Step 1: Input Validation
The system validates:
- Symbol, side, and type fields
- Price and quantity (for limit orders)
- Market order constraints (no price field required)

Invalid inputs immediately return a `400 Bad Request`.

---

#### Step 2: Price Discovery and Matching
When a new order arrives:
1. Identify the **counter side** (Buy → match with Sell, Sell → match with Buy)
2. Check if matching prices exist:
   - For **Buy orders**: Match if `BuyPrice >= BestSellPrice`
   - For **Sell orders**: Match if `SellPrice <= BestBuyPrice`
3. If both heaps are non-empty and price overlap exists, trades are executed.

Otherwise, the order is placed in the order book.

---

#### Step 3: Trade Execution
Matching is performed iteratively:
1. Fetch the best price from the counter-side heap top.
2. Retrieve the earliest order (FIFO) at that price from `TimeQueue`.
3. Calculate:
   ```
   matchedQty = min(incomingOrder.Quantity, restingOrder.Quantity)
   ```
4. Record a trade entry with:
   - Price = resting order’s price
   - Quantity = matchedQty
   - Timestamps and order IDs

5. Update quantities:
   - Reduce both orders’ remaining quantities
   - Remove from book when fully filled
6. Repeat until:
   - Incoming order is fully filled, or
   - No more eligible price levels exist

---

### 3. Heap Maintenance

Custom heap implementation ensures:
- **BuyHeap** (max-heap): Highest bid price remains at index 0  
  `Less(i, j)` returns `true` if `Price[i] > Price[j]`
- **SellHeap** (min-heap): Lowest ask price remains at index 0  
  `Less(i, j)` returns `true` if `Price[i] < Price[j]`

This enables constant-time access to the top of the order book and logarithmic insertion/deletion.

---

## Latency Optimization Techniques

### 1. In-Memory Data Structures
All operations (insert, match, remove) are executed entirely in memory — no external I/O or database calls occur in the matching path.

**Effect:** Sub-millisecond response times for order matching in typical loads.

---

### 2. Heap-Based Price Levels
The use of Go’s heap interface with custom logic allows:
- Constant-time retrieval of best bid/ask
- Logarithmic updates (`heap.Push`, `heap.Pop`)
- Minimal lock contention when scaled with concurrent symbols

**Effect:** Reduces matching complexity from O(N) to O(log N).

---

### 3. FIFO Queues for Time Priority
Each price level maintains a linked list (`OrderList`), ensuring earliest orders are filled first.  
This avoids sorting overhead and guarantees deterministic execution.

---

### 4. Optimized Error and Response Handling
Errors (like “No match found”) are short-circuited early, avoiding unnecessary computations.  
Responses are directly serialized using Echo’s JSON encoder for low-overhead API delivery.

---

### 5. Stateless Matching Pipeline
The `RunPipeline` function is designed as a pure computation pipeline:
1. Accepts an order input
2. Mutates in-memory order book
3. Returns trades, updated order state, and errors

No shared global state is modified outside of synchronized order book maps, ensuring consistent and parallel execution across multiple instruments.

---


## Why It’s Low Latency

- Fully **in-memory** — no database or disk interaction
- **Heap + FIFO hybrid** for O(log N) price management and O(1) matching
- **Stateless REST handlers** — direct JSON encoding/decoding

---

## Summary

This matching engine achieves low latency through:
1. **Efficient in-memory design**
2. **Custom heap-based priority queues**
3. **Price-time matching logic**
4. **Lightweight REST layer**

It provides a robust and fast foundation for high-throughput trading systems or simulation environments.
