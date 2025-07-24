package routers

type routers struct {
	User        user
	Email       email
	Account     account
	Application application
	Group       group
}

var Routers = new(routers)
