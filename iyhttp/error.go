package iyhttp

import (
	"net/http"
)

// Common standard errors.
var (
	ErrApplicationError = Error{
		Message:    "application error",
		StatusCode: http.StatusInternalServerError,
	}
)

// Error is an error that also contains an HTTP status code and a
// message for a client.
//
// An Error can optionally have some context, which is to be used for
// for providing more information about the error for internal logging
// purposes. Context will not be encoded when presenting the error to a
// client or user.
//
type Error struct {
	Context    string `json:"-"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// Code returns the HTTP status code associated with this Error.
// If the StatusCode was not set on the Error, then
// http.StatusInternalServerError is returned.
func (e Error) Code() int {
	if e.StatusCode == 0 {
		return http.StatusInternalServerError
	}
	return e.StatusCode
}

// Error implements the error interface and returns the Error's Context
// if it has one, or Message if it does not.
func (e Error) Error() string {
	if e.Context != "" {
		return e.Context
	}
	return e.Message
}
