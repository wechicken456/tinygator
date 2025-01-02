package main

import (
	"context"
	"fmt"
	"gator/internal/database"
)

// not sorted - print the feeds as they appear in the table.
func handlerListFeeds(s *state, cmd command) error {
	var (
		allFeeds []database.Feed
		curFeed  database.Feed
		err      error
		user     database.User
	)

	allFeeds, err = s.db.GetFeeds(context.Background())
	for _, curFeed = range allFeeds {
		user, _ = s.db.GetUserById(context.Background(), curFeed.UserID)

		fmt.Printf("Feed ID: %v\n", curFeed.ID)
		fmt.Printf("Created by user: %v\n", user.Name)
		fmt.Printf("CreatedAt: %v\n", curFeed.CreatedAt)
		fmt.Printf("UpdatedAt: %v\n", curFeed.UpdatedAt)
		fmt.Printf("Name: %v\n", curFeed.Name)
		fmt.Printf("Url: %v\n", curFeed.Url)
		fmt.Println()
	}
	return err
}
