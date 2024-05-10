package server

import (
	"git.homebank.kz/homebank-oauth/halykid-events/controllers"
)

func (srv *impl) routes() {

	srv.router.GET("/health", controllers.Ping)

	srv.router.GET("/ws", srv.controller.Connect)

	srv.router.GET("/redirect-url", srv.controller.RedirectURL)
	srv.router.GET("/user-auth-status", srv.controller.UserAuthStatus)
}
