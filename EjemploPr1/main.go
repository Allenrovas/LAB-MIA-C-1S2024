package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Partition struct {
	Part_status [1]byte
	Part_type   [1]byte
	Part_fit    [1]byte
	Part_start  int32
	Part_size   int32
	Part_name   [16]byte
}

type MBR struct {
	Mbr_tamano         int32
	Mbr_fecha_creacion [19]byte
	Mbr_disk_signature int32
	Dsk_fit            [1]byte
	//Particiones [4]Partition
	Mbr_partition1 Partition
	Mbr_partition2 Partition
	Mbr_partition3 Partition
	Mbr_partition4 Partition
}

func NewPartition() Partition {
	return Partition{
		Part_status: [1]byte{'0'},
		Part_type:   [1]byte{'p'},
		Part_fit:    [1]byte{'w'},
		Part_start:  -1,
		Part_size:   -1,
		Part_name:   [16]byte{'~'},
	}
}

func NewMBR() MBR {
	return MBR{
		Mbr_tamano:         0,
		Mbr_fecha_creacion: [19]byte{},
		Mbr_disk_signature: 0,
		Dsk_fit:            [1]byte{'w'},
		Mbr_partition1:     NewPartition(),
		Mbr_partition2:     NewPartition(),
		Mbr_partition3:     NewPartition(),
		Mbr_partition4:     NewPartition(),
	}
}
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
			} else if valor == "rep" {
				fmt.Println("Ejecutando comando rep")
				Rep(&comandoSeparado)
			} else if valor == "\n" {
				continue
			} else if valor == "\r" {
				continue
			} else if valor == "" {
				continue
			} else {
				fmt.Println("Comando No reconocido")
			}
		}
	}
}

