package models

import "errors"

var (
	ErrLoadEnvFailed      = errors.New("failed to load environment")
	ErrConnectionDBFailed = errors.New("failed to connect to database")
	ErrServerFailed       = errors.New("failed to connect to server")
	ErrClientFailed       = errors.New("failed to create client")
)

var (
	ErrInvalidPassword        = errors.New("invalid password")
	ErrUserDeleteResponse     = errors.New("failed to delete user")
	ErrChangeUserInfoResponse = errors.New("failed to change user info")
	ErrGetUserResponse        = errors.New("failed to get user")
	ErrUserExists             = errors.New("user already exists")
	ErrCreateUserResponse     = errors.New("failed to create user")
)

var (
	ErrCreateTaskResponse = errors.New("failed to create task")
	ErrGetTaskResponse    = errors.New("failed to get task")
	ErrCheckTimerStatus   = errors.New("failed to check timer status")
	ErrTimerStarted       = errors.New("timer already started")
	ErrStartTimer         = errors.New("failed to start timer")
	ErrStopTimer          = errors.New("failed to stop timer")
)
