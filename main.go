package main

import (
	"errors"
	"fmt"
	"gator/internal/config"
	"log"
	"os"
)

type state struct {
	conf *config.Config
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Empty arguments!! Need a username.")
	}
	return s.conf.SetUser(cmd.args[0])
}

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	STATE := &state{conf: conf}
	COMMANDS := &commands{handler: make(map[string]func(*state, command) error)}
	COMMANDS.register("login", handlerLogin)
	_args := os.Args
	if len(_args) < 2 {
		log.Fatal("Forgot an argument?...")
	}

	COMMAND := command{name: _args[1], args: _args[2:]}
	if err := COMMANDS.run(STATE, COMMAND); err != nil {
		log.Fatal(err)
	}
	conf, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)
}
