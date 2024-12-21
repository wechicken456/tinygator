package main

import "context"

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetDatabase(context.Background())
	return err
}
