package main

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/database"
	"net/url"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("[!] Bad arguments!! Need a feed URL.")
	}
	feedURL := cmd.args[0]

	if _, err := url.ParseRequestURI(feedURL); err != nil {
		return errors.New("[!] Invalid URL!!")
	}

	// check if a feed with feedURL is present in the database
	var feed database.Feed = database.Feed{}
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return err
	}

	// if the combination of this user and and this feed URL already exists in the database
	// return gracefully
	_, err = s.db.GetFeedFollowWithUserID(context.Background(), database.GetFeedFollowWithUserIDParams{
		UserID: user.ID,
		FeedID: feed.ID})
	if err != nil {
		return errors.New("[!] User is not following this feed.")
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID})
	if err != nil {
		return err
	}
	fmt.Printf("[+] User %v has unfollowed %v\n", user.Name, feedURL)
	return nil
}
