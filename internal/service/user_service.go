package service

import (
	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"
	"log"
	"time"
)

type UserService struct {
	userRepo     *repository.UserRepository
	cacheWrapper *cache.CacheWrapper
}

// ИСПРАВЛЕНО: принимает CacheWrapper вместо RedisClient
func NewUserService(userRepo *repository.UserRepository, cacheWrapper *cache.CacheWrapper) *UserService {
	return &UserService{
		userRepo:     userRepo,
		cacheWrapper: cacheWrapper,
	}
}

// GetAllUsers с кешированием и ЛОГАМИ
func (s *UserService) GetAllUsers() ([]model.User, error) {
	var users []model.User

	// Пробуем получить из кэша
	err := s.cacheWrapper.Get("users:all", &users)
	if err == nil {
		log.Println("📦 DATA FROM CACHE") // ← ДОБАВИТЬ (из кэша)
		return users, nil
	}

	// Нет в кэше - идём в БД
	log.Println("💾 DATA FROM DATABASE") // ← ДОБАВИТЬ (из базы)
	users, err = s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш на 5 минут
	if len(users) > 0 {
		s.cacheWrapper.Set("users:all", users, 5*time.Minute)
	}

	return users, nil
}

// GetUserByID с кешированием и ЛОГАМИ
func (s *UserService) GetUserByID(id int) (*model.User, error) {
	var user model.User
	key := "user:" + string(rune(id))

	err := s.cacheWrapper.Get(key, &user)
	if err == nil {
		log.Println("📦 USER FROM CACHE") // ← ДОБАВИТЬ
		return &user, nil
	}

	log.Println("💾 USER FROM DATABASE") // ← ДОБАВИТЬ
	userFromDB, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	s.cacheWrapper.Set(key, userFromDB, 10*time.Minute)
	return userFromDB, nil
}

// CreateUser создаёт пользователя и очищает кеш
func (s *UserService) CreateUser(user *model.User) error {
	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// Очищаем кеш списка
	log.Println("🗑️ CLEARING USERS CACHE") // ← ДОБАВИТЬ
	s.cacheWrapper.Delete("users:all")
	return nil
}

// UpdateUser обновляет пользователя и очищает кеш
func (s *UserService) UpdateUser(user *model.User) error {
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Очищаем кеш
	s.cacheWrapper.Delete("user:" + string(rune(user.ID)))
	s.cacheWrapper.Delete("users:all")
	return nil
}

// DeleteUser удаляет пользователя и очищает кеш
func (s *UserService) DeleteUser(id int) error {
	if err := s.userRepo.Delete(id); err != nil {
		return err
	}

	// Очищаем кеш
	s.cacheWrapper.Delete("user:" + string(rune(id)))
	s.cacheWrapper.Delete("users:all")
	return nil
}
