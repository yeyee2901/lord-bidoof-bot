package telegram

import "google.golang.org/grpc/codes"

type TelegramError struct {
	GrpcCode codes.Code
	details  string
}

func (e *TelegramError) Error() string {
	return e.details
}

type ServerError struct {
	details string
}

func (e *ServerError) Error() string {
	return e.details
}
