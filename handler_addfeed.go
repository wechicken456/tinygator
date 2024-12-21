package main

import (
	"context"
	"errors"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("Empty arguments!! Need a username.")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	var (
		curUser database.User
		err     error
	)

	curUser, err = s.db.GetUser(context.Background(), s.conf.Current_user_name)
	if err != nil {
		return err
	}

	// insert new feed into database
	feed := database.CreateFeedParams{}
	feed.ID = uuid.New()
	feed.CreatedAt = time.Now()
	feed.UpdatedAt = time.Now()
	feed.Name = feedName
	feed.Url = feedURL
	feed.UserID = curUser.ID

	_, err = s.db.CreateFeed(context.Background(), feed)
	return err
}
