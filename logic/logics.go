package logic

type logics struct {
	User        user
	Email       email
	Account     account
	Application application
	Group       group
	Setting     setting
	File        file
	Notify      notify
	Message     message
}

var Logics = new(logics)
