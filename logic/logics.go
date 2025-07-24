package logic

type logics struct {
	User        user
	Email       email
	Account     account
	Application application
	Group       group
}

var Logics = new(logics)
