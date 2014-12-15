package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/Scalingo/cli/users"
	"gopkg.in/errgo.v1"
)

type AuthConfigData struct {
	LastUpdate        time.Time              `json:"last_update"`
	AuthConfigPerHost map[string]*users.User `json:"auth_config_data"`
}

func StoreAuth(user *users.User) error {
	// Check ~/.config/scalingo
	if _, err := os.Stat(C.ConfigDir); err != nil {
		if err, ok := err.(*os.PathError); ok {
			if err := os.MkdirAll(C.ConfigDir, 0755); err != nil {
				return errgo.Mask(err, errgo.Any)
			}
		} else {
			return errgo.Notef(err, "error reaching config directory")
		}
	}

	var authConfig AuthConfigData
	content, err := ioutil.ReadFile(C.AuthFile)
	if err == nil || (err != nil && os.IsNotExist(err)) {
		err = json.Unmarshal(content, &authConfig)
		if err != nil {
			return errgo.Mask(err)
		}

		// For backward compatiblity
		if authConfig.AuthConfigPerHost == nil {
			authConfig.AuthConfigPerHost = make(map[string]*users.User)
		}
	} else {
		return errgo.Mask(err)
	}

	authConfig.AuthConfigPerHost[C.apiHost] = user
	authConfig.LastUpdate = time.Now()

	file, err := os.OpenFile(C.AuthFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(authConfig); err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}

func LoadAuth() (*users.User, error) {
	file, err := os.OpenFile(C.AuthFile, os.O_RDONLY, 0644)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	var authConfig AuthConfigData
	dec := json.NewDecoder(file)
	if err := dec.Decode(&authConfig); err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	if user, ok := authConfig.AuthConfigPerHost[C.apiHost]; !ok {
		return nil, nil
	} else {
		return user, nil
	}
}