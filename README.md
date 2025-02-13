# 🚀 Workshop - Cache in App Memory

Learn and practice Cache in App memory with Golang through hands-on exercises.

## 📌  Overview

This workshop demonstrates how to implement **In-App Memory Caching** in **Golang** to optimize API performance and reduce database load.  
It covers different caching techniques using.

- **query sql**
- **sync.Map**
- **atomic.Value**

## 🚀 Getting Started

### **1️⃣ Clone the Repository**

```sh
git clone https://github.com/kaweel/workshop-cache-in-app-memory.git

cd kaweel_workshop_cache_in_app_memory
```

### 2️⃣ Set Up the Database (MSSQL via Docker)

```ssh
docker-compose up -d
```

### 3️⃣ Run the Application

```ssh
go run main.go
```

### 4️⃣ Run Benchmark Tests

```ssh
go test -bench=. -benchmem
```

### 🔥 Benchmark Results (Example)

| Method                  | Execution Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|-------------------------|----------------------:|--------------:|------------------------:|
| **MSSQL Query**         | 537,067 ns (~0.537 ms) | 3,451 B       | 42 allocs               |
| **sync.Map Cache**      | 33,790 ns (~0.0338 ms) | 3,459 B       | 42 allocs               |
| **atomic.Value Cache**  | 33,755 ns (~0.0337 ms) | 3,467 B       | 42 allocs               |

 ✅ Cache is ~15.9x faster than querying the database! 🚀

---

### 📜 License

This project is licensed under the MIT License.

### 🚀 Enjoy Coding! 😊
