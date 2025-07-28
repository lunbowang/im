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
	File        file
	Notify      notify
	Message     message
}

var Apis = new(apis)
