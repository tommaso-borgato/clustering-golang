package node

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"fmt"
	"bytes"
)

type Node struct {
	PublickeyFile string
	IPV4 string
	User string
	SshPort int
}

func (n *Node) Ssh(script string) {
	log.Printf("IPV4: %s - publickeyFile: %s - script: %s", n.IPV4, n.PublickeyFile, script)

	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(n.PublickeyFile)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: n.User,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // accept any host key
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", n.IPV4, n.SshPort), sshConfig)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on the remote side using the Run method
	// TODO: https://pkg.go.dev/golang.org/x/crypto/ssh#Session
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(script); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())

	defer client.Close()
}