package node

import (
    "golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"fmt"
	"os"
	"github.com/pkg/sftp"
	"path/filepath"
	"strings"
	"io"
)

// TODO: https://sftptogo.com/blog/go-sftp/ implementare anche la get

func (n *Node) Put(localFile, remoteFile string) error {	
	log.Printf("IPV4: %s - publickeyFile: %s - local: %s, remote: %s", n.IPV4, n.PublickeyFile, localFile, remoteFile)

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

	// Create new SFTP client
    sftpClient, err := sftp.NewClient(conn)
    if err != nil {
        log.Fatalf("Unable to start SFTP subsystem: %v\n", err)
        return err
    }
	defer sftpClient.Close()
	
	// open local file
	srcFile, err := os.Open(localFile)
    if err != nil {
        log.Fatalf("Unable to open local file: %v\n", err)
        return err
    }
	defer srcFile.Close()
	
	// Make remote directories recursion
    parent := filepath.Dir(remoteFile)
    path := string(filepath.Separator)
    dirs := strings.Split(parent, path)
    for _, dir := range dirs {
        path = filepath.Join(path, dir)
        sftpClient.Mkdir(path)
	}
	
	// Note: SFTP To Go doesn't support O_RDWR mode
    dstFile, err := sftpClient.OpenFile(remoteFile, (os.O_WRONLY|os.O_CREATE|os.O_TRUNC))
    if err != nil {
        log.Fatalf("Unable to open remote file: %v\n", err)
        return err
    }
	defer dstFile.Close()
	
	bytes, err := io.Copy(dstFile, srcFile)
    if err != nil {
        log.Fatalf("Unable to upload local file: %v\n", err)
        return err
    }
    log.Printf("%d bytes copied\n", bytes)

	defer conn.Close()

	return nil
}