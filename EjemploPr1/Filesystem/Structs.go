package Filesystem

import "encoding/binary"

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

type EBR struct {
	Part_mount [1]byte
	Part_fit   [1]byte
	Part_start int32
	Part_size  int32
	Part_next  int32
	Part_name  [16]byte
}

type Inodes struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [19]byte
	I_ctime [19]byte
	I_mtime [19]byte
	I_block [16]int32
	I_type  [1]byte
	I_perm  int32
}

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type FolderBlock struct {
	B_content [4]Content
}

type Fileblock struct {
	B_content [64]byte
}

type SuperBlock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [19]byte
	S_umtime            [19]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_first_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

type ContenidoJournal struct {
	Tipo_operacion [10]byte
	Path           [100]byte
	Contenido      [100]byte
	Time           [19]byte
}

type Journal struct {
	Journal_size int32
	Journal_last int32
	Journal      [50]ContenidoJournal
}

func NewJournal() Journal {
	return Journal{
		Journal_size: 0,
		Journal_last: 1,
		Journal:      [50]ContenidoJournal{},
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

func NewSuperBlock() SuperBlock {
	return SuperBlock{
		S_filesystem_type:   0,
		S_inodes_count:      0,
		S_blocks_count:      0,
		S_free_blocks_count: 0,
		S_free_inodes_count: 0,
		S_mtime:             [19]byte{},
		S_umtime:            [19]byte{},
		S_mnt_count:         0,
		S_magic:             0xEF53,
		S_inode_size:        int32(binary.Size(Inodes{})),
		S_block_size:        int32(binary.Size(FolderBlock{})),
		S_first_ino:         0,
		S_first_blo:         0,
		S_bm_inode_start:    0,
		S_bm_block_start:    0,
		S_inode_start:       0,
		S_block_start:       0,
	}
}

func NewContent() Content {
	return Content{
		B_name:  [12]byte{},
		B_inodo: -1,
	}
}

func NewInodes() Inodes {
	return Inodes{
		I_uid:   -1,
		I_gid:   -1,
		I_size:  -1,
		I_atime: [19]byte{},
		I_ctime: [19]byte{},
		I_mtime: [19]byte{},
		I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'-'},
		I_perm:  -1,
	}
}

func NewPartition() Partition {
	return Partition{
		Part_status: [1]byte{'0'},
		Part_type:   [1]byte{'p'},
		Part_fit:    [1]byte{'w'},
		Part_start:  -1,
		Part_size:   -1,
		Part_name:   [16]byte{'~', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func NewEBR() EBR {
	return EBR{
		Part_mount: [1]byte{'0'},
		Part_fit:   [1]byte{'w'},
		Part_start: -1,
		Part_size:  0,
		Part_next:  -1,
		Part_name:  [16]byte{'~', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}
