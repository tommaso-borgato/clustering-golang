package main

import (
	"log"
	"fmt"
	"github.com/tommaso-borgato/clustering-golang/config"
	"github.com/tommaso-borgato/clustering-golang/node"
)

const configFile = "clustering.properties"

func main() { 
	log.Println("START SERVER") 

	config, err := config.ReadConfig( "clustering.properties" )
	if err != nil {
		panic(fmt.Sprintf("cannot find configuration file %s", configFile))
	}

	node := node.Node{ PublickeyFile: config["PUBLIC_KEY_FILE"], IPV4: "10.0.145.55", User: "hudson", SshPort: 22 }

	// ssh
	err = node.Run("ls -l")
	if err != nil {
		panic(fmt.Sprintf("cannot run 'ls -l' on node %s", node.IPV4))
	}

	// sftp
	err = node.Put("/tmp/prova.txt", "/tmp/prova.txt")
	if err != nil {
		panic(fmt.Sprintf("cannot put '/tmp/prova.txt' to node %s", node.IPV4))
	}
}