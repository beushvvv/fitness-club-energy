package service

import (
	"testing"

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
	// Инициализируем логгер для тестов
	logger.Init("debug")
}

// setupTestDB подключается к реальной PostgreSQL в Docker
func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	db, err := sqlx.Connect("postgres", "host=localhost port=5433 user=postgres password=postgres123 dbname=fitness_club sslmode=disable")
	require.NoError(t, err)

	// Очищаем таблицу перед тестами
	_, err = db.Exec("DELETE FROM users")
	require.NoError(t, err)

	cleanup := func() {
		db.Exec("DELETE FROM users")
		db.Close()
	}
	return db, cleanup
}

// setupTestCache подключается к реальному Redis в Docker
func setupTestCache(t *testing.T) (*cache.RedisClient, func()) {
	client := cache.NewRedisClient("localhost:6379", "", 0)
	err := client.Ping()
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
	// Очищаем кэш
	client.Delete("users:all")
	client.Delete("user:1")
	client.Delete("user:2")

	cleanup := func() {
		client.Delete("users:all")
		client.Delete("user:1")
		client.Delete("user:2")
		client.Close()
	}
	return client, cleanup
}

func TestUserService_GetAllUsers_Integration(t *testing.T) {
	db, dbCleanup := setupTestDB(t)
	defer dbCleanup()

	redisClient, redisCleanup := setupTestCache(t)
	defer redisCleanup()

	cacheWrapper := cache.NewCacheWrapper(redisClient)
	userRepo := repository.NewUserRepository(db)
	service := NewUserService(userRepo, cacheWrapper)

	// Создаём тестового пользователя
	testUser := &model.User{Name: "Test", Email: "test@test.com", Phone: "123"}
	err := userRepo.Create(testUser)
	require.NoError(t, err)

	// Первый запрос — данные из БД
	users1, err := service.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users1, 1)

	// Второй запрос — данные из кэша
	users2, err := service.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users2, 1)

	assert.Equal(t, users1[0].Name, users2[0].Name)
}

func TestUserService_GetUserByID_Integration(t *testing.T) {
	db, dbCleanup := setupTestDB(t)
	defer dbCleanup()

	redisClient, redisCleanup := setupTestCache(t)
	defer redisCleanup()

	cacheWrapper := cache.NewCacheWrapper(redisClient)
	userRepo := repository.NewUserRepository(db)
	service := NewUserService(userRepo, cacheWrapper)

	// Создаём пользователя
	user := &model.User{Name: "Test User", Email: "testid@test.com", Phone: "789"}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Получаем по ID
	found, err := service.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, found.Name)
	assert.Equal(t, user.Email, found.Email)

	// Второй раз — из кэша
	found2, err := service.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, found2.Name)
}

func TestUserService_CreateUser_Integration(t *testing.T) {
	db, dbCleanup := setupTestDB(t)
	defer dbCleanup()

	redisClient, redisCleanup := setupTestCache(t)
	defer redisCleanup()

	cacheWrapper := cache.NewCacheWrapper(redisClient)
	userRepo := repository.NewUserRepository(db)
	service := NewUserService(userRepo, cacheWrapper)

	user := &model.User{Name: "New User", Email: "new@test.com", Phone: "456"}

	err := service.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Проверяем что пользователь создался
	created, err := userRepo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, created.Name)

	// Проверяем что кэш очищен — следующий GetAllUsers должен пойти в БД
	users, err := service.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	db, dbCleanup := setupTestDB(t)
	defer dbCleanup()

	redisClient, redisCleanup := setupTestCache(t)
	defer redisCleanup()

	cacheWrapper := cache.NewCacheWrapper(redisClient)
	userRepo := repository.NewUserRepository(db)
	service := NewUserService(userRepo, cacheWrapper)

	// Пытаемся получить несуществующего пользователя
	user, err := service.GetUserByID(99999)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "no rows")
}
