package routers

type routers struct {
	User        user
	Email       email
	Account     account
	Application application
	Group       group
	Setting     setting
	Chat        ws
	File        file
	Notify      notify
	Message     message
}

var Routers = new(routers)
