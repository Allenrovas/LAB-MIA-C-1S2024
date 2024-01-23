package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Profesor struct {
	Tipo        int32 //4 bytes
	Id_profesor int32 //4 bytes
	CUI         [13]byte
	Nombre      [25]byte
	Curso       [25]byte
}

type Estudiante struct {
	Tipo          int32
	Id_estudiante int32
	CUI           [13]byte
	Nombre        [25]byte
	Carnet        [25]byte
}

func main() {
	CrearArchivo()
	Menu()
}

// Ejecuta el Menu
func Menu() {
	var ValorMenu string

	fmt.Println("")
	fmt.Println("Sistema de registro de estudiantes y profesores")
	fmt.Println("")
	fmt.Println("1. Registro de profesores")
	fmt.Println("2. Registro de estudiantes")
	fmt.Println("3. Ver registros")
	fmt.Println("4. Salir")
	fmt.Println("")
	fmt.Println("Por favor seleccione una opcion: ")
	fmt.Scan(&ValorMenu)

	if ValorMenu == "1" {
		RegistroProfesor()
	} else if ValorMenu == "2" {
		RegistroEstudiante()
	} else if ValorMenu == "3" {
		VerRegistros()
	} else if ValorMenu == "4" {
		os.Exit(0)
	} else {
		fmt.Println("Opcion no valida")
	}
	Menu()
}

// Funcion para registrar a un profesor en el archivo
func RegistroProfesor() {
	var id int32
	var cui string
	var nombre string
	var curso string

	//Abrir archivo en modo escritura
	arch, err := os.OpenFile("Registros.bin", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644) //
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arch.Close()

	arch.Seek(0, io.SeekEnd)

	//Crear un profesor
	var profesorNuevo Profesor
	profesorNuevo.Tipo = int32(1)
	fmt.Println("Ingrese el ID del profesor: ")
	fmt.Scan(&id)
	profesorNuevo.Id_profesor = id

	fmt.Println("Ingrese el CUI del profesor: ")
	fmt.Scan(&cui)
	copy(profesorNuevo.CUI[:], cui)

	fmt.Println("Ingrese el nombre del profesor: ")
	fmt.Scan(&nombre)
	copy(profesorNuevo.Nombre[:], nombre)

	fmt.Println("Ingrese el curso del profesor: ")
	fmt.Scan(&curso)
	copy(profesorNuevo.Curso[:], curso)

	//Escribir en el archivo

	// Escribir la estructura completa del profesor
	binary.Write(arch, binary.LittleEndian, &profesorNuevo)
	arch.Close()
	fmt.Println("Profesor registrado exitosamente")

}

// Funcion para registrar a un estudiante en el archivo
func RegistroEstudiante() {
	var id int32
	var cui string
	var nombre string
	var carnet string

	//Abrir archivo en modo escritura
	arch, err := os.OpenFile("Registros.bin", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644) //
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arch.Close()

	arch.Seek(0, io.SeekEnd)

	//Crear un estudiante
	var estudianteNuevo Estudiante
	estudianteNuevo.Tipo = 2
	fmt.Println("Ingrese el ID del estudiante: ")
	fmt.Scan(&id)
	estudianteNuevo.Id_estudiante = id

	fmt.Println("Ingrese el CUI del estudiante: ")
	fmt.Scan(&cui)
	copy(estudianteNuevo.CUI[:], cui)

	fmt.Println("Ingrese el nombre del estudiante: ")
	fmt.Scan(&nombre)
	copy(estudianteNuevo.Nombre[:], nombre)

	fmt.Println("Ingrese el carnet del estudiante: ")
	fmt.Scan(&carnet)
	copy(estudianteNuevo.Carnet[:], carnet)

	//Escribir en el archivo
	err = binary.Write(arch, binary.LittleEndian, &estudianteNuevo)
	if err != nil {
		fmt.Println(err)
		return
	}
	arch.Close()
	fmt.Println("Estudiante registrado exitosamente")
}

// Función para ver los registros
func VerRegistros() {
	fmt.Println("Registros")

	// Abrir archivo en modo lectura
	arch, err := os.OpenFile("Registros.bin", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arch.Close()

	// Leer el archivo con un bucle para leer todos los registros
	for {
		//Leer como profesor y separar con el tipo
		var profesor Profesor
		err = binary.Read(arch, binary.LittleEndian, &profesor)
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		//Si el tipo es 1, es profesor
		fmt.Println("Tipo: ", profesor.Tipo)
		if profesor.Tipo == 1 {
			fmt.Println("Profesor")
			fmt.Println("ID: ", profesor.Id_profesor)
			fmt.Println("CUI: ", string(profesor.CUI[:]))
			fmt.Println("Nombre: ", string(profesor.Nombre[:]))
			fmt.Println("Curso: ", string(profesor.Curso[:]))
			fmt.Println("")
		} else if profesor.Tipo == 2 {
			fmt.Println("Estudiante")
			fmt.Println("ID: ", profesor.Id_profesor)
			fmt.Println("CUI: ", string(profesor.CUI[:]))
			fmt.Println("Nombre: ", string(profesor.Nombre[:]))
			fmt.Println("Carnet: ", string(profesor.Curso[:]))
			fmt.Println("")
		}

	}
}

// Crear un archivo binario
func CrearArchivo() {
	if _, err := os.Stat("Registros.bin"); os.IsNotExist(err) {
		arch, err := os.Create("Registros.bin")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer arch.Close()
	}
}
