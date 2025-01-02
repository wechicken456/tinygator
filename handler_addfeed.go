package main

import (
	"context"
	"errors"
	"gator/internal/database"
	"net/url"
	"time"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("[!] Bad arguments!! Need 2 arguments: feed name and feed URL.")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	var (
		curUser database.User
		err     error
	)

	if _, err = url.ParseRequestURI(feedURL); err != nil {
		return errors.New("[!] Invalid URL!!")
	}

	curUser, err = s.db.GetUser(context.Background(), s.conf.Current_user_name)
	if err != nil {
		return err
	}

	// check if URL is already in feeds table
	feed := database.Feed{}
	feed, err = s.db.GetFeedByURL(context.Background(), feedURL)
	if err == nil {
		// if the combination of this user and and this feed URL already exists in the database
		// return gracefully
		_, err = s.db.GetFeedFollowWithUserID(context.Background(), database.GetFeedFollowWithUserIDParams{
			UserID: curUser.ID,
			FeedID: feed.ID})
		if err == nil {
			return nil
		}
	}

	// insert new feed into database
	feedCreate := database.CreateFeedParams{}
	feedCreate.ID = uuid.New()
	feedCreate.CreatedAt = time.Now()
	feedCreate.UpdatedAt = time.Now()
	feedCreate.Name = feedName
	feedCreate.Url = feedURL
	feedCreate.UserID = curUser.ID

	_, err = s.db.CreateFeed(context.Background(), feedCreate)
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    curUser.ID,
		FeedID:    feedCreate.ID,
	})
	return err
}
