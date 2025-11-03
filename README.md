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
├── algo
│   ├── handlers
│   │   └── validateInput.go
│   ├── matchingAlgorithm.go
│   ├── orchestrator.go
│   └── submitOrderEntry.go
├── go.mod
├── go.sum
├── main.go
├── README.md
├── routes
│   ├── api
│   │   ├── get_orderbook.go
│   │   ├── order_cancel.go
│   │   ├── order.go
│   │   └── order_status.go
│   └── health.go
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


## Future Enhancements
- Metrics API
- WebSocket order/trade stream
- Persistent orderbook using Redis or PostgreSQL
- Cancel/Amend orders

