package logic

type logics struct {
	User        user
	Email       email
	Account     account
	Application application
	Group       group
	Setting     setting
	File        file
}

var Logics = new(logics)
