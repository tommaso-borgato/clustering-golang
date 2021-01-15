package node

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"fmt"
	"bytes"
)

// TODO: http://networkbit.ch/golang-ssh-client/ implementare anche "Multiple Command"

func (n *Node) Run(script string) error {
	log.Printf("IPV4: %s - publickeyFile: %s - script: %s", n.IPV4, n.PublickeyFile, script)

	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(n.PublickeyFile)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
		return err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
		return err
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
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", n.IPV4, n.SshPort), sshConfig)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
		return err
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
		return err
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on the remote side using the Run method
	// TODO: https://pkg.go.dev/golang.org/x/crypto/ssh#Session
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(script); err != nil {
		log.Fatal("Failed to run: " + err.Error())
		return err
	}
	fmt.Println(b.String())

	defer conn.Close()

	return nil
}