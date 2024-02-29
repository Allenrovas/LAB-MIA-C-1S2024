package Filesystem

import (
	"encoding/binary"
	"math"
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

}
