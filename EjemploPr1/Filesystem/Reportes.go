package Filesystem

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
)

func ReporteDisk(comandoSeparado *[]string) {
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
	sizeMBR := int(disk.Mbr_tamano)
	libre := int(disk.Mbr_tamano)

	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte Disk \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "node[shape=record, color=lightgrey]a0[label=\"MBR"

	if disk.Mbr_partition1.Part_size != 0 {
		libre -= int(disk.Mbr_partition1.Part_size)
		Dot += "|"
		if disk.Mbr_partition1.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition1.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition1.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			eb := NewEBR()
			Desplazamiento := int(disk.Mbr_partition1.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &eb)
			if err != nil {
				fmt.Println("Error al leer el EBR del disco: ", err)
				return
			}
			if eb.Part_size != 0 {
				Dot += "|{"
				PrimerEbr := true
				for {
					if !PrimerEbr {
						Dot += "|EBR"
					} else {
						PrimerEbr = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					porcentaje := (float64(eb.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libre -= int(eb.Part_size)

					Desplazamiento += int(eb.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					binary.Read(archivo, binary.LittleEndian, &eb)
					if eb.Part_size == 0 {
						break
					}
				}
				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if disk.Mbr_partition2.Part_size != 0 {
		libre -= int(disk.Mbr_partition2.Part_size)
		Dot += "|"
		if disk.Mbr_partition2.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition2.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition2.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			eb := NewEBR()
			Desplazamiento := int(disk.Mbr_partition2.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &eb)
			if err != nil {
				fmt.Println("Error al leer el EBR del disco: ", err)
				return
			}
			if eb.Part_size != 0 {
				Dot += "|{"
				PrimerEbr := true
				for {
					if !PrimerEbr {
						Dot += "|EBR"
					} else {
						PrimerEbr = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					porcentaje := (float64(eb.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libre -= int(eb.Part_size)

					Desplazamiento += int(eb.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					binary.Read(archivo, binary.LittleEndian, &eb)
					if eb.Part_size == 0 {
						break
					}
				}
				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if disk.Mbr_partition3.Part_size != 0 {
		libre -= int(disk.Mbr_partition3.Part_size)
		Dot += "|"
		if disk.Mbr_partition3.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition3.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition3.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			eb := NewEBR()
			Desplazamiento := int(disk.Mbr_partition3.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &eb)
			if err != nil {
				fmt.Println("Error al leer el EBR del disco: ", err)
				return
			}
			if eb.Part_size != 0 {
				Dot += "|{"
				PrimerEbr := true
				for {
					if !PrimerEbr {
						Dot += "|EBR"
					} else {
						PrimerEbr = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					porcentaje := (float64(eb.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libre -= int(eb.Part_size)

					Desplazamiento += int(eb.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					binary.Read(archivo, binary.LittleEndian, &eb)
					if eb.Part_size == 0 {
						break
					}
				}
				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if disk.Mbr_partition4.Part_size != 0 {
		libre -= int(disk.Mbr_partition4.Part_size)
		Dot += "|"
		if disk.Mbr_partition4.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition4.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition4.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			eb := NewEBR()
			Desplazamiento := int(disk.Mbr_partition4.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &eb)
			if err != nil {
				fmt.Println("Error al leer el EBR del disco: ", err)
				return
			}
			if eb.Part_size != 0 {
				Dot += "|{"
				PrimerEbr := true
				for {
					if !PrimerEbr {
						Dot += "|EBR"
					} else {
						PrimerEbr = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					porcentaje := (float64(eb.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libre -= int(eb.Part_size)

					Desplazamiento += int(eb.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					binary.Read(archivo, binary.LittleEndian, &eb)
					if eb.Part_size == 0 {
						break
					}
				}
				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if libre > 0 {
		Dot += "|Libre"
		porcentaje := (float64(libre) * float64(100)) / float64(sizeMBR)
		Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
	}
	Dot += "\"];\n}"

	//Crear el archivo .dot
	DotName := "Reportes/ReporteDisk.dot"
	archivoDot, err := os.Create(DotName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot: ", err)
		return
	}
	defer archivoDot.Close()
	_, err = archivoDot.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot: ", err)
		return
	}
	//Generar la imagen
	cmd := exec.Command("dot", "-T", "png", DotName, "-o", "Reportes/ReporteDisk.png")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte generado con exito")

}
