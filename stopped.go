package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/ssh"
)

func main() {
	lambda.Start(initiating)
}

/*
	Starting Application
*/
func initiating() {
	fmt.Println("Starting Application ======= goPods ")

	nodes := getNodes()
	signer := getPPKKey()
	config := configureSSHClient(signer)

	for _, node := range nodes {

		session := getSession(node, config)
		stoppedContainers(session)
	}

}

/*
	Checks the number of containers in the node
*/
func stoppedContainers(session *ssh.Session) {
	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run("docker stop $(docker ps -qa)"); err != nil {
		log.Println("Erro:", err.Error())
	}

	fmt.Println(b.String())
}

/*
	SSH access settings
*/
func configureSSHClient(signer ssh.Signer) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: "rancher",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

/*
	Get the ssh key and set up a signer for connection
*/
func getPPKKey() ssh.Signer {
	pk, _ := ioutil.ReadFile("ssh-testes.pem")

	signer, err := ssh.ParsePrivateKey(pk)

	if err != nil {
		panic(err)
	}

	return signer
}

/*
	Returns all nodes that will start the containers
*/
func getNodes() []string {
	return []string{"34.204.90.20:22"}
}

/*
	Realiza a conexão ssh com o node especificado na assintura do metodo
*/
func getSession(node string, config *ssh.ClientConfig) *ssh.Session {

	client, err := ssh.Dial("tcp", node, config)

	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()

	if err != nil {
		panic("Failed to create session: " + err.Error())
	}

	return session
}
