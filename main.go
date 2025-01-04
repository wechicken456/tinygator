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

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		curUser, err := s.db.GetUser(context.Background(), s.conf.Current_user_name)
		if err != nil {
			return err
		}
		return handler(s, cmd, curUser)
	}
}

func main() {
	var (
		conf *config.Config
		err  error
		db   *sql.DB
	)
	conf, err = config.Read() // read config file
	if err != nil {
		log.Fatal(err)
	}

	// open connection to our database
	db, err = sql.Open("postgres", conf.DB_url)
	dbQueries := database.New(db)

	STATE := &state{conf: conf, db: dbQueries}
	COMMANDS := &commands{handler: make(map[string]func(*state, command) error)}
	COMMANDS.register("login", handlerLogin)
	COMMANDS.register("register", handlerRegister)
	COMMANDS.register("reset", handlerReset)
	COMMANDS.register("users", handlerGetUsers)
	COMMANDS.register("agg", handlerAgg)
	COMMANDS.register("feeds", handlerListAllFeeds)
	COMMANDS.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	COMMANDS.register("following", middlewareLoggedIn(handlerFollowingFeeds))
	COMMANDS.register("follow", middlewareLoggedIn(handlerFollow))
	COMMANDS.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	_args := os.Args
	if len(_args) < 2 {
		log.Fatal("Forgot an argument?...")
	}

	COMMAND := command{name: _args[1], args: _args[2:]}
	if err := COMMANDS.run(STATE, COMMAND); err != nil { // run command
		log.Fatal(err)
	}

}
