package main

import (
	"bytes"
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
	log.Println("Starting Application ======= goPods ")

	nodes := getNodes()
	signer := getPPKKey()
	config := configureSSHClient(signer)

	for _, node := range nodes {

		session := getSession(node, config)
		stoppedContainers(session)
	}

	log.Println("Ending application ======= goPods")

}

/*
	Checks the number of containers in the node
*/
func stoppedContainers(session *ssh.Session) {
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run("docker pause $(docker ps -qa --filter status=running --filter ancestor=highlanderdantas/snk-jiva-w:v1.5 --filter ancestor=highlanderdantas/snk-jiva-w:v1.4)"); err != nil {
		log.Println("Nenhum container a para ser pausado")
	} else {
		log.Println("#############################")
		log.Println("Pausing containers W\n")
		log.Println(b.String())
	}
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
	return []string{"54.236.187.107:22", "3.211.88.71:22"}
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
