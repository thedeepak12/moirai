# Moirai

A concurrent, context-aware worker pool in Go that supports dynamic scaling, polymorphic tasks, and real-time telemetry.

## Features

- **Polymorphic Tasks**: Run any task structure by implementing a simple Go interface.
- **Context Integration**: Cancel tasks and handle timeouts cleanly across all running workers.
- **Dynamic Scaling**: Add or remove active worker goroutines on-the-fly at runtime.
- **Real-Time Telemetry**: Track progress and performance metrics safely using atomic counters.
- **Graceful Shutdown**: Drain the pending job queue completely before terminating workers.

## Tech Stack

- **Language**: Go (1.26.4)
- **Standard Library Primitives**: Goroutines, Channels, `sync.WaitGroup`, `sync.Mutex`, `sync/atomic`, `context.Context`, `time.Ticker`.
- **Dependencies**: None. Pure Go Standard Library.

## Project Structure

```text
moirai/
├── cmd/
│   └── demo/
│       └── main.go       # Demo runner & task definitions
├── internal/
│   └── pool/
│       ├── pool.go       # Pool orchestration (ScaleUp, ScaleDown, Wait)
│       ├── worker.go     # Worker execution loops and signal handling
│       ├── task.go       # Generic Task interface, Job, and Result types
│       └── metrics.go    # Thread-safe atomic counters and progress metrics
└── go.mod                # Go module descriptor
```

## Setup

1. Clone the repository:
```bash
git clone https://github.com/thedeepak12/moirai.git
cd moirai
```

2. Tidy the Go modules:
```bash
go mod tidy
```

3. Run the application:
```bash
go run ./cmd/demo/main.go
```
## License
Distributed under the MIT License. See [LICENSE](https://github.com/thedeepak12/moirai/blob/main/LICENSE) for more information.
