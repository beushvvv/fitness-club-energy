package cache

import (
	"testing"
	"time"
)

func TestRedisSetGet(t *testing.T) {
	// Подключаемся
	client := NewRedisClient("localhost:6379", "", 0)

	// Тестовые данные
	testKey := "test:key"
	testValue := "Hello Redis"

	// Set
	err := client.Set(testKey, testValue, 1*time.Minute)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Get
	var result string
	err = client.Get(testKey, &result)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}

	if result != testValue {
		t.Errorf("Expected %s, got %s", testValue, result)
	}

	// Cleanup
	client.Delete(testKey)
	t.Log("✅ Redis test passed")
}
