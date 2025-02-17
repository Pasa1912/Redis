# GoRedis 🚀

A Redis clone implemented in Go with support for basic Redis commands, RESP (Redis Serialization Protocol) protocol, and AOF (Append Only File) persistence. This implementation focuses on simplicity and reliability while maintaining compatibility with standard Redis clients.

<img src="https://retool.com/_next/image?url=https%3A%2F%2Fcdn.sanity.io%2Fimages%2Fbclf52sw%2Fproduction%2F76750a5eba5df22b3abe63904b770da4a4473557-2048x800.webp&w=3840&q=75" height="400">

## ✨ Features

- **RESP Protocol Implementation**: Full support for Redis wire protocol
- **Command Support**:
  - String Operations: `GET`, `SET`
  - Hash Operations: `HGET`, `HSET`, `HGETALL`
  - Server Operations: `PING`
- **Persistence**: Append-Only File (AOF) based persistence with automatic background syncing
- **Concurrent Access**: Thread-safe operations using mutex locks
- **Standard Compatibility**: Works with official Redis CLI and clients

## 🚀 Quick Start

### Prerequisites

- Go 1.19 or higher
- Redis CLI (for testing)

### Running the Server

1. Clone the repository:

```bash
git clone https://github.com/Pasa1912/Redis
cd Redis
```

2. Start the server:

```bash
go run *.go
```

The server will start on port 6379 (default Redis port).

### Connecting with Redis CLI

```bash
redis-cli
```

## 💡 Usage Examples

### Basic Operations

```bash
# Ping the server
redis-cli> PING
"PONG"

# Set a key
redis-cli> SET mykey "Hello, GoRedis!"
"OK"

# Get a key
redis-cli> GET mykey
"Hello, GoRedis!"
```

### Hash Operations

```bash
# Set hash fields
redis-cli> HSET user:1 name "John Doe"
"OK"
redis-cli> HSET user:1 email "john@example.com"
"OK"

# Get a specific hash field
redis-cli> HGET user:1 name
"John Doe"

# Get all hash fields
redis-cli> HGETALL user:1
1) "name"
2) "John Doe"
3) "email"
4) "john@example.com"
```

## 🔧 Technical Implementation

### Architecture

The project is organized into several key components:

```
Redis/
├── main.go          # Server initialization and connection handling
├── handler.go       # Command implementations
├── aof.go          # Persistence layer
├── serializer.go    # RESP protocol serialization
└── deserializer.go  # RESP protocol deserialization
```

### RESP Protocol

The RESP protocol implementation supports five data types:

- Simple Strings ("+")
- Errors ("-")
- Integers (":")
- Bulk Strings ("$")
- Arrays ("\*")

### Persistence Layer

The AOF persistence implementation features:

- Automatic background syncing every second
- Atomic writes using mutex locks
- Crash recovery through AOF replay
- Append-only format for data integrity

### Concurrency Handling

Thread safety is ensured through:

- Read-Write mutexes for data stores
- Separate mutex for AOF operations
- Connection limiting (5 concurrent connections)

## 🛠️ Code Structure

Key components and their responsibilities:

- **main.go**:

  - Server initialization
  - Connection handling
  - Command routing

- **handler.go**:

  - Command implementations
  - Data store management
  - Thread-safe operations

- **aof.go**:

  - Persistence management
  - Background syncing
  - Recovery operations

- **serializer.go & deserializer.go**:
  - RESP protocol implementation
  - Data type handling
  - Wire format processing

## 🤝 Contributing

Contributions are welcome! Some areas for potential improvement:

- Additional Redis commands
- Cluster support
- RDB persistence format
- Command pipelining
- Transaction support

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⭐️ Show your support

Give a ⭐️ if this project helped you!

---

Made with ❤️ using Go
