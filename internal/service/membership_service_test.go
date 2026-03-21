package service

import (
	"fmt"
	"testing"
	"time"

	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/logger"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.Init("debug")
}

func setupTestDBMembership(t *testing.T) (*sqlx.DB, func()) {
	db, err := sqlx.Connect("postgres", "host=localhost port=5433 user=postgres password=postgres123 dbname=fitness_club sslmode=disable")
	require.NoError(t, err)

	// Очищаем таблицы
	_, err = db.Exec("DELETE FROM memberships")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM users")
	require.NoError(t, err)

	cleanup := func() {
		db.Exec("DELETE FROM memberships")
		db.Exec("DELETE FROM users")
		db.Close()
	}
	return db, cleanup
}

func setupTestCacheMembership(t *testing.T) (*cache.RedisClient, func()) {
	client := cache.NewRedisClient("localhost:6379", "", 0)
	err := client.Ping()
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
	client.Delete("memberships:all")
	client.Delete("membership:1")

	cleanup := func() {
		client.Delete("memberships:all")
		client.Delete("membership:1")
		client.Close()
	}
	return client, cleanup
}

func createTestUserForMembership(t *testing.T, db *sqlx.DB, email string) int {
	var userID int
	err := db.QueryRow(
		"INSERT INTO users (name, email, phone, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id",
		"Test User", email, "123456789",
	).Scan(&userID)
	require.NoError(t, err)
	return userID
}

func TestMembershipService_GetAllMemberships_Integration(t *testing.T) {
	db, dbCleanup := setupTestDBMembership(t)
	defer dbCleanup()

	redisClient, redisCleanup := setupTestCacheMembership(t)
	defer redisCleanup()

	cacheWrapper := cache.NewCacheWrapper(redisClient)
	membershipRepo := repository.NewMembershipRepository(db)
	service := NewMembershipService(membershipRepo, cacheWrapper)

	// Создаём уникального пользователя
	userID := createTestUserForMembership(t, db, fmt.Sprintf("test_membership_%d@test.com", time.Now().UnixNano()))

	// Создаём тестовый абонемент с допустимым типом
	testMembership := &model.Membership{
		UserID:    userID,
		Type:      "standard", // допустимый тип
		Price:     1000,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 1, 0),
		IsActive:  true,
	}
	err := membershipRepo.Create(testMembership)
	require.NoError(t, err)

	// Первый запрос — из БД
	memberships1, err := service.GetAllMemberships()
	assert.NoError(t, err)
	assert.Len(t, memberships1, 1)

	// Второй запрос — из кэша
	memberships2, err := service.GetAllMemberships()
	assert.NoError(t, err)
	assert.Len(t, memberships2, 1)

	assert.Equal(t, memberships1[0].Type, memberships2[0].Type)
}

func TestMembershipService_CreateMembership_Integration(t *testing.T) {
	db, dbCleanup := setupTestDBMembership(t)
	defer dbCleanup()

	redisClient, redisCleanup := setupTestCacheMembership(t)
	defer redisCleanup()

	cacheWrapper := cache.NewCacheWrapper(redisClient)
	membershipRepo := repository.NewMembershipRepository(db)
	service := NewMembershipService(membershipRepo, cacheWrapper)

	// Создаём уникального пользователя
	userID := createTestUserForMembership(t, db, fmt.Sprintf("test_create_%d@test.com", time.Now().UnixNano()))

	// Используем допустимый тип
	membership := &model.Membership{
		UserID:    userID,
		Type:      "premium", // допустимый тип
		Price:     5000,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(1, 0, 0),
		IsActive:  true,
	}

	err := service.CreateMembership(membership)
	assert.NoError(t, err)
	assert.NotZero(t, membership.ID)

	// Проверяем что создалось
	created, err := membershipRepo.FindByID(membership.ID)
	assert.NoError(t, err)
	assert.Equal(t, membership.Type, created.Type)
}
