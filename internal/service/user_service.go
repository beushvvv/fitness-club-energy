package service

import (
	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"
	"time"
)

type UserService struct {
	userRepo  *repository.UserRepository
	cache     *cache.RedisClient
	cacheWrap *cache.CacheWrapper
}

func NewUserService(userRepo *repository.UserRepository, redisClient *cache.RedisClient) *UserService {
	return &UserService{
		userRepo:  userRepo,
		cache:     redisClient,
		cacheWrap: cache.NewCacheWrapper(redisClient),
	}
}

// GetAllUsers с кешированием
func (s *UserService) GetAllUsers() ([]model.User, error) {
	var users []model.User

	err := s.cacheWrap.GetOrSet(
		cache.UsersAllPrefix,
		5*time.Minute,
		&users,
		func() (interface{}, error) {
			return s.userRepo.FindAll()
		},
	)

	return users, err
}

// GetUserByID с кешированием
func (s *UserService) GetUserByID(id int) (*model.User, error) {
	var user model.User

	err := s.cacheWrap.GetOrSet(
		cache.UserKey(id),
		10*time.Minute,
		&user,
		func() (interface{}, error) {
			return s.userRepo.FindByID(id)
		},
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser создаёт пользователя и очищает кеш
func (s *UserService) CreateUser(user *model.User) error {
	// Сохраняем в БД
	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// Очищаем кеш списка пользователей
	for _, key := range cache.KeysToInvalidate("user") {
		_ = s.cache.Delete(key)
	}

	return nil
}

// UpdateUser обновляет пользователя и очищает кеш
func (s *UserService) UpdateUser(user *model.User) error {
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Очищаем кеш конкретного пользователя и списка
	_ = s.cache.Delete(cache.UserKey(user.ID))
	for _, key := range cache.KeysToInvalidate("user") {
		_ = s.cache.Delete(key)
	}

	return nil
}

// DeleteUser удаляет пользователя и очищает кеш
func (s *UserService) DeleteUser(id int) error {
	if err := s.userRepo.Delete(id); err != nil {
		return err
	}

	// Очищаем кеш
	_ = s.cache.Delete(cache.UserKey(id))
	for _, key := range cache.KeysToInvalidate("user") {
		_ = s.cache.Delete(key)
	}

	return nil
}
