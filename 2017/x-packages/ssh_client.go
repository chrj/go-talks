package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

// START HOSTKEY OMIT

var hostKey, _, _, _, _ = ssh.ParseAuthorizedKey([]byte("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBIwT1IFZq5h0gs0FigWvTPJYFvK+3CMZj+1PiIFzwgDW09s3jnn5veNrN2oW/WZLbDY67zsg+70aC8kfCumqY40="))

// END HOSTKEY OMIT

func loadRSAKey() ssh.Signer {

	key, err := ioutil.ReadFile("/home/razor/.ssh/id_ed25519")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	return signer
}

func main() {

	// START CLIENT SETUP OMIT

	config := &ssh.ClientConfig{
		User: "razor",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(loadRSAKey()),
			ssh.Password("hunter2"),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	client, err := ssh.Dial("tcp", "cerebro.technobabble.dk:22", config)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	// END CLIENT SETUP OMIT

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		log.Fatalf("failed to run: %v", err)
	}

	fmt.Println(b.String())

}
