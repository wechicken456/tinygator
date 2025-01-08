# tinygator

My tiny [RSS](https://en.wikipedia.org/wiki/RSS) feed aggregator in Go. RSS feeds are a way for websites to publish updates to their content. You can use this project to keep up with your favorite blogs, news sites, comic books, and more!

Technology stack: `go` standard library, PostgreSQL as the database, [`goose`](https://github.com/pressly/goose) migration, [`sqlc`](https://docs.sqlc.dev/en/latest/overview/install.html) (to generate Go code that our application can use to interact with the database).

It is be a CLI tool that allows users to:
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post
- Give live update of followed feeds.


## Installation
Pre-req: Postgres and Go. 

1. To install `gator`, first clone this repo.
2. `cd` into the cloned repo.
3. `go build`. This will compile a binary called `gator` in the current directory. Verify that it works before continuing.
4. `go list -f '{{.Target}}'`. This will print where the `go install` command will install the packages. Note this path.
5. Add the path in step 4 to your `$PATH` environment variable if it isn't there.
6. `go install`. You should now be able to invoke `gator` in any directory.

## Config file

By default, `gator` assumes the config file is created/found at `~/.gatorconfig.json`. You can change its name by modifying the `configFileName` variable in the file `config/config.go`.

## Commands

For a clean install, run `gator reset` first.

Then, run `gator register <your-name>`. If the name already existed, run `gator login <your-name>`.

To list all users, `gator users`.

To add a feed, `gator addfeed <feed-name> <feed-url>`.

To follow a feed, `gator follow <feed-url>`.

To unfollow a feed, `gator unfollow <feed-url>`.

To list all feeds (you don't need to be logged in for this), `gator feeds`.

To list only your following feeds (needs to be logged in first), `gator following`.

To start aggregating, `go agg <agg-interval>`. The interval should be in the format `1s`, `1m`, `1h`, etc.

## Contribution
Feel free to contribute!

`sqlc` should always run from the root of the project. Its configuration file can be found in `sqlc.yaml`. By default, it reads SQL queries from `sql/schema` and `sql/queries` and generates code into `internal/database`.

You can reach me via my email tinvuong2003@gmail.com or [Twitter](https://x.com/pwnPHOfun).