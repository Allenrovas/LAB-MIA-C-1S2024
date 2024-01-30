package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("------------------------")
	fmt.Println("-- Ejemplo Proyecto 1 --")
	fmt.Println("------------------------")
	fmt.Println("--Allen Roman-202004745-")
	fmt.Println("------------------------")

	for {
		leerComando()
	}
}

func leerComando() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Ingrese un comando: ")
	comando, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al ingresar el comando: ", err)
		return
	}

	comando = strings.TrimSpace(comando)

	//Funcion para analizar
	fmt.Println(comando)

	analizar(comando)
}

func analizar(comando string) {
	//mkdisk -size=3000 -unit=K login
	//Mkdisk
	//MKDISK
	//login
	//#inicio mkdisk del proyecto1

	//var comando ----> comando := "Hola este es un comando" ----> cambio en memoria
	//Apuntador al cambiar la variable, la cambia para todo el programa
	comandoSeparado := strings.Split(comando, " ")
	if strings.Contains(comandoSeparado[0], "#") {
		//Imprimir el comentario
		fmt.Println("Comentario: ")
		//Eliminar el # del comentario
		comandoSeparado[0] = strings.Replace(comandoSeparado[0], "#", "", -1)
		for _, comentario := range comandoSeparado {
			fmt.Println(comentario + " ")
		}
	} else {
		//Si no es un comentario, entonces es un comando
		//Iterar sobre el comando separado
		for _, valor := range comandoSeparado {
			//el primer valor del comando lo pasamos a minusculas
			valor = strings.ToLower(valor)
			//Si el valor es igual a mkdisk, entonces es un comando de creacion de disco
			if valor == "mkdisk" {
				fmt.Println("Ejecutando comando mkdisk")
				//Analizar Comando Mkdisk
				analizarMkdisk(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				analizar(comandoSeparadoString)
			} else if valor == "\n" {
				continue
			} else if valor == "\r" {
				continue
			} else if valor == "" {
				continue
			} else {
				fmt.Println("ComandoNoreconocido")
			}
		}
	}
}

func analizarMkdisk(comandoSeparado *[]string) {
	//-size -unit
	*comandoSeparado = (*comandoSeparado)[1:]

}
