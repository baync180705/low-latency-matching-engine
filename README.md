# Low-Latency Matching Engine

A lightweight, high-performance matching engine built in Go with REST APIs using Echo.  
Supports basic limit and market order matching with an in-memory order book.

---

## Features

- REST API built with **Echo**
- **Limit** and **Market** order support
- In-memory order book using **heap-based priority queues**
- Real-time trade execution
- Clean modular design (API, Algo, Types packages)
- Ready for WebSocket integration

---

## Project Structure

```
├── engine
│   ├── handlers
│   │   └── validateInput.go
│   ├── matchingAlgorithm.go
│   ├── orchestrator.go
│   └── submitOrderEntry.go
├── go.mod
├── go.sum
├── main.go
├── Matching_Engine_Documentation.md
├── metrics
│   └── metrics.go
├── README.md
├── routes
│   ├── api
│   │   ├── get_orderbook.go
│   │   ├── order_cancel.go
│   │   ├── order.go
│   │   └── order_status.go
│   ├── health.go
│   └── metrics.go
└── types
    ├── heap.go
    ├── order.go
    ├── registry.go
    └── response.go
```

---

## Setup

### 1. Clone the Repository
```bash
git clone https://github.com/baync180705/low-latency-matching-engine.git
cd low-latency-matching-engine
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Run the Server
```bash
go run main.go
```

Server runs by default on **http://localhost:8080**

---

## Load Testing

Load Testing using wrk is supported

### Test the Server

Make the shell script executable
```bash
chmod +x run_bench.sh
```
Run the script
```bash
./run_bench.sh
```

## Benchmarks Obtained 
Below are the benchmark results obtained from one of the test runs:

| Metric                      | Value              |
|------------------------------|--------------------:|
| **Orders Received**          | 391,896            |
| **Orders Matched**           | 503,497            |
| **Orders Cancelled**         | 74,380             |
| **Orders In Book**           | 75                 |
| **Trades Executed**          | 316,321            |
| **Latency P50 (ms)**         | 0.003              |
| **Latency P99 (ms)**         | 0.073              |
| **Latency P99.9 (ms)**       | 0.595315           |
| **Throughput (orders/sec)**  | 22,870.64          |


## Future Enhancements
- Metrics API
- WebSocket order/trade stream
- Persistent orderbook using Redis or PostgreSQL
- Amend orders

