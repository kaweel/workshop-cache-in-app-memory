package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"sync/atomic"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/redis/go-redis/v9"
)

var db *sql.DB
var inMemoryCache = &sync.Map{}
var atomicCache atomic.Value
var rdb *redis.Client
var ctx context.Context

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func init() {
	ctx = context.Background()
	var err error
	db, err = sql.Open("sqlserver", "server=localhost;user id=sa;password=SuperStrong@Passw0rd;database=cache")
	if err != nil {
		log.Fatal(err)
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // เปลี่ยนเป็น Redis address ของคุณ
		Password: "",               // ถ้ามี password ให้ใส่ที่นี่
		DB:       0,                // ใช้ database หมายเลข 0
	})

	loadFromSyncCache()
	loadFromAtomicCache()
	loadFromRedis(rdb)

}

func loadFromSyncCache() {
	fmt.Println("Loading 1000 records into In-Memory Cache...")
	for i := 1; i <= 1000; i++ {
		inMemoryCache.Store(strconv.Itoa(i), &Product{ID: strconv.Itoa(i), Name: fmt.Sprintf("Product %d", i)})
	}
	fmt.Println("✅ In-Memory Cache Ready!")
}

func loadFromAtomicCache() {
	fmt.Println("Loading 1000 records into In-Memory Atomic Cache...")
	cache := make(map[string]*Product)
	for i := 1; i <= 1000; i++ {
		cache[strconv.Itoa(i)] = &Product{ID: strconv.Itoa(i), Name: fmt.Sprintf("Product %d", i)}
	}
	atomicCache.Store(cache)
	fmt.Println("✅ In-Memory Atomic Cache Ready!")
}

func loadFromRedis(rdb *redis.Client) {
	fmt.Println("Loading 1000 records into In-Memory Redis Cache...")
	for i := 1; i <= 1000; i++ {
		product, _ := json.Marshal(Product{ID: strconv.Itoa(i), Name: fmt.Sprintf("Product %d", i)})
		rdb.Set(ctx, strconv.Itoa(i), product, 0)
	}
	fmt.Println("✅ In-Memory Redis Cache Ready!")
}

func queryFromDBHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	product := &Product{}
	err := db.QueryRow("SELECT id,name FROM products WHERE id = @p1", productID).Scan(&product.ID, &product.Name)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func queryFromCacheHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	val, _ := inMemoryCache.Load(productID)
	json.NewEncoder(w).Encode(val)
}

func queryFromAtomicCache(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	cache := atomicCache.Load().(map[string]*Product)
	product := cache[productID]
	json.NewEncoder(w).Encode(*product)
}

func queryFromRedisCache(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	product, _ := rdb.Get(ctx, productID).Result()
	json.NewEncoder(w).Encode(product)
}

func main() {
	http.HandleFunc("/query/db", queryFromDBHandler)
	http.HandleFunc("/query/cache", queryFromCacheHandler)
	http.HandleFunc("/query/atomic-cache", queryFromAtomicCache)
	http.HandleFunc("/query/redis-cache", queryFromRedisCache)

	fmt.Println("✅ Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
