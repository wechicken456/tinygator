package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"html"
	"io"
	"log"
	"net/http"
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
	if err == nil {
		err = s.conf.SetUser(user.Name)
	}
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

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetDatabase(context.Background())
	return err
}

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

func parseXML(resp *http.Response, feed *RSSFeed) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(body, feed)
	if err != nil {
		return err
	}

	// unescape HTML characters in the struct fields in both Channel and Items.
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Link = html.UnescapeString(feed.Channel.Link)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := 0; i < len(feed.Channel.Item); i++ {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Link = html.UnescapeString(feed.Channel.Item[i].Link)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
		feed.Channel.Item[i].PubDate = html.UnescapeString(feed.Channel.Item[i].PubDate)
	}
	return err
}

func fetchFeed(_ctx context.Context, feedURL string) (*RSSFeed, error) {
	var (
		client *http.Client
		ctx    context.Context
		cancel context.CancelFunc
		req    *http.Request
		resp   *http.Response
		err    error
		feed   *RSSFeed = &RSSFeed{}
	)
	ctx, cancel = context.WithTimeout(_ctx, 3*time.Second) // request timeout after 3s
	defer cancel()

	client = &http.Client{}
	req, err = http.NewRequestWithContext(ctx, "GET", feedURL, nil) // create a new request with the timeout context
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	c := make(chan error, 1)
	go func() {
		var _err error
		resp, _err = client.Do(req) // send the request
		if _err != nil {
			c <- _err
			return
		}

		defer resp.Body.Close()
		err = parseXML(resp, feed)
		c <- _err
	}()

	select {
	case <-ctx.Done(): // case timeout
		go func() { <-c }()
		return feed, ctx.Err()
	case err = <-c: // case finished parsing before timeout
		return feed, err
	}
}

func handlerAgg(s *state, cmd command) error {
	var (
		feed *RSSFeed
		err  error
	)
	feed, err = fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return err
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

	__args := os.Args
	if len(__args) < 2 {
		log.Fatal("Forgot an argument?...")
	}

	var _args []string = nil // argument to command
	if len(__args) > 2 {
		_args = __args[2:]
	}

	COMMAND := command{name: __args[1], args: _args}
	if err := COMMANDS.run(STATE, COMMAND); err != nil { // run command
		log.Fatal(err)
	}

	conf, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)
}
