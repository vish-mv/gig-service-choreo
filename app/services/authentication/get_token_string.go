package authentication

import (
	"GIG/app/constants/error_messages"
	"errors"

	"github.com/revel/revel"
)

func getTokenString(header *revel.RevelHeader, headerName string) (tokenString string, err error) {
	authHeader := header.Get(headerName)
	if authHeader == "" {
		return "", errors.New(error_messages.AuthHeaderNotFound)
	}
	return authHeader, nil

}
