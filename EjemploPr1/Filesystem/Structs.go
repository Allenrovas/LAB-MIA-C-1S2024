package Filesystem

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
