package handlers

import "go.uber.org/dig"

type Routes struct {
	dig.In
	AuthHandler *AuthHandler
}
