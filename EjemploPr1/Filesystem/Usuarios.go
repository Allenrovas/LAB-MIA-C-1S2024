package Filesystem

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Login es una función que recibe un usuario y una contraseña e inicia sesión
func Login(userValor string, pwdValor string, idValor string) {
	//Verificar que el id exista en la lista de particiones montadas
	indice := VerificarParticionMontada(idValor)
	if indice == -1 {
		println("La partición no está montada")
		return
	}

	if Usr_sesion.Uid != -1 {
		println("Ya hay una sesión iniciada")
		return
	}

	MountActual := particionesMontadas[indice]
	SuperBlock := NewSuperBlock()

	//Leer el superbloque
	file, err := os.OpenFile("Discos/"+MountActual.LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al leer el disco")
		return
	}
	defer file.Close()

	//Leer el superbloque
	file.Seek(int64(MountActual.Start), 0)
	err = binary.Read(file, binary.LittleEndian, &SuperBlock)
	if err != nil {
		fmt.Println("Error al leer el superbloque")
		return
	}

	if !(SuperBlock.S_filesystem_type == 2 || SuperBlock.S_filesystem_type == 3) {
		println("El sistema de archivos no es 2fs ni 3fs o no está formateado")
		return
	}

	//Leer la ruta /users.txt

	//Verificar que el archivo /users.txt exista
	numeroInodo := BuscarArchivo("/users.txt", MountActual, SuperBlock, file)
	if numeroInodo == -1 {
		println("No se encontró el archivo /users.txt")
		return
	}

	//Leer el archivo /users.txt
	contenido := LeerArchivo(numeroInodo, SuperBlock, file)
	if contenido == "" {
		println("No se pudo leer el archivo /users.txt")
		return
	}

	//Dividir el archivo en líneas
	lineas := strings.Split(contenido, "\n")
	//Recorrer las líneas
	for _, linea := range lineas {
		if len(linea) == 0 {
			break
		}
		if linea[2] == 'U' || linea[2] == 'u' {
			in := strings.Split(linea, ",")
			if in[3] == userValor && in[4] == pwdValor {
				uid, _ := strconv.Atoi(in[0])
				Usr_sesion.Uid = int32(uid)
				Usr_sesion.Usr = userValor
				Usr_sesion.Pwd = pwdValor
				Usr_sesion.Pid = idValor
				Usr_sesion.Grp = in[2]
				break
			}
		}
		if len(linea) == 0 {
			break
		}
		if linea[2] == 'G' || linea[2] == 'g' {
			in := strings.Split(linea, ",")
			if in[1] == Usr_sesion.Grp {
				gid, _ := strconv.Atoi(in[0])
				Usr_sesion.Gid = int32(gid)
			}
		}
	}
	fmt.Println("Sesión iniciada con éxito en la partición: " + idValor + " con el usuario: " + userValor)
}

func BuscarArchivo(ruta string, MountActual Mount, SuperBlock SuperBlock, file *os.File) int {
	pathSplit := strings.Split(ruta, "/")
	var newPath []string
	for _, s := range pathSplit {
		if s != "" {
			newPath = append(newPath, s)
		}
	}

	pathSplit = newPath
	//Leer el inodo raíz
	inodoRaiz := NewInodes()
	file.Seek(int64(SuperBlock.S_inode_start), 0)
	err := binary.Read(file, binary.LittleEndian, &inodoRaiz)
	if err != nil {
		fmt.Println("Error al leer el inodo raíz")
		return -1
	}

	//Buscar el numero de inodo del archivo
	numeroInodo := BuscarIndiceInodo(inodoRaiz, pathSplit, SuperBlock, file)
	return numeroInodo
}

func BuscarIndiceInodo(inodo Inodes, pathSplit []string, SuperBlock SuperBlock, file *os.File) int {
	contador := 0
	if len(pathSplit) == 0 {
		return contador
	}
	actual := pathSplit[0]
	path := pathSplit[1:]
	for _, i := range inodo.I_block {
		if i != -1 {
			Desplazamiento := (SuperBlock.S_block_start) + (int32(i) * int32(binary.Size(Fileblock{})))
			file.Seek(int64(Desplazamiento), 0)
			var folder FolderBlock
			err := binary.Read(file, binary.LittleEndian, &folder)
			if err != nil {
				fmt.Println("Error al leer el bloque")
				return -1
			}
			for _, j := range folder.B_content {
				if j.B_inodo != -1 && strings.Contains(string(j.B_name[:]), actual) {
					if len(path) == 0 {
						return int(j.B_inodo)
					}
					//Bucar el siguiente inodo
					inodoSiguiente := NewInodes()
					file.Seek(int64(SuperBlock.S_inode_start)+int64(j.B_inodo*int32(binary.Size(Inodes{}))), 0)
					err := binary.Read(file, binary.LittleEndian, &inodoSiguiente)
					if err != nil {
						fmt.Println("Error al leer el inodo")
						return -1
					}
					return BuscarIndiceInodo(inodoSiguiente, path, SuperBlock, file)
				}
			}
		}
	}
	return -1
}

func LeerArchivo(numeroInodo int, SuperBlock SuperBlock, file *os.File) string {
	var contenido string
	inodo := NewInodes()
	file.Seek(int64(SuperBlock.S_inode_start+int32(numeroInodo)*int32(binary.Size(Inodes{}))), 0)
	err := binary.Read(file, binary.LittleEndian, &inodo)
	if err != nil {
		fmt.Println("Error al leer el inodo")
		return ""
	}

	if inodo.I_size == 0 {
		fmt.Println("La particion no tiene contenido")
		return ""
	}

	//Buscar el inodo del archivo
	for _, i := range inodo.I_block {
		if i != -1 {
			Desplazamiento := (SuperBlock.S_block_start) + (int32(i) * int32(binary.Size(Fileblock{})))
			var bloque Fileblock
			file.Seek(int64(Desplazamiento), 0)
			binary.Read(file, binary.LittleEndian, &bloque)
			lectura := strings.TrimRight(string(bloque.B_content[:]), string(rune(0)))
			lectura = ObtenerContenido(lectura, 64)
			contenido += lectura
		}
	}
	return contenido
}

func ObtenerContenido(contenido string, size int) string {
	var contenidoFinal string
	cantidadCaracteres := len(contenido)
	if cantidadCaracteres < size {
		contenidoFinal = contenido
	} else {
		for i := 0; i < size; i++ {
			contenidoFinal += string(contenido[i])
			contenido = contenido[1:]
		}
	}
	return contenidoFinal
}
