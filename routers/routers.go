package routers

type routers struct {
	User        user
	Email       email
	Account     account
	Application application
}

var Routers = new(routers)
