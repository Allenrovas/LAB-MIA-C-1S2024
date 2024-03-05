package main

import (
	"LAB-MIA-C-1S2024/EjemploPr1/Filesystem"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("------------------------")
	fmt.Println("-- Ejemplo Proyecto 1 --")
	fmt.Println("------------------------")
	fmt.Println("--Allen Roman-202004745-")
	fmt.Println("------------------------")
	//BuscarMontadas()
	//Leer el MBR de cada disco,
	//ir buscando en cada particion el status para saber si una particion es montada
	//Si es montada, agregarla a la lista de particiones montadas
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
			} else if valor == "rmdisk" {
				fmt.Println("Ejecutando comando rmdisk")
				//Analizar Comando Rmdisk
				analizarRmdisk(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				analizar(comandoSeparadoString)
			} else if valor == "fdisk" {
				fmt.Println("Ejecutando comando fdisk")
				//Analizar Comando Fdisk
				analizarFdisk(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				analizar(comandoSeparadoString)
			} else if valor == "rep" {
				fmt.Println("Ejecutando comando rep")
				Filesystem.ReporteDisk(&comandoSeparado)
			} else if valor == "mount" {
				fmt.Println("Ejecutando comando mount")
				//Analizar Comando Mount
				analizarMount(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				analizar(comandoSeparadoString)
			} else if valor == "mkfs" {
				fmt.Println("Ejecutando comando mkfs")
				//Analizar Comando Mkfs
				analizarMkfs(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				analizar(comandoSeparadoString)
			} else if valor == "login" {
				fmt.Println("Ejecutando comando login")
				//Analizar Comando Login
				analizarLogin(&comandoSeparado)
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
		//CrearDisco(sizeInt, fitValor, unitValor)
		Filesystem.CrearDisco(sizeInt, fitValor, unitValor)
	}

}

func analizarRmdisk(comandoSeparado *[]string) {
	//rmdisk -driveletter=A
	*comandoSeparado = (*comandoSeparado)[1:]
	//Iterar sobre el comando separado
	var driveletter string
	var drive bool
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		bandera := ObtenerBandera(valor)
		banderaValor := ObtenerBanderaValor(valor)
		if bandera == "-driveletter" {
			driveletter = banderaValor
			driveletter = strings.ToUpper(driveletter)
			drive = true
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("Parametro no reconocido: ", bandera)
		}
	}
	//Verificar si se ingresaron los parametros obligatorios
	if !drive {
		fmt.Println("El parametro -driveletter es obligatorio")
		return
	} else {
		//Imprimir los valores de los parametros
		fmt.Println("Driveletter: ", driveletter)
		//Llamar a la funcion para eliminar el disco
		//Buscar el disco con la letra en el directorio Discos
		//EliminarDisco(driveletter)
		//os.Remove("Discos/" + driveletter + ".dsk")
	}
}

func analizarFdisk(comandoSeparado *[]string) {
	//fdisk -size=300 -driveletter=A -name=Particion1
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingresaron los parametros
	var banderaLetter, banderaName, banderaFit, banderaUnit, banderaType, banderaDelete, banderaAdd bool
	//Variables para almacenar los valores de los parametros
	var sizeValor, letterValor, nameValor, fitValor, unitValor, typeValor, deleteValor, addValor string
	//Setear valores por defecto
	fitValor = "w"
	unitValor = "k"
	typeValor = "p"
	deleteValor = "0"
	addValor = "0"
	sizeValor = "0"
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		//-size
		bandera := ObtenerBandera(valor)
		//300
		banderaValor := ObtenerBanderaValor(valor)
		if bandera == "-size" {
			sizeValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-driveletter" {
			banderaLetter = true
			letterValor = banderaValor
			letterValor = strings.ToUpper(letterValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-name" {
			banderaName = true
			nameValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-fit" {
			banderaFit = true
			fitValor = banderaValor
			fitValor = strings.ToLower(fitValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-unit" {
			banderaUnit = true
			unitValor = banderaValor
			unitValor = strings.ToLower(unitValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-type" {
			banderaType = true
			typeValor = banderaValor
			typeValor = strings.ToLower(typeValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-delete" {
			banderaDelete = true
			deleteValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-add" {
			banderaAdd = true
			addValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("Parametro no reconocido: ", bandera)
		}
	}
	//Obligatorios: -size(al crear), driveletter, name
	//Verificar si se ingresaron los parametros obligatorios
	if !banderaLetter {
		fmt.Println("El parametro -driveletter es obligatorio")
		return
	} else if !banderaName {
		fmt.Println("El parametro -name es obligatorio")
		return
	} else {
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

		if banderaFit {
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
		if !banderaUnit {
			unitValor = "k"
		} else {
			if unitValor != "k" && unitValor != "m" && unitValor != "b" {
				fmt.Println("El valor del parametro -unit no es valido")
				return
			}
		}
		if !banderaType {
			typeValor = "p"
		} else {
			if typeValor != "p" && typeValor != "e" && typeValor != "l" {
				fmt.Println("El valor del parametro -type no es valido")
				return
			}
		}
		if banderaDelete {
			if deleteValor != "full" {
				fmt.Println("El valor del parametro -delete no es valido")
				return

			}
		}
		var addInt int
		if banderaAdd {
			//Intentar pasar a entero el valor del size a entero
			addInt, err := strconv.Atoi(addValor)
			if err != nil {
				fmt.Println("El valor del parametro -add no es valido")
				return
			}
			if addInt != 0 {
				fmt.Println("El valor del parametro -add no es valido")
				return
			}
		}
		//Imprimir los valores de los parametros
		fmt.Println("Size: ", sizeInt)
		fmt.Println("Driveletter: ", letterValor)
		fmt.Println("Name: ", nameValor)
		fmt.Println("Fit: ", fitValor)
		fmt.Println("Unit: ", unitValor)
		fmt.Println("Type: ", typeValor)
		fmt.Println("Delete: ", deleteValor)
		fmt.Println("Add: ", addInt)
		//Llamar a la funcion para crear la particion
		Filesystem.Fdisk(sizeInt, letterValor, nameValor, fitValor, unitValor, typeValor, deleteValor, addInt)
	}
}

func analizarMkfs(comandoSeparado *[]string) {
	// mkfs -type=full -id=B145 -fs=3fs

	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingresaron los parametros
	var banderaType, banderaId, banderaFs bool
	//Variables para almacenar los valores de los parametros
	var typeValor, idValor, fsValor string
	typeValor = "full"
	fsValor = "2fs"
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		bandera := ObtenerBandera(valor)
		banderaValor := ObtenerBanderaValor(valor)
		if bandera == "-type" {
			banderaType = true
			typeValor = banderaValor
			typeValor = strings.ToLower(typeValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-id" {
			banderaId = true
			idValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-fs" {
			banderaFs = true
			fsValor = banderaValor
			fsValor = strings.ToLower(fsValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("Parametro no reconocido: ", bandera)
		}
	}
	//Obligatorios: -id
	//Verificar si se ingresaron los parametros obligatorios
	if !banderaId {
		fmt.Println("El parametro -id es obligatorio")
		return
	} else {
		//Verificar si se ingresaron los parametros aceptados
		if banderaType {
			if typeValor != "full" {
				fmt.Println("El valor del parametro -type no es valido")
				return
			}
		}
		if banderaFs {
			if fsValor != "2fs" && fsValor != "3fs" {
				fmt.Println("El valor del parametro -fs no es valido")
				return
			}
		}
		//Imprimir los valores de los parametros
		fmt.Println("Type: ", typeValor)
		fmt.Println("Id: ", idValor)
		fmt.Println("Fs: ", fsValor)
		//Llamar a la funcion para formatear la particion
		Filesystem.Mkfs(typeValor, idValor, fsValor)
	}

}

func analizarMount(comandoSeparado *[]string) {
	//mount -driveletter=A -name=Part1 #id=A118
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingresaron los parametros
	var banderaLetter, banderaName bool
	//Variables para almacenar los valores de los parametros
	var letterValor, nameValor string
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		bandera := ObtenerBandera(valor)
		banderaValor := ObtenerBanderaValor(valor)
		if bandera == "-driveletter" {
			banderaLetter = true
			letterValor = banderaValor
			letterValor = strings.ToUpper(letterValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-name" {
			banderaName = true
			nameValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("Parametro no reconocido: ", bandera)
		}
	}
	//Obligatorios: -driveletter, -name
	//Verificar si se ingresaron los parametros obligatorios
	if !banderaLetter {
		fmt.Println("El parametro -driveletter es obligatorio")
		return
	} else if !banderaName {
		fmt.Println("El parametro -name es obligatorio")
		return
	} else {
		//Imprimir los valores de los parametros
		fmt.Println("Driveletter: ", letterValor)
		fmt.Println("Name: ", nameValor)
		//Llamar a la funcion para montar la particion
		Filesystem.MountPartition(letterValor, nameValor)
	}
}

func analizarLogin(comandoSeparado *[]string) {
	//mount -driveletter=A -name=Part1 #id=A118
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingresaron los parametros
	var banderaUser, banderaPassword, banderaId bool
	//Variables para almacenar los valores de los parametros
	var userValor, passwordValor, idValor string
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		bandera := ObtenerBandera(valor)
		banderaValor := ObtenerBanderaValor(valor)
		if bandera == "-user" {
			banderaUser = true
			userValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-pass" {
			banderaPassword = true
			passwordValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-id" {
			banderaId = true
			idValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("Parametro no reconocido: ", bandera)
		}
	}
	//Obligatorios: -user, -pass, -id
	//Verificar si se ingresaron los parametros obligatorios
	if !banderaUser {
		fmt.Println("El parametro -user es obligatorio")
		return
	}
	if !banderaPassword {
		fmt.Println("El parametro -pass es obligatorio")
		return
	}
	if !banderaId {
		fmt.Println("El parametro -id es obligatorio")
		return
	} else {
		//Imprimir los valores de los parametros
		fmt.Println("User: ", userValor)
		fmt.Println("Password: ", passwordValor)
		fmt.Println("Id: ", idValor)
		//Llamar a la funcion para montar la particion
		Filesystem.Login(userValor, passwordValor, idValor)
	}
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
