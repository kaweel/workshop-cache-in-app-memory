package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"sync/atomic"

	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB
var inMemoryCache = &sync.Map{}
var atomicCache atomic.Value

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	var err error
	db, err = sql.Open("sqlserver", "server=localhost;user id=sa;password=SuperStrong@Passw0rd;database=cache")
	if err != nil {
		log.Fatal(err)
	}
	loadFromSyncCache()
	loadFromAtomicCache()
}

func loadFromSyncCache() {
	fmt.Println("Loading 1000 records into In-Memory Cache...")
	for i := 1; i <= 1000; i++ {
		inMemoryCache.Store(i, &Product{ID: i, Name: fmt.Sprintf("Product %d", i)})
	}
	fmt.Println("✅ In-Memory Cache Ready!")
}

func loadFromAtomicCache() {
	fmt.Println("Loading 1000 records into In-Memory Atomic Cache...")
	cache := make(map[int]*Product)
	for i := 1; i <= 1000; i++ {
		cache[i] = &Product{ID: i, Name: fmt.Sprintf("Product %d", i)}
	}
	atomicCache.Store(cache)
	fmt.Println("✅ In-Memory Atomic Cache Ready!")
}

func queryFromDBHandler(w http.ResponseWriter, r *http.Request) {
	productID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	product := &Product{}
	err := db.QueryRow("SELECT id,name FROM products WHERE id = @p1", productID).Scan(&product.ID, &product.Name)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func queryFromCacheHandler(w http.ResponseWriter, r *http.Request) {
	productID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	val, _ := inMemoryCache.Load(productID)
	json.NewEncoder(w).Encode(val)
}

func queryFromAtomicCache(w http.ResponseWriter, r *http.Request) {
	productID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	cache := atomicCache.Load().(map[int]*Product)
	product := cache[productID]
	json.NewEncoder(w).Encode(*product)
}

func main() {
	http.HandleFunc("/query/db", queryFromDBHandler)
	http.HandleFunc("/query/cache", queryFromCacheHandler)
	http.HandleFunc("/query/atomic-cache", queryFromAtomicCache)

	fmt.Println("✅ Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
