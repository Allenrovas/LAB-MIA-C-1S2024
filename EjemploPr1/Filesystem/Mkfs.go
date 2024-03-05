package Filesystem

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"time"
)

// Mkfs crea un sistema de archivos en una partición
func Mkfs(typeValor string, idValor string, fsValor string) {
	//Verificar que el id exista en la lista de particiones montadas
	indice := VerificarParticionMontada(idValor)
	if indice == -1 {
		println("La partición no está montada")
		return
	}

	MountActual := particionesMontadas[indice]

	if MountActual.Part_type != [1]byte{'p'} {
		println("La partición no es primaria")
		return
	}

	//Cantidad de estructuras que caben en la partición
	var n int
	if fsValor == "2fs" {
		n = int(math.Floor(float64(int(MountActual.Size)-int(binary.Size(SuperBlock{}))) / float64(4+int(binary.Size(Inodes{}))+3*int(binary.Size(Fileblock{})))))

	} else {
		n = int(math.Floor(float64(int(MountActual.Size)-int(binary.Size(SuperBlock{}))) / float64(4+int(binary.Size(Inodes{}))+3*int(binary.Size(FolderBlock{}))+binary.Size(Journal{}))))

	}

	//Crear el superbloque
	sb := NewSuperBlock()
	sb.S_inodes_count = int32(n)
	sb.S_blocks_count = int32(n * 3)
	sb.S_free_blocks_count = int32(n * 3)
	sb.S_free_inodes_count = int32(n)
	fechaActual := time.Now()
	fecha := fechaActual.Format("2006-01-02 15:04:05")
	copy(sb.S_mtime[:], fecha)
	copy(sb.S_umtime[:], fecha)
	sb.S_mnt_count = 1
	if fsValor == "2fs" {
		Crear2FS(sb, MountActual, n)
	} else {
		Crear3FS(sb, MountActual, n)
	}

}

