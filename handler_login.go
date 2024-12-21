package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Empty arguments!! Need a username.")
	}
	// check if user already existsed
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return errors.New(fmt.Sprintf("[!] User %v doesn't exists!", cmd.args[0]))
	}
	return s.conf.SetUser(cmd.args[0])
}
