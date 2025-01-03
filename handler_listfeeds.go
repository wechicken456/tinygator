package main

import (
	"context"
	"fmt"
	"gator/internal/database"
)

// not sorted - print the feeds (that the current user is following) as they appear in the table.
func handlerFollowingFeeds(s *state, cmd command, user database.User) error {
	var (
		allFeeds []database.GetFeedFollowsForUserRow
		curFeed  database.GetFeedFollowsForUserRow
		err      error
	)

	allFeeds, err = s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	fmt.Printf("[+] User %v is following:\n\n", user.Name)
	for _, curFeed = range allFeeds {
		fmt.Printf("Feed ID: %v\n", curFeed.ID)
		fmt.Printf("Feed name: %v\n", curFeed.Feedname)
		fmt.Printf("Created by user: %v\n", user.Name)
		fmt.Printf("CreatedAt: %v\n", curFeed.CreatedAt)
		fmt.Printf("UpdatedAt: %v\n", curFeed.UpdatedAt)
		fmt.Println("--------------------")
	}
	return err
}

// not sorted - print ALL feeds as they appear in the table.
func handlerListAllFeeds(s *state, cmd command) error {
	var (
		allFeeds []database.Feed
		curFeed  database.Feed
		err      error
		curUser  database.User
	)

	allFeeds, err = s.db.GetFeeds(context.Background())
	for _, curFeed = range allFeeds {
		curUser, _ = s.db.GetUserById(context.Background(), curFeed.UserID)

		fmt.Printf("Feed ID: %v\n", curFeed.ID)
		fmt.Printf("Created by user: %v\n", curUser.Name)
		fmt.Printf("CreatedAt: %v\n", curFeed.CreatedAt)
		fmt.Printf("UpdatedAt: %v\n", curFeed.UpdatedAt)
		fmt.Printf("Name: %v\n", curFeed.Name)
		fmt.Printf("Url: %v\n", curFeed.Url)
		fmt.Println()
	}
	return err
}
