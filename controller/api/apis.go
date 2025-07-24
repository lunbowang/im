package api

type apis struct {
	User        user
	Email       email
	Account     account
	Application application
	Group       group
}

var Apis = new(apis)
