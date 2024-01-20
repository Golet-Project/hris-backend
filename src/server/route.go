package server

import (
	"hroost/server/middleware"
)

func (s *Server) route() {
	rest := s.presentation.rest

	// attendance
	attendance := s.app.Group("/attendance", middleware.Jwt())
	attendance.Post("/", rest.attendance.AddAttendance)
	attendance.Get("/", rest.attendance.FindAllAttendance)
	attendance.Put("/", rest.attendance.Checkout)
	attendance.Get("/today", rest.attendance.GetTodayAttendance)

	// auth
	auth := s.app.Group("/auth")
	auth.Post("/login", rest.auth.BasicAuthLogin)
	auth.Post("/forgot-password", rest.auth.ForgotPassword)
	auth.Post("/password-recovery/check", rest.auth.PasswordRecoveryCheck)
	auth.Put("/password", rest.auth.ChangePassword)

	// oauth
	oauth := s.app.Group("/oauth")
	oauth.Post("/google/login", rest.auth.OAuthLogin)
	oauth.Get("/google/callback", rest.auth.OAuthCallback)

	// employee
	employee := s.app.Group("/employee", middleware.Jwt())
	employee.Get("/", rest.employee.FindAllEmployee)
	employee.Post("/", rest.employee.CreateEmployee)

	// profile
	profile := s.app.Group("/profile", middleware.Jwt())
	profile.Get("/profile", rest.employee.GetProfile)

	// homepage
	homepage := s.app.Group("/homepage", middleware.Jwt())
	homepage.Get("/", middleware.Jwt(), rest.homepage.HomePage)

	// region
	region := s.app.Group("/region")
	region.Get("/provinces", rest.region.FindAllProvince)

	// tenant management
	tenantManagement := s.app.Group("/tenant", middleware.Jwt())
	tenantManagement.Post("/", rest.tenantManagement.CreateTenant)
}
