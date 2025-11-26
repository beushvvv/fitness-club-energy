package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Главная функция - точка входа в программу
func main() {
	// Настраиваем обработчики URL
	http.HandleFunc("/", homeHandler)                   // Главная страница
	http.HandleFunc("/login", loginHandler)             // Страница входа
	http.HandleFunc("/register", registerHandler)       // Страница регистрации
	http.HandleFunc("/memberships", membershipsHandler) // Страница с абонементами
	http.HandleFunc("/about", aboutHandler)             // Страница "О нас"
	http.HandleFunc("/contact", contactHandler)         // Страница контактов

	// Запускаем веб-сервер на порту 8080
	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Обработчик главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос к главной странице
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Загружаем HTML шаблон
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		// Если шаблон не найден, показываем простую страницу
		simpleHomePage(w)
		return
	}

	// Отправляем HTML страницу пользователю
	tmpl.Execute(w, nil)
}

// Простая главная страница (если шаблон не работает)
func simpleHomePage(w http.ResponseWriter) {
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Фитнес-клуб "Энергия"</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 40px; }
			.menu { margin: 20px 0; }
			.menu a { margin-right: 15px; text-decoration: none; color: blue; }
			.menu a:hover { text-decoration: underline; }
		</style>
	</head>
	<body>
		<h1>Добро пожаловать в фитнес-клуб "Энергия"!</h1>
		
		<div class="menu">
			<a href="/">Главная</a>
			<a href="/about">О нас</a>
			<a href="/memberships">Абонементы</a>
			<a href="/contact">Контакты</a>
			<a href="/login">Вход</a>
			<a href="/register">Регистрация</a>
		</div>

		<h2>Наши услуги:</h2>
		<ul>
			<li>Современное кардио-оборудование</li>
			<li>Силовые тренажеры</li>
			<li>Групповые занятия (йога, пилатес)</li>
			<li>Персональные тренировки</li>
		</ul>

		<p>Присоединяйтесь к нам сегодня!</p>
	</body>
	</html>
	`)
}

// Обработчик страницы входа
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Показываем форму входа
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Вход в систему</title>
		</head>
		<body>
			<h1>Вход в личный кабинет</h1>
			
			<form method="POST">
				<div>
					<label>Email:</label>
					<input type="email" name="email" required>
				</div>
				<div>
					<label>Пароль:</label>
					<input type="password" name="password" required>
				</div>
				<button type="submit">Войти</button>
			</form>
			
			<br>
			<a href="/register">Нет аккаунта? Зарегистрируйтесь</a>
			<br>
			<a href="/">На главную</a>
		</body>
		</html>
		`)
	} else if r.Method == "POST" {
		// Обрабатываем данные формы
		// Получаем данные из формы
		email := r.FormValue("email")
		password := r.FormValue("password")

		fmt.Printf("Пользователь %s пытается войти (пароль: %s)\n", email, password)

		// Пока просто показываем сообщение об успехе
		fmt.Fprintf(w, `
		<h2>Вход выполнен!</h2>
		<p>Добро пожаловать, %s!</p>
		<p>Это упрощенная версия. В реальном приложении здесь была бы проверка пароля.</p>
		<a href="/">На главную</a>
		`, email)
	}
}

// Обработчик страницы регистрации
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Регистрация</title>
		</head>
		<body>
			<h1>Регистрация</h1>
			<form method="POST">
				<div>
					<label>Имя:</label>
					<input type="text" name="name" placeholder="Имя" required>
				</div>
				<div>
					<label>Email:</label>
					<input type="email" name="email" placeholder="Email" required>
				</div>
				<div>
					<label>Пароль:</label>
					<input type="password" name="password" placeholder="Пароль" required>
				</div>
				<div>
					<label>Телефон:</label>
					<input type="tel" name="phone" placeholder="Телефон">
				</div>
				<button type="submit">Зарегистрироваться</button>
			</form>
			<a href="/">На главную</a>
		</body>
		</html>
		`)
	} else if r.Method == "POST" {
		// Получаем данные из формы
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		phone := r.FormValue("phone")

		fmt.Printf("Новый пользователь: %s, email: %s, телефон: %s, пароль: %s\n",
			name, email, phone, password)

		// Показываем сообщение об успешной регистрации
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Регистрация успешна</title>
		</head>
		<body>
			<h2>Регистрация успешна!</h2>
			<p>Данные пользователя:</p>
			<ul>
				<li>Имя: %s</li>
				<li>Email: %s</li>
				<li>Телефон: %s</li>
			</ul>
			<p>Теперь вы можете <a href="/login">войти</a> в систему</p>
			<a href="/">На главную</a>
		</body>
		</html>
		`, name, email, phone)
	}
}

// Обработчик страницы с абонементами
func membershipsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Абонементы</title>
		<style>
			.membership { border: 1px solid #ccc; padding: 20px; margin: 10px; border-radius: 5px; }
			button { background: blue; color: white; padding: 10px; border: none; border-radius: 3px; }
		</style>
	</head>
	<body>
		<h1>Наши абонементы</h1>
		
		<div class="membership">
			<h2>Разовый визит</h2>
			<p><strong>Цена: 200 руб.</strong></p>
			<p>Однократное посещение всех зон клуба</p>
			<button>Купить</button>
		</div>
		
		<div class="membership">
			<h2>Месячный абонемент</h2>
			<p><strong>Цена: 3000 руб.</strong></p>
			<p>Неограниченное посещение в течение месяца</p>
			<button>Купить</button>
		</div>
		
		<div class="membership">
			<h2>Годовой абонемент</h2>
			<p><strong>Цена: 25000 руб.</strong></p>
			<p>Неограниченное посещение в течение года + 2 персональные тренировки</p>
			<button>Купить</button>
		</div>
		
		<br>
		<a href="/">На главную</a>
	</body>
	</html>
	`)
}

// Обработчик для страницы "О нас"
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<head>
		<title>О нас</title>
	</head>
	<body>
		<h1>Фитнес-клуб "Энергия"</h1>
		<p>Современный фитнес-клуб с лучшим оборудованием!</p>
		<h2>Наши преимущества:</h2>
		<ul>
			<li>✅ Современная кардио-зона</li>
			<li>✅ Профессиональные силовые тренажеры</li>
			<li>✅ Просторный зал для групповых занятий</li>
			<li>✅ Опытные тренеры</li>
			<li>✅ Чистые раздевалки и душевые</li>
		</ul>
		<a href="/">На главную</a>
	</body>
	</html>
	`)
}

// Обработчик для страницы контактов
func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Контакты</title>
	</head>
	<body>
		<h1>Контакты</h1>
		<p><strong>Телефон:</strong> +7 (123) 456-78-90</p>
		<p><strong>Email:</strong> info@energy-fitness.ru</p>
		<p><strong>Адрес:</strong> г. Москва, ул. Спортивная, д. 1</p>
		<p><strong>Часы работы:</strong> ежедневно с 7:00 до 23:00</p>
		<br>
		<a href="/">На главную</a>
	</body>
	</html>
	`)
}
