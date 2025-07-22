package routers

type routers struct {
	User    user
	Email   email
	Account account
}

var Routers = new(routers)
