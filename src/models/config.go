package models

import (
	"github.com/BurntSushi/toml"
	"fmt"
)

const (
	filepath string = "/go/src/github.com/stevenmahana/OrganizationServiceTemplate/config.toml"
)

type tomlConfig struct {
	Title   string
	Owner   ownerInfo
	DB      database `toml:"database"`
	Mongo   mongo `toml:"mongo"`
	Neo4j   neo4j `toml:"neo4j"`
	Postgres postgres `toml:"postgres"`
	Servers map[string]server
	Clients clients
}

type ownerInfo struct {
	Name  string
	Org   string `toml:"organization"`
	Email string
}

type database struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}

type mongo struct {
	Server  string
	Port    string
	User  	string
	Pass  	string
	Database string
	Enabled bool
}

type neo4j struct {
	Server  string
	Port    string
	User  	string
	Pass  	string
	Enabled bool
}

type postgres struct {
	Server  string
	Port    string
	User  	string
	Pass  	string
	Database string
	Enabled bool
}

type server struct {
	IP string
	DC string
}

type clients struct {
	Data  [][]interface{}
	Hosts []string
}

type (
	Configuration struct{
		session tomlConfig
	}
)


/*
	Service Configuration

	https://blog.gopheracademy.com/advent-2014/reading-config-files-the-go-way/
	https://github.com/BurntSushi/toml
*/
func Config() *Configuration {
	var config tomlConfig
	if _, err := toml.DecodeFile(filepath, &config); err != nil {
		panic(err)
	}

	return &Configuration{config}
}

func (c Configuration) test() {
	config := c.session

	fmt.Printf("Title: %s\n", config.Title)
	fmt.Printf("Owner: %s (%s, %s), Born: %s\n",
		config.Owner.Name, config.Owner.Org)
	fmt.Printf("Database: %s %v (Max conn. %d), Enabled? %v\n",
		config.DB.Server, config.DB.Ports, config.DB.ConnMax,
		config.DB.Enabled)
	for serverName, server := range config.Servers {
		fmt.Printf("Server: %s (%s, %s)\n", serverName, server.IP, server.DC)
	}
	fmt.Printf("Client data: %v\n", config.Clients.Data)
	fmt.Printf("Client hosts: %v\n", config.Clients.Hosts)
}
