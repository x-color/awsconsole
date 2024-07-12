package awsconsole

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Instance struct {
	SsoURL       string        `json:"ssoUrl"`
	AccountRoles []AccountRole `json:"accountRoles"`
}

type Store struct {
	Instances []Instance `json:"instances"`
}

func (c *Store) addAndUpdateInstance(instance Instance) {
	for i, ins := range c.Instances {
		if ins.SsoURL == instance.SsoURL {
			c.Instances[i] = instance
			return
		}
	}
	c.Instances = append(c.Instances, instance)
}

func (c *Store) findInstance(ssoURL string) *Instance {
	for _, ins := range c.Instances {
		if ins.SsoURL == ssoURL {
			return &ins
		}
	}
	return nil
}

func loadStore() (*Store, error) {
	fp := defaultCacheFilePath()
	b, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	store := &Store{}
	if err := json.Unmarshal(b, store); err != nil {
		return nil, err
	}
	return store, nil
}

func saveStore(store *Store) error {
	b, err := json.Marshal(store)
	if err != nil {
		return err
	}
	fp := defaultCacheFilePath()
	return os.WriteFile(fp, b, 0644)
}

func defaultCacheFilePath() string {
	return filepath.Join(os.ExpandEnv("$HOME"), ".aws", "cli", "awsconsole.json")
}
