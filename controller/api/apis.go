package api

type apis struct {
	User        user
	Email       email
	Account     account
	Application application
}

var Apis = new(apis)