// Crear2FS crea un sistema de archivos 2fs
func Crear2FS(sb SuperBlock, MountActual Mount, n int) {

	sb.S_filesystem_type = 2
	sb.S_bm_inode_start = int32(MountActual.Start) + int32(binary.Size(SuperBlock{}))
	sb.S_bm_block_start = sb.S_bm_inode_start + int32(n)
	sb.S_inode_start = sb.S_bm_block_start + int32(3*n)
	sb.S_block_start = sb.S_inode_start + int32(n*int(binary.Size(Inodes{})))
	//Crear el bloque 0, inodo 0 y el usuario root
	sb.S_free_blocks_count--
	sb.S_free_inodes_count--
	sb.S_free_blocks_count--
	sb.S_free_inodes_count--

	//Creación del super bloque
	//Abrir el archivo

	file, err := os.OpenFile("Discos/"+MountActual.LetterValor+".dsk", os.O_WRONLY, 0777)
	if err != nil {
		println("Error al abrir el archivo")
		return
	}
	defer file.Close()

	file.Seek(int64(MountActual.Start), 0)
	binary.Write(file, binary.LittleEndian, &sb)

	//Crear el bitmap de inodos
	var llenar byte = 0
	file.Seek(int64(sb.S_bm_inode_start), 0)
	for i := 0; i < n; i++ {
		binary.Write(file, binary.LittleEndian, &llenar)
	}

	//Crear el bitmap de bloques
	file.Seek(int64(sb.S_bm_block_start), 0)
	for i := 0; i < n*3; i++ {
		binary.Write(file, binary.LittleEndian, &llenar)
	}

	//Crear el inodo 0
	inodo0 := NewInodes()

	//Crear el bloque 0
	var bloque0 Fileblock

	//Formatear inodos
	file.Seek(int64(sb.S_inode_start), 0)
	for i := 0; i < n; i++ {
		binary.Write(file, binary.LittleEndian, &inodo0)
	}

	//Formatear bloques
	file.Seek(int64(sb.S_block_start), 0)
	for i := 0; i < n*3; i++ {
		binary.Write(file, binary.LittleEndian, &bloque0)
	}

	//Crear el directorio raíz
	//Crear el inodo
	inodo0.I_uid = 1
	inodo0.I_gid = 1
	fechaActual := time.Now()
	fecha := fechaActual.Format("2006-01-02 15:04:05")
	copy(inodo0.I_atime[:], fecha)
	copy(inodo0.I_ctime[:], fecha)
	copy(inodo0.I_mtime[:], fecha)
	inodo0.I_type = [1]byte{'0'}
	inodo0.I_perm = 664
	inodo0.I_block[0] = 0

	//Crear el bloque carpeta

	var bloqueCarpeta FolderBlock
	bloqueCarpeta.B_content[0].B_inodo = 0
	copy(bloqueCarpeta.B_content[0].B_name[:], ".")
	bloqueCarpeta.B_content[1].B_inodo = 0
	copy(bloqueCarpeta.B_content[1].B_name[:], "..")
	bloqueCarpeta.B_content[2].B_inodo = 1
	copy(bloqueCarpeta.B_content[2].B_name[:], "users.txt")
	bloqueCarpeta.B_content[3].B_inodo = -1

	data := "1,G,root\n1,U,root,root,123\n"

	//Escribir el inodo y el bloque en el archivo

	inodo1 := NewInodes()
	inodo1.I_uid = 1
	inodo1.I_gid = 1
	fechaActual = time.Now()
	fecha = fechaActual.Format("2006-01-02 15:04:05")
	copy(inodo1.I_atime[:], fecha)
	copy(inodo1.I_ctime[:], fecha)
	copy(inodo1.I_mtime[:], fecha)
	inodo1.I_type = [1]byte{'1'}
	inodo1.I_perm = 664
	inodo1.I_block[0] = 1
	inodo1.I_size = int32(len(data)) + int32(binary.Size(Fileblock{}))

	inodo0.I_size = inodo1.I_size + int32(binary.Size(FolderBlock{})) + int32(binary.Size(FolderBlock{}))

	var bloqueArchivo Fileblock
	copy(bloqueArchivo.B_content[:], data)

	//Escribir el inodo en el archivo
	file.Seek(int64(sb.S_bm_inode_start), 0)
	var bit byte = 1
	binary.Write(file, binary.LittleEndian, &bit)
	binary.Write(file, binary.LittleEndian, &bit)

	file.Seek(int64(sb.S_bm_block_start), 0)
	binary.Write(file, binary.LittleEndian, &bit)
	binary.Write(file, binary.LittleEndian, &bit)

	file.Seek(int64(sb.S_inode_start), 0)
	binary.Write(file, binary.LittleEndian, &inodo0)
	binary.Write(file, binary.LittleEndian, &inodo1)

	file.Seek(int64(sb.S_block_start), 0)
	binary.Write(file, binary.LittleEndian, &bloqueCarpeta)
	binary.Write(file, binary.LittleEndian, &bloqueArchivo)

	fmt.Println("Sistema de archivos 2FS creado con éxito en el disco: " + MountActual.Id)

}

