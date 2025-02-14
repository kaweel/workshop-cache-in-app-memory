package main

import (
	"io"
	"net/http"
	"testing"
)

func BenchmarkAPIQueryFromDB(b *testing.B) {
	client := &http.Client{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := client.Get("http://localhost:8080/query/db?id=1")
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkAPIQueryFromCache(b *testing.B) {
	client := &http.Client{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := client.Get("http://localhost:8080/query/cache?id=500")
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkAPIQueryFromAtomicCache(b *testing.B) {
	client := &http.Client{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := client.Get("http://localhost:8080/query/atomic-cache?id=500")
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkAPIQueryFromRedisCache(b *testing.B) {
	client := &http.Client{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := client.Get("http://localhost:8080/query/redis-cache?id=500")
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
	}
}
