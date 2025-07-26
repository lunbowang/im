package api

import "im/controller/api/chat"

type apis struct {
	User        user
	Email       email
	Account     account
	Application application
	Group       group
	Setting     setting
	Chat        chat.Group
}

var Apis = new(apis)
