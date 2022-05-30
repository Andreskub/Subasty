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

	//Se debe indicar el host y el port
	if len(arguments) == 1 {
		fmt.Println("Porfavor ingrese host:port.")
		return
	}

	dir := arguments[1]
	//Se crea y conecta un socket al servidor en la direccion especificada
	socket_connected, err := net.Dial("tcp", dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		//Se lee de entrada estandar
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		//Se envia al servidor
		fmt.Fprintf(socket_connected, text+"\n")

		//Se recibe una respuesta del servidor
		message, _ := bufio.NewReader(socket_connected).ReadString('\n')
		fmt.Print("->: " + message)

		//Si el usuario escribio STOP, finaliza
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("Finalizando cliente")
			return
		}
	}
}
