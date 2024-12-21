package main

import (
	"context"
	"fmt"
	"gator/internal/database"
)

// get all users and print them in the format:
// * lane
// * allan (current)
// * hunter
func handlerGetUsers(s *state, cmd command) error {
	var (
		allUsers []database.User
		err      error
	)
	allUsers, err = s.db.GetUsers(context.Background())
	for _, curUser := range allUsers {
		fmt.Printf("* %v", curUser.Name)
		if curUser.Name == s.conf.Current_user_name {
			fmt.Printf(" (current)")
		}
		fmt.Println()
	}
	return err
}
