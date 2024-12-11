package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // underscore tells Go that you're importing it for its side effects, not because you need to use it.
)

type state struct {
	conf *config.Config
	db   *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	handler map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handler[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.handler[cmd.name]; ok {
		return f(s, cmd)
	}
	return errors.New(fmt.Sprintf("[!] Command '%v' doesn't exist!", cmd.name))
}

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
	return err
}

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

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", conf.DB_url)
	dbQueries := database.New(db)

	STATE := &state{conf: conf, db: dbQueries}
	COMMANDS := &commands{handler: make(map[string]func(*state, command) error)}
	COMMANDS.register("login", handlerLogin)
	COMMANDS.register("register", handlerRegister)
	_args := os.Args
	if len(_args) < 2 {
		log.Fatal("Forgot an argument?...")
	}

	COMMAND := command{name: _args[1], args: _args[2:]}
	if err := COMMANDS.run(STATE, COMMAND); err != nil {
		log.Fatal(err)
	}
	conf.SetUser(_args[2])
	conf, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)
}
