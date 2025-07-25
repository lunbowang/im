package logic

type logics struct {
	User        user
	Email       email
	Account     account
	Application application
	Group       group
	Setting     setting
}

var Logics = new(logics)
