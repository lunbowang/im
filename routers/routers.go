package routers

type routers struct {
	User        user
	Email       email
	Account     account
	Application application
	Group       group
	Setting     setting
	Chat        ws
}

var Routers = new(routers)
