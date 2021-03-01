package service

import (
	"errors"
	"net/url"
)

type webhookJob struct {
	id string
	accounts []string
	recipient url.URL
}

type accountSet = map[string]webhookJob

func newAccountSet() accountSet {
	return make(map[string]webhookJob)
}

type cache struct {
	store map[string]accountSet
}

func newCache() cache {
	return cache{
		store: make(map[string]accountSet),
	}
}

func (c *cache) get(account string) (accountSet, bool) {
	accountSet, ok := c.store[account]
	return accountSet, ok
}

func (c *cache) add(job webhookJob) {
	for _, account := range job.accounts {
		accountSet, ok := c.store[account]
		if ok {
			accountSet[job.id] = job
			c.store[account] = accountSet
		} else {
			accountSet = newAccountSet()
			accountSet[job.id] = job
			c.store[account] = accountSet
		}
	}
}

func (c *cache) remove(job webhookJob) error {
	var err error = nil

	for _, account := range job.accounts {
		accountSet, ok := c.store[account]
		if !ok {
			err = errors.New("account not registered")
		} else {
			delete(accountSet, job.id)
		}
	}

	return err
}
