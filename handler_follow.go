package main

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/database"
	"net/url"
	"time"

	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("[!] Bad arguments!! Need a feed URL.")
	}
	feedURL := cmd.args[0]

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

	// check if a feed with feedURL is present in the database
	var feed database.Feed = database.Feed{}
	feed, err = s.db.GetFeedByURL(context.Background(), feedURL)
	if err == nil {
		// if the combination of this user and and this feed URL already exists in the database
		// return gracefully
		_, err = s.db.GetFeedFollowWithUserID(context.Background(), database.GetFeedFollowWithUserIDParams{
			UserID: curUser.ID,
			FeedID: feed.ID})
		fmt.Println(err)
		if err == nil {
			return nil
		}
	} else {
		return err
	}

	feedFollow := database.CreateFeedFollowParams{}
	feedFollow.ID = uuid.New()
	feedFollow.CreatedAt = time.Now()
	feedFollow.UpdatedAt = time.Now()
	feedFollow.UserID = curUser.ID
	feedFollow.FeedID = feed.ID

	reply := database.CreateFeedFollowRow{}
	reply, err = s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return err
	}
	fmt.Printf("[+] User %v is now following %v\n", reply.Username, reply.Feedname)
	return err
}
