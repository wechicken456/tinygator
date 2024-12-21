package main

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Empty arguments!! Need a username.")
	}

	// check if user already existsed
	_user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if _user != (database.User{}) {
		return errors.New(fmt.Sprintf("[!] User %v already exists!", cmd.args[0]))
	}

	// insert new user into database
	user := database.CreateUserParams{}
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Name = cmd.args[0]

	_, err = s.db.CreateUser(context.Background(), user)
	if err == nil {
		err = s.conf.SetUser(user.Name)
	}
	return err
}
