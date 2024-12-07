package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

// struct fields MUST be captialized for the json encoder (which lives in another package) to see them.
// then we use field `tags` to specify how we want to the corresponding field in json.
// in this case, since we don't want to captitalize the struct fields, we have to use tags.
type Config struct {
	DB_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (conf *Config) getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + configFileName, nil
}

// read ~/.gatorconfig.json and return the corresponding `Config` struct in it.
func Read() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(homeDir + "/" + configFileName)
	if err != nil {
		return nil, err
	}

	var conf *Config = &Config{}
	dec := json.NewDecoder(file)
	err = dec.Decode(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (conf *Config) SetUser(new_user_name string) error {
	conf.Current_user_name = new_user_name

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(homeDir+"/"+configFileName, os.O_RDWR|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	err = enc.Encode(conf)
	if err != nil {
		return err
	}
	return nil
}

func (conf *Config) write(cfg Config) error {
	conf.DB_url = cfg.DB_url
	conf.Current_user_name = cfg.Current_user_name
	return nil
}
