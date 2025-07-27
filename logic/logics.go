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
}

var Logics = new(logics)
