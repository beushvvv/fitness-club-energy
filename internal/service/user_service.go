package service

import (
	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/logger"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"
	"strconv"
	"time"
)

type UserService struct {
	userRepo     repository.UserRepositoryInterface // изменено
	cacheWrapper cache.CacheWrapperInterface        // изменено
}

func NewUserService(userRepo repository.UserRepositoryInterface, cacheWrapper cache.CacheWrapperInterface) *UserService {
	return &UserService{
		userRepo:     userRepo,
		cacheWrapper: cacheWrapper,
	}
}

// GetAllUsers возвращает всех пользователей с кэшированием и логированием
func (s *UserService) GetAllUsers() ([]model.User, error) {
	var users []model.User
	sugar := logger.Log.Sugar()

	// Пытаемся получить данные из кэша
	err := s.cacheWrapper.Get("users:all", &users)
	if err == nil {
		// Данные найдены в кэше — логируем на уровне Debug
		sugar.Debugw("📦 DATA FROM CACHE", "key", "users:all")
		return users, nil
	}

	// Данных в кэше нет — идём в базу данных
	sugar.Debug("💾 DATA FROM DATABASE")
	users, err = s.userRepo.FindAll()
	if err != nil {
		// Ошибка при получении из БД — логируем на уровне Error
		sugar.Errorw("Failed to get users from DB", "error", err)
		return nil, err
	}

	// Сохраняем полученные данные в кэш
	if len(users) > 0 {
		s.cacheWrapper.Set("users:all", users, 5*time.Minute)
		sugar.Debugw("✅ Users cached", "count", len(users))
	}

	return users, nil
}

// GetUserByID возвращает пользователя по ID с кэшированием
func (s *UserService) GetUserByID(id int) (*model.User, error) {
	var user model.User
	sugar := logger.Log.Sugar()
	key := "user:" + strconv.Itoa(id) // исправлено

	// Пытаемся получить из кэша
	err := s.cacheWrapper.Get(key, &user)
	if err == nil {
		sugar.Debugw("📦 USER FROM CACHE", "user_id", id)
		return &user, nil
	}

	// Нет в кэше — идём в БД
	sugar.Debugw("💾 USER FROM DATABASE", "user_id", id)
	userFromDB, err := s.userRepo.FindByID(id)
	if err != nil {
		sugar.Errorw("Failed to get user by ID from DB", "user_id", id, "error", err)
		return nil, err
	}

	// Сохраняем в кэш
	s.cacheWrapper.Set(key, userFromDB, 10*time.Minute)
	sugar.Debugw("✅ User cached", "user_id", id)

	return userFromDB, nil
}

// CreateUser создаёт нового пользователя и очищает кэш
func (s *UserService) CreateUser(user *model.User) error {
	sugar := logger.Log.Sugar()

	if err := s.userRepo.Create(user); err != nil {
		// Ошибка создания — логируем с контекстом
		sugar.Errorw("Failed to create user",
			"error", err,
			"email", user.Email,
			"name", user.Name)
		return err
	}

	// Инвалидируем кэш списка пользователей
	s.cacheWrapper.Delete("users:all")
	sugar.Infow("User created and cache cleared",
		"user_id", user.ID,
		"email", user.Email,
		"name", user.Name)

	return nil
}

// UpdateUser обновляет данные пользователя
func (s *UserService) UpdateUser(user *model.User) error {
	sugar := logger.Log.Sugar()

	if err := s.userRepo.Update(user); err != nil {
		sugar.Errorw("Failed to update user",
			"user_id", user.ID,
			"error", err)
		return err
	}

	// Очищаем кэш конкретного пользователя и списка
	key := "user:" + strconv.Itoa(user.ID) // исправлено
	s.cacheWrapper.Delete(key)
	s.cacheWrapper.Delete("users:all")
	sugar.Infow("User updated and cache cleared",
		"user_id", user.ID,
		"email", user.Email)

	return nil
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(id int) error {
	sugar := logger.Log.Sugar()

	if err := s.userRepo.Delete(id); err != nil {
		sugar.Errorw("Failed to delete user",
			"user_id", id,
			"error", err)
		return err
	}

	// Очищаем кэш
	key := "user:" + strconv.Itoa(id) // исправлено
	s.cacheWrapper.Delete(key)
	s.cacheWrapper.Delete("users:all")
	sugar.Infow("User deleted and cache cleared", "user_id", id)

	return nil
}
