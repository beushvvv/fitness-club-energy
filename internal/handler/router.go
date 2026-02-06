package handler

import "github.com/gorilla/mux"

func SetupRoutes(r *mux.Router, userHandler *UserHandler, membershipHandler *MembershipHandler) {
	api := r.PathPrefix("/api/v1").Subrouter()

	// Users
	api.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	// Memberships
	api.HandleFunc("/memberships", membershipHandler.GetMemberships).Methods("GET")
	api.HandleFunc("/memberships", membershipHandler.CreateMembership).Methods("POST") // Добавляем POST

	// Workouts
	api.HandleFunc("/workouts", GetWorkouts).Methods("GET")
	api.HandleFunc("/workouts", CreateWorkout).Methods("POST") // Добавляем POST
}