func analizarMkdisk(comandoSeparado *[]string) {
	//mkdisk -size=3000 -unit=K -fit
	*comandoSeparado = (*comandoSeparado)[1:]
	//Iterar sobre el comando separado
	var size, fit, unit bool
	//Variables para almacenar los valores de los parametros
	var sizeValor, fitValor, unitValor string
	fitValor = "f"
	unitValor = "m"
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		bandera := ObtenerBandera(valor)
		banderaValor := ObtenerBanderaValor(valor)
		if bandera == "-size" {
			size = true
			sizeValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-fit" {
			fit = true
			fitValor = banderaValor
			fitValor = strings.ToLower(fitValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-unit" {
			unit = true
			unitValor = banderaValor
			unitValor = strings.ToLower(unitValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("Parametro no reconocido: ", bandera)
		}
	}

	//Verificar si se ingresaron los parametros obligatorios
	if !size {
		fmt.Println("El parametro -size es obligatorio")
		return
	} else {
		if fit {
			if fitValor != "bf" && fitValor != "ff" && fitValor != "wf" {
				fmt.Println("El valor del parametro -fit no es valido")
				return
			} else {
				if fitValor == "bf" {
					fitValor = "b"
				} else if fitValor == "ff" {
					fitValor = "f"
				} else if fitValor == "wf" {
					fitValor = "w"
				}
			}
		}
		if unit {
			if unitValor != "k" && unitValor != "m" {
				fmt.Println("El valor del parametro -unit no es valido")
				return
			}
		}
		//Pasar a entero el valor del size
		sizeInt, err := strconv.Atoi(sizeValor)
		if err != nil {
			fmt.Println("El valor del parametro -size no es valido")
			return
		}
		if sizeInt <= 0 {
			fmt.Println("El valor del parametro -size no es valido")
			return
		}

		//Imprimir los valores de los parametros
		fmt.Println("Size: ", sizeValor)
		fmt.Println("Fit: ", fitValor)
		fmt.Println("Unit: ", unitValor)
		//Llamar a la funcion para crear el disco
		CrearDisco(sizeInt, fitValor, unitValor)
	}

}

func CrearDisco(sizeValor int, fitValor string, unitValor string) {
	//Tamano en bytes
	if unitValor == "k" && sizeValor != 0 {
		sizeValor = sizeValor * 1024
	} else if unitValor == "m" && sizeValor != 0 {
		sizeValor = sizeValor * 1024 * 1024
	} else {
		fmt.Println("El valor del parametro -unit no es valido")
		return
	}
	//Crear Directorio Discos si no existe para almacenar los discos

	//Si no existe el directorio Discos, entonces crearlo
	if _, err := os.Stat("Discos"); os.IsNotExist(err) {
		err = os.Mkdir("Discos", 0664)
		if err != nil {
			fmt.Println("Error al crear el directorio Discos: ", err)
			return
		}
	}
	//Contar la cantidad de discos para asignar el nombre
	archivos, err := ioutil.ReadDir("Discos")
	if err != nil {
		fmt.Println("Error al leer el directorio: ", err)
		return
	}
	//Declarar las letras del abecedario
	letras := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//Nombre del disco a partir de la cantidad de discos, por ejemplo A=1, B=2, C=3
	nombreDisco := string(letras[len(archivos)])
	//Crear el archivo del disco
	archivo, err := os.Create("Discos/" + nombreDisco + ".dsk")
	if err != nil {
		fmt.Println("Error al crear el archivo del disco: ", err)
		return
	}
	defer archivo.Close()

	//Escribir el MBR en el disco
	randomNum := rand.Intn(99) + 1
	var disk MBR
	fmt.Println("Size: ", sizeValor)
	disk.Mbr_tamano = int32(sizeValor)
	disk.Mbr_disk_signature = int32(randomNum)
	fitAux := []byte(fitValor)
	disk.Dsk_fit = [1]byte{fitAux[0]}
	fechaActual := time.Now()
	fecha := fechaActual.Format("2006-01-02 15:04:05")
	copy(disk.Mbr_fecha_creacion[:], fecha)

	//Escribir en las particiones

	disk.Mbr_partition1.Part_status = [1]byte{'0'}
	disk.Mbr_partition2.Part_status = [1]byte{'0'}
	disk.Mbr_partition3.Part_status = [1]byte{'0'}
	disk.Mbr_partition4.Part_status = [1]byte{'0'}

	disk.Mbr_partition1.Part_type = [1]byte{'0'}
	disk.Mbr_partition2.Part_type = [1]byte{'0'}
	disk.Mbr_partition3.Part_type = [1]byte{'0'}
	disk.Mbr_partition4.Part_type = [1]byte{'0'}

	disk.Mbr_partition1.Part_fit = [1]byte{'0'}
	disk.Mbr_partition2.Part_fit = [1]byte{'0'}
	disk.Mbr_partition3.Part_fit = [1]byte{'0'}
	disk.Mbr_partition4.Part_fit = [1]byte{'0'}

	disk.Mbr_partition1.Part_start = 0
	disk.Mbr_partition2.Part_start = 0
	disk.Mbr_partition3.Part_start = 0
	disk.Mbr_partition4.Part_start = 0

	disk.Mbr_partition1.Part_size = 0
	disk.Mbr_partition2.Part_size = 0
	disk.Mbr_partition3.Part_size = 0
	disk.Mbr_partition4.Part_size = 0

	disk.Mbr_partition1.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	disk.Mbr_partition2.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	disk.Mbr_partition3.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	disk.Mbr_partition4.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}

	bufer := new(bytes.Buffer)
	for i := 0; i < 1024; i++ {
		bufer.WriteByte(0)
	}

	var totalBytes int = 0
	for totalBytes < int(sizeValor) {
		c, err := archivo.Write(bufer.Bytes())
		if err != nil {
			fmt.Println("Error al escribir en el archivo: ", err)
			return
		}
		totalBytes += c
	}
	fmt.Println("Archivo llenado con 0s")
	//Escribir el MBR en el disco
	archivo.Seek(0, 0)
	err = binary.Write(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al escribir el MBR en el disco: ", err)
		return
	}
	fmt.Println("Disco", nombreDisco, "creado con exito")
}

func Rep(comandoSeparado *[]string) {
	*comandoSeparado = (*comandoSeparado)[1:]
	//Abrir el disco A
	archivo, err := os.Open("Discos/A.dsk")
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()
	disk := NewMBR()
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		return
	}
	fmt.Println("Tamaño: ", disk.Mbr_tamano)
	fmt.Println("Fecha: ", string(disk.Mbr_fecha_creacion[:]))
	fmt.Println("Signature: ", disk.Mbr_disk_signature)
	fmt.Println("Fit: ", string(disk.Dsk_fit[:]))
	fmt.Println("Partition1: ", string(disk.Mbr_partition1.Part_status[:]))
	fmt.Println("Partition2: ", string(disk.Mbr_partition2.Part_status[:]))
	fmt.Println("Partition3: ", string(disk.Mbr_partition3.Part_status[:]))
	fmt.Println("Partition4: ", string(disk.Mbr_partition4.Part_status[:]))
}

func ObtenerBandera(bandera string) string {
	//mkdisk -size=3000 -unit=K
	var banderaValor string
	for _, valor := range bandera {
		if valor == '=' {
			break
		}
		banderaValor += string(valor)
	}
	banderaValor = strings.ToLower(banderaValor)
	return banderaValor
}

func ObtenerBanderaValor(bandera string) string {
	//mkdisk -size=3000 -unit=K
	var banderaValor string
	var banderaEncontrada bool
	for _, valor := range bandera {
		if banderaEncontrada {
			banderaValor += string(valor)
		}
		if valor == '=' {
			banderaEncontrada = true
		}
	}
	return banderaValor
}
