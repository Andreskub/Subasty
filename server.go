package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args

	//Se debe indicar el port
	if len(arguments) == 1 {
		fmt.Println("Porfavor, ingrese numero de puerto.")
		return
	}

	PORT := ":" + arguments[1]
	//Se crea un socket aceptador en el puerto indicado
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	//Se crea un socket con el cliente aceptado
	socket_accepted, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		//Se recibe un mensaje del cliente
		netData, err := bufio.NewReader(socket_accepted).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		//Si el cliente envio STOP finaliza
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Finalizando servidor")
			return
		}
		//Imprime lo recibido y se lo envia al cliente
		fmt.Print("-> ", string(netData))
		fmt.Fprintf(socket_accepted, netData+"\n")
	}
}
