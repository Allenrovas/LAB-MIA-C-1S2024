package Filesystem

import (
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

// Funcion para crear los discos binarios
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

// Funcion para crear las particiones
func Fdisk(sizeValor int, letterValor string, nameValor string, fitValor string, unitValor string, typeValor string, deleteValor string, addValor int) {
	//Abrir el archivo del disco
	archivo, err := os.OpenFile("Discos/"+letterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()
	//Leer el MBR del disco
	var disk MBR
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		return
	}
	//Verificar si se va a eliminar la particion
	if deleteValor == "0" && addValor == 0 {
		TemporalDesplazamiento := 1 + binary.Size(MBR{})
		var ParticionExtendida Partition
		indiceParticion := 0
		var nombreRepetido, verificarEspacio bool
		if disk.Mbr_partition1.Part_size != 0 {
			if disk.Mbr_partition1.Part_type == [1]byte{'e'} {
				ParticionExtendida = disk.Mbr_partition1
			}
			if strings.Contains(string(disk.Mbr_partition1.Part_name[:]), nameValor) {
				nombreRepetido = true
			}
			TemporalDesplazamiento += int(disk.Mbr_partition1.Part_size) + 1
		} else {
			indiceParticion = 1
			verificarEspacio = true
		}
		if disk.Mbr_partition2.Part_size != 0 {
			if disk.Mbr_partition2.Part_type == [1]byte{'e'} {
				ParticionExtendida = disk.Mbr_partition2
			}
			if strings.Contains(string(disk.Mbr_partition2.Part_name[:]), nameValor) {
				nombreRepetido = true
			}
			TemporalDesplazamiento += int(disk.Mbr_partition2.Part_size) + 1
		} else if !verificarEspacio {
			indiceParticion = 2
			verificarEspacio = true
		}
		if disk.Mbr_partition3.Part_size != 0 {
			if disk.Mbr_partition3.Part_type == [1]byte{'e'} {
				ParticionExtendida = disk.Mbr_partition3
			}
			if strings.Contains(string(disk.Mbr_partition3.Part_name[:]), nameValor) {
				nombreRepetido = true
			}
			TemporalDesplazamiento += int(disk.Mbr_partition3.Part_size) + 1
		} else if !verificarEspacio {
			indiceParticion = 3
			verificarEspacio = true
		}
		if disk.Mbr_partition4.Part_size != 0 {
			if disk.Mbr_partition4.Part_type == [1]byte{'e'} {
				ParticionExtendida = disk.Mbr_partition4
			}
			if strings.Contains(string(disk.Mbr_partition4.Part_name[:]), nameValor) {
				nombreRepetido = true
			}
			TemporalDesplazamiento += int(disk.Mbr_partition4.Part_size) + 1
		} else if !verificarEspacio {
			indiceParticion = 4
			verificarEspacio = true
		}
		//Si el indice sigue siendo 0, entonces no hay espacio
		if indiceParticion == 0 && typeValor != "l" {
			fmt.Println("Error: No hay espacio para crear la particion")
			return
		}
		//Si el nombre ya existe, entonces no se puede crear la particion
		if nombreRepetido {
			fmt.Println("Error: El nombre de la particion ya existe")
			return
		}
		//Si el tipo es extendida y ya existe una extendida entonces no se puede crear
		if typeValor == "e" && ParticionExtendida.Part_type == [1]byte{'e'} {
			fmt.Println("Error: Ya existe una particion extendida")
			return
		}
		//Si es diferente a la logica
		if typeValor != "l" {
			particionNueva := NewPartition()
			particionNueva.Part_status = [1]byte{'1'}
			particionNueva.Part_type = [1]byte{typeValor[0]}
			particionNueva.Part_fit = [1]byte{fitValor[0]}
			particionNueva.Part_start = int32(TemporalDesplazamiento)
			var size int32
			if unitValor == "k" {
				size = int32(sizeValor * 1024)
			} else if unitValor == "m" {
				size = int32(sizeValor * 1024 * 1024)
			} else {
				size = int32(sizeValor)
			}
			particionNueva.Part_size = size
			copy(particionNueva.Part_name[:], nameValor)
			//Verificar si hay espacio para la particion
			if int32(TemporalDesplazamiento)+particionNueva.Part_size+1 > disk.Mbr_tamano {
				fmt.Println("Error: No hay espacio para crear la particion")
				return
			}
			if indiceParticion == 1 {
				disk.Mbr_partition1 = particionNueva
			} else if indiceParticion == 2 {
				disk.Mbr_partition2 = particionNueva
			} else if indiceParticion == 3 {
				disk.Mbr_partition3 = particionNueva
			} else if indiceParticion == 4 {
				disk.Mbr_partition4 = particionNueva
			}
			archivo.Seek(0, 0)
			binary.Write(archivo, binary.LittleEndian, &disk)
			archivo.Close()
			if typeValor == "p" {
				fmt.Println("Particion primaria creada con exito")
			} else {
				fmt.Println("Particion extendida creada con exito")
			}
		} else {
			//Particion logica
			//Verificar si hay particion extendida
			if ParticionExtendida.Part_type != [1]byte{'e'} {
				fmt.Println("Error: No hay particion extendida")
				return
			}
			ebr := NewEBR()
			TemporalDesplazamiento = int(ParticionExtendida.Part_start)
			//Leer el EBR a traves de un for
			for {
				//Intentar leer el EBR
				archivo.Seek(int64(TemporalDesplazamiento), 0)
				binary.Read(archivo, binary.LittleEndian, &ebr)

				if ebr.Part_size != 0 {
					//Comprobacion de nombre
					if strings.Contains(string(ebr.Part_name[:]), nameValor) {
						fmt.Println("Error: El nombre de la particion ya existe")
						return
					}
					//Desplazar al siguiente EBR
					TemporalDesplazamiento += int(ebr.Part_size) + 1 + binary.Size(EBR{})
				}
				if ebr.Part_next == 0 {
					break
				}
			}
			//Crear el nuevo EBR
			var size int
			if unitValor == "k" {
				size = sizeValor * 1024
			} else if unitValor == "m" {
				size = sizeValor * 1024 * 1024
			} else {
				size = sizeValor
			}
			//Verificar si hay espacio para la particion
			if int32(TemporalDesplazamiento)+int32(size)+1 > ParticionExtendida.Part_start+ParticionExtendida.Part_size {
				fmt.Println("Error: No hay espacio para crear la particion")
				return
			}
			//Crear el nuevo EBR
			ebrNueva := NewEBR()
			ebrNueva.Part_mount = [1]byte{'1'}
			ebrNueva.Part_fit = [1]byte{fitValor[0]}
			ebrNueva.Part_start = int32(TemporalDesplazamiento) + 1 + int32(binary.Size(EBR{}))
			ebrNueva.Part_size = int32(size)
			ebrNueva.Part_next = int32(TemporalDesplazamiento) + 1 + int32(binary.Size(EBR{})) + ebrNueva.Part_size
			copy(ebrNueva.Part_name[:], nameValor)
			//Escribir el nuevo EBR
			archivo.Seek(int64(TemporalDesplazamiento), 0)
			binary.Write(archivo, binary.LittleEndian, &ebrNueva)
			archivo.Close()
			fmt.Println("Particion logica creada con exito")
			return
		}
	}
}

// Funcion para hacer mount
func MountPartition(letterValor string, nameValor string) {
	//Abrir el archivo del disco
	archivo, err := os.OpenFile("Discos/"+letterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()
	//Leer el MBR del disco
	var disk MBR
	archivo.Seek(int64(0), 0)
	binary.Read(archivo, binary.LittleEndian, &disk)
	//Verificar si el MBR es valido
	if disk.Mbr_tamano == 0 {
		fmt.Println("Error: El disco no es valido")
		return
	}
	//Verificar si la particion existe
	indiceParticion := 0
	if strings.Contains(string(disk.Mbr_partition1.Part_name[:]), nameValor) {
		indiceParticion = 1
	} else if strings.Contains(string(disk.Mbr_partition2.Part_name[:]), nameValor) {
		indiceParticion = 2
	} else if strings.Contains(string(disk.Mbr_partition3.Part_name[:]), nameValor) {
		indiceParticion = 3
	} else if strings.Contains(string(disk.Mbr_partition4.Part_name[:]), nameValor) {
		indiceParticion = 4
	}
	if indiceParticion != 0 {
		if indiceParticion == 1 {
			particion := disk.Mbr_partition1
			//Verificar si la particion esta montada
			for i := 0; i < len(particionesMontadas); i++ {
				if particionesMontadas[i].LetterValor == letterValor && particionesMontadas[i].Name == nameValor {
					fmt.Println("Error: La particion ya esta montada")
					return
				}
			}
			//Montar la particion
			var particionMontada Mount
			particionMontada.LetterValor = letterValor
			particionMontada.Name = nameValor
			particionMontada.Part_type = particion.Part_type

			//ID:Letra Del Disco + Correlativo Partición + *Últimos dos dígitos del Carne

			//Obtener el correlativo de la particion
			contador := 1
			for i := 0; i < len(particionesMontadas); i++ {
				if particionesMontadas[i].LetterValor == letterValor {
					contador++
				}
			}
			particionMontada.Id = letterValor + strconv.Itoa(contador) + "45"
			particionMontada.Start = particion.Part_start
			particionMontada.Size = particion.Part_size
			//Agregar la particion montada
			particionesMontadas = append(particionesMontadas, particionMontada)
			//Escribir en el MBR y cambiar el status
			disk.Mbr_partition1.Part_status = [1]byte{'1'}
			archivo.Seek(0, 0)
			binary.Write(archivo, binary.LittleEndian, &disk)
			fmt.Println("Particion montada con exito con ID: ", particionMontada.Id)
		} else if indiceParticion == 2 {
			particion := disk.Mbr_partition2
			//Verificar si la particion esta montada
			for i := 0; i < len(particionesMontadas); i++ {
				if particionesMontadas[i].LetterValor == letterValor && particionesMontadas[i].Name == nameValor {
					fmt.Println("Error: La particion ya esta montada")
					return
				}
			}
			//Montar la particion
			var particionMontada Mount
			particionMontada.LetterValor = letterValor
			particionMontada.Name = nameValor
			particionMontada.Part_type = particion.Part_type

			//ID:Letra Del Disco + Correlativo Partición + *Últimos dos dígitos del Carne

			//Obtener el correlativo de la particion
			contador := 1
			for i := 0; i < len(particionesMontadas); i++ {
				if particionesMontadas[i].LetterValor == letterValor {
					contador++
				}
			}
			particionMontada.Id = letterValor + strconv.Itoa(contador) + "45"
			particionMontada.Start = particion.Part_start
			particionMontada.Size = particion.Part_size
			//Agregar la particion montada
			particionesMontadas = append(particionesMontadas, particionMontada)
			//Escribir en el MBR y cambiar el status
			disk.Mbr_partition2.Part_status = [1]byte{'1'}
			archivo.Seek(0, 0)
			binary.Write(archivo, binary.LittleEndian, &disk)
			fmt.Println("Particion montada con exito con ID: ", particionMontada.Id)
		} else if indiceParticion == 3 {
			particion := disk.Mbr_partition3
			//Verificar si la particion esta montada
			for i := 0; i < len(particionesMontadas); i++ {
				if particionesMontadas[i].LetterValor == letterValor && particionesMontadas[i].Name == nameValor {
					fmt.Println("Error: La particion ya esta montada")
					return
				}
			}
			//Montar la particion
			var particionMontada Mount
			particionMontada.LetterValor = letterValor
			particionMontada.Name = nameValor
			particionMontada.Part_type = particion.Part_type

			//ID:Letra Del Disco + Correlativo Partición + *Últimos dos dígitos del Carne

			//Obtener el correlativo de la particion
			contador := 1
			for i := 0; i < len(particionesMontadas); i++ {
				if particionesMontadas[i].LetterValor == letterValor {
					contador++
				}
			}
			particionMontada.Id = letterValor + strconv.Itoa(contador) + "45"
			particionMontada.Start = particion.Part_start
			particionMontada.Size = particion.Part_size
			//Agregar la particion montada
			particionesMontadas = append(particionesMontadas, particionMontada)
			//Escribir en el MBR y cambiar el status
			disk.Mbr_partition3.Part_status = [1]byte{'1'}
			archivo.Seek(0, 0)
			binary.Write(archivo, binary.LittleEndian, &disk)
			fmt.Println("Particion montada con exito con ID: ", particionMontada.Id)
		} else if indiceParticion == 4 {
			particion := disk.Mbr_partition4
			//Verificar si la particion esta montada
			for i := 0; i < len(particionesMontadas); i++ {
				if particionesMontadas[i].LetterValor == letterValor && particionesMontadas[i].Name == nameValor {
					fmt.Println("Error: La particion ya esta montada")
					return
				}
			}
			//Montar la particion
			var particionMontada Mount
			particionMontada.LetterValor = letterValor
			particionMontada.Name = nameValor
			particionMontada.Part_type = particion.Part_type

			//ID:Letra Del Disco + Correlativo Partición + *Últimos dos dígitos del Carne

			//Obtener el correlativo de la particion
			contador := 1
			for i := 0; i < len(particionesMontadas); i++ {
				if particionesMontadas[i].LetterValor == letterValor {
					contador++
				}
			}
			particionMontada.Id = letterValor + strconv.Itoa(contador) + "45"
			particionMontada.Start = particion.Part_start
			particionMontada.Size = particion.Part_size
			//Agregar la particion montada
			particionesMontadas = append(particionesMontadas, particionMontada)
			//Escribir en el MBR y cambiar el status
			disk.Mbr_partition4.Part_status = [1]byte{'1'}
			archivo.Seek(0, 0)
			binary.Write(archivo, binary.LittleEndian, &disk)
			fmt.Println("Particion montada con exito con ID: ", particionMontada.Id)
			return
		}
	} else {
		fmt.Println("Error: La particion no existe")
		return
	}
}
