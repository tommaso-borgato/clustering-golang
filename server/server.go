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

	n := node.Node{ PublickeyFile: config["PUBLIC_KEY_FILE"], IPV4: "10.0.145.55", User: "hudson", SshPort: 22 }
	n.Ssh("ls -l")
}