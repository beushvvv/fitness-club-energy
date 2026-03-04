package handler

import "github.com/gorilla/mux"

func SetupRoutes(
	r *mux.Router,
	userHandler *UserHandler,
	membershipHandler *MembershipHandler,
	workoutHandler *WorkoutHandler, // Добавлен новый параметр
) {
	api := r.PathPrefix("/api/v1").Subrouter()

	// Users
	api.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")    // Добавлено
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE") // Добавлено

	// Memberships
	api.HandleFunc("/memberships", membershipHandler.GetMemberships).Methods("GET")
	api.HandleFunc("/memberships/{id}", membershipHandler.GetMembership).Methods("GET") // Добавлено
	api.HandleFunc("/memberships", membershipHandler.CreateMembership).Methods("POST")
	api.HandleFunc("/memberships/{id}", membershipHandler.UpdateMembership).Methods("PUT") // Добавлено

	// Workouts
	api.HandleFunc("/workouts", workoutHandler.GetWorkouts).Methods("GET")
	api.HandleFunc("/workouts", workoutHandler.CreateWorkout).Methods("POST")
}