// Crear3FS crea un sistema de archivos 3fs
func Crear3FS(sb SuperBlock, MountActual Mount, n int) {
	sb.S_filesystem_type = 3
	sb.S_bm_inode_start = int32(MountActual.Start) + int32(binary.Size(SuperBlock{})) + int32(binary.Size(Journal{}))
	sb.S_bm_block_start = sb.S_bm_inode_start + int32(n)
	sb.S_inode_start = sb.S_bm_block_start + int32(3*n)
	sb.S_block_start = sb.S_inode_start + int32(n*int(binary.Size(Inodes{})))
	//Crear el bloque 0, inodo 0 y el usuario root
	sb.S_free_blocks_count--
	sb.S_free_inodes_count--
	sb.S_free_blocks_count--
	sb.S_free_inodes_count--

	//Creacion Journaling
	var journal Journal

	inicioJournal := "mkdir"
	pathJournal := "/"
	contenidoJournaling := "-"

	fechaActual := time.Now()
	fecha := fechaActual.Format("2006-01-02 15:04:05")

	copy(journal.Journal[0].Tipo_operacion[:], inicioJournal)
	copy(journal.Journal[0].Path[:], pathJournal)
	copy(journal.Journal[0].Contenido[:], contenidoJournaling)
	copy(journal.Journal[0].Time[:], fecha)

	inicioJournal = "mkfile"
	pathJournal = "/users.txt"
	contenidoJournaling = "1,G,root\n1,U,root,root,123\n"

	journal.Journal_size = 2
	journal.Journal_last = 1

	copy(journal.Journal[1].Tipo_operacion[:], inicioJournal)
	copy(journal.Journal[1].Path[:], pathJournal)
	copy(journal.Journal[1].Contenido[:], contenidoJournaling)
	copy(journal.Journal[1].Time[:], fecha)

	//Creación del superbloque
	//Abrir el archivo
	file, err := os.OpenFile("Discos/"+MountActual.LetterValor, os.O_WRONLY, 0777)
	if err != nil {
		println("Error al abrir el archivo")
		return
	}
	defer file.Close()

	file.Seek(int64(MountActual.Start), 0)
	file.Seek(int64(MountActual.Start), 0)
	binary.Write(file, binary.LittleEndian, &sb)
	binary.Write(file, binary.LittleEndian, &journal)

	//Crear el bitmap de inodos
	var llenar byte = 0
	file.Seek(int64(sb.S_bm_inode_start), 0)
	for i := 0; i < n; i++ {
		binary.Write(file, binary.LittleEndian, &llenar)
	}

	//Crear el bitmap de bloques
	file.Seek(int64(sb.S_bm_block_start), 0)
	for i := 0; i < n*3; i++ {
		binary.Write(file, binary.LittleEndian, &llenar)
	}

	inodo0 := NewInodes()
	var bloque0 Fileblock

	//Formatear inodos
	file.Seek(int64(sb.S_inode_start), 0)
	for i := 0; i < n; i++ {
		binary.Write(file, binary.LittleEndian, &inodo0)
	}

	//Formatear bloques
	file.Seek(int64(sb.S_block_start), 0)
	for i := 0; i < n*3; i++ {
		binary.Write(file, binary.LittleEndian, &bloque0)
	}

	//Crear el inodo 0
	inodo0.I_uid = 1
	inodo0.I_gid = 1
	fechaActual = time.Now()
	fecha = fechaActual.Format("2006-01-02 15:04:05")
	copy(inodo0.I_atime[:], fecha)
	copy(inodo0.I_ctime[:], fecha)
	copy(inodo0.I_mtime[:], fecha)
	inodo0.I_type = [1]byte{'0'}
	inodo0.I_perm = 664
	inodo0.I_block[0] = 0

	//Crear el bloque carpeta

	var bloqueCarpeta FolderBlock
	bloqueCarpeta.B_content[0].B_inodo = 0
	copy(bloqueCarpeta.B_content[0].B_name[:], ".")
	bloqueCarpeta.B_content[1].B_inodo = 0
	copy(bloqueCarpeta.B_content[1].B_name[:], "..")
	bloqueCarpeta.B_content[2].B_inodo = 1
	copy(bloqueCarpeta.B_content[2].B_name[:], "users.txt")
	bloqueCarpeta.B_content[3].B_inodo = -1

	data := "1,G,root\n1,U,root,root,123\n"

	inodo1 := NewInodes()
	inodo1.I_uid = 1
	inodo1.I_gid = 1
	fechaActual = time.Now()
	fecha = fechaActual.Format("2006-01-02 15:04:05")
	copy(inodo1.I_atime[:], fecha)
	copy(inodo1.I_ctime[:], fecha)
	copy(inodo1.I_mtime[:], fecha)
	inodo1.I_type = [1]byte{'1'}
	inodo1.I_perm = 664
	inodo1.I_block[0] = 1
	inodo1.I_size = int32(len(data)) + int32(binary.Size(Fileblock{}))

	inodo0.I_size = inodo1.I_size + int32(binary.Size(FolderBlock{})) + int32(binary.Size(FolderBlock{}))

	var bloqueArchivo Fileblock
	copy(bloqueArchivo.B_content[:], data)

	//Escribir el inodo en el archivo
	file.Seek(int64(sb.S_bm_inode_start), 0)
	var bit byte = 1
	binary.Write(file, binary.LittleEndian, &bit)
	binary.Write(file, binary.LittleEndian, &bit)

	file.Seek(int64(sb.S_bm_block_start), 0)
	binary.Write(file, binary.LittleEndian, &bit)
	binary.Write(file, binary.LittleEndian, &bit)

	file.Seek(int64(sb.S_inode_start), 0)
	binary.Write(file, binary.LittleEndian, &inodo0)
	binary.Write(file, binary.LittleEndian, &inodo1)

	file.Seek(int64(sb.S_block_start), 0)
	binary.Write(file, binary.LittleEndian, &bloqueCarpeta)
	binary.Write(file, binary.LittleEndian, &bloqueArchivo)

	fmt.Println("Sistema de archivos 3FS creado con éxito en el disco: " + MountActual.Id)

}
