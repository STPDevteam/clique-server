package config

import (
	"github.com/patrickmn/go-cache"
)

var (
	GitServer  string = "unknown"
	ConfigFile string
	ServerName string = "myClique"
	ServerMark string

	MyCache *cache.Cache
)
