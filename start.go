package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

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

		amount := getAmountContainers(getSession(node, config))

		for number := 0; number < amount; number++ {
			startContainers(getSession(node, config), 4)
		}
	}

	fmt.Println("Ending application ======= goPods")
}

/*
	Starts containers in status stopped, based on quantity spent on subscription @{quantity}
*/
func startContainers(session *ssh.Session, amount int) {
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	command := fmt.Sprintf("docker unpause $(docker ps -qa --last %d  --filter status=exited --filter ancestor=highlanderdantas/snk-jiva-w:v1.5 --filter ancestor=highlanderdantas/snk-jiva-w:v1.4) ", amount)

	if err := session.Run(command); err != nil {
		log.Println("Erro:", err.Error())
	}

	fmt.Println("\n#############################")
	fmt.Println("Initiating ", amount, "containers W\n")
	fmt.Println(b.String())

}

/*
	Checks the number of containers in the node
*/
func getAmountContainers(session *ssh.Session) int {
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run("docker ps -qa --filter status=exited --filter ancestor=highlanderdantas/snk-jiva-w:v1.5 --filter ancestor=highlanderdantas/snk-jiva-w:v1.4 | wc -l"); err != nil {
		log.Println("Erro:", err.Error())
	}

	number := strings.Replace(b.String(), "\n", " ", 2)
	amount, _ := strconv.ParseInt(strings.TrimSpace(number), 10, 64)

	if int(amount) == 0 {
		fmt.Println("Nenhum container W pausado")
		os.Exit(0)
	}

	return int(amount)
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
