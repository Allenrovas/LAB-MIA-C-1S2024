package Filesystem

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func ReporteDisk(idValor string, pathValor string) {
	//Abrir el disco A

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		fmt.Println("Error al crear el directorio: ", err)
		return
	}

	//Buscar la particion montada con el ID
	Particion := VerificarParticionMontada(idValor)
	if Particion == -1 {
		fmt.Println("No se encontro la particion montada con el ID: ", idValor)
		return
	}
	MountActual := particionesMontadas[Particion]
	//Abrir el disco
	archivo, err := os.OpenFile("Discos/"+MountActual.LetterValor+".dsk", os.O_RDWR, 0664)
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

	//Quitar la extension al archivo (pdf, etc, )
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

func RepTree(idValor string, pathValor string) {
	//Crear el directorio si no existe
	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		fmt.Println("Error al crear el directorio: ", err)
		return
	}

	//Buscar la particion montada con el ID
	Particion := VerificarParticionMontada(idValor)
	if Particion == -1 {
		fmt.Println("No se encontro la particion montada con el ID: ", idValor)
		return
	}
	MountActual := particionesMontadas[Particion]
	//Abrir el disco
	archivo, err := os.OpenFile("Discos/"+MountActual.LetterValor+".dsk", os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	archivo.Seek(int64(MountActual.Start), 0)
	//Leer el superbloque
	sb := NewSuperBlock()
	err = binary.Read(archivo, binary.LittleEndian, &sb)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	//Buscar el inodo raiz
	raiz := NewInodes()
	archivo.Seek(int64(sb.S_inode_start), 0)
	binary.Read(archivo, binary.LittleEndian, &raiz)
	Dot := "digraph H {\n"
	Dot += "node [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"1\"];\n"
	Dot += "node [shape=plaintext];\n"
	Dot += "graph [bb=\"0,0,352,154\"];\n"
	Dot += "rankdir=LR;\n"
	Dot += RecursivoTree(raiz, sb, archivo, 0)
	Dot += "}"

	//Quitar la extension al archivo (pdf, etc, )
	extension := path.Ext(pathValor)
	//Archivo sin extension
	fileName = strings.TrimSuffix(fileName, extension)
	DotName := dirPath + fileName + ".dot"
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
	//Quitar el punto
	extensionSinPunto := strings.TrimPrefix(extension, ".")
	//Correr con todos los permisos
	cmd := exec.Command("dot", "-T", extensionSinPunto, DotName, "-o", dirPath+fileName+extension)
	fmt.Println("dot", "-T", extensionSinPunto, DotName, "-o", dirPath+fileName+extension)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte tree generado con exito")

}

func RecursivoTree(inodo Inodes, sb SuperBlock, archivo *os.File, numeroInodo int) string {
	Dot := "Inodo" + strconv.Itoa(numeroInodo) + "[label = <\n"
	Dot += "<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
	Dot += "<tr><td bgcolor=\"lightgrey\">Inodo" + strconv.Itoa(numeroInodo) + "</td></tr>\n"
	Dot += "<tr><td>i_uid</td><td>" + strconv.Itoa(int(inodo.I_uid)) + "</td></tr>\n"
	Dot += "<tr><td>i_gid</td><td>" + strconv.Itoa(int(inodo.I_gid)) + "</td></tr>\n"
	Dot += "<tr><td>i_size</td><td>" + strconv.Itoa(int(inodo.I_size)) + "</td></tr>\n"
	Dot += "<tr><td>i_atime</td><td>" + string(inodo.I_atime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_ctime</td><td>" + string(inodo.I_ctime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_mtime</td><td>" + string(inodo.I_mtime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_type</td><td>" + string(inodo.I_type[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_perm</td><td>" + strconv.Itoa(int(inodo.I_perm)) + "</td></tr>\n"
	//Recorrer los bloques
	Contador := 0
	for _, i := range inodo.I_block {
		Dot += "<tr><td>i_block" + strconv.Itoa(Contador+1) + "</td><td port='" + strconv.Itoa(Contador+1) + "'>" + strconv.Itoa(int(i)) + "</td></tr>\n"
		Contador++
	}
	Dot += "</table>>];\n"
	//Recorrer los bloques
	Contador = 0
	for _, i := range inodo.I_block {
		if i != -1 {
			//Leer el bloque
			Dot += "Inodo" + strconv.Itoa(numeroInodo) + ":" + strconv.Itoa(Contador+1) + " -> Bloque" + strconv.Itoa(int(i)) + ":0;\n"
			Dot += "Bloque" + strconv.Itoa(int(i)) + "[label = <\n"
			Dot += "<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
			DesplazamientoBloque := int(sb.S_block_start) + (int(i) * binary.Size(FolderBlock{}))
			carpeta := FolderBlock{}
			archivo.Seek(int64(DesplazamientoBloque), 0)
			binary.Read(archivo, binary.LittleEndian, &carpeta)
			if inodo.I_type == [1]byte{'0'} {
				Dot += "<tr><td colspan=\"2\" port='0'>Bloque" + strconv.Itoa(int(i)) + "</td></tr>\n"
				Contador2 := 0
				for _, j := range carpeta.B_content {
					fmt.Println("Nombre: ", string(j.B_name[:]))
					nam := strings.TrimRight(string(j.B_name[:]), string(rune(0)))

					if Contador2 == 0 {
						nam = "."
					}
					if Contador2 == 1 {
						nam = ".."
					}
					if j.B_inodo == -1 {
						nam = ""
					}
					fmt.Println("Nombre: ", nam)
					Dot += "<tr><td>" + nam + "</td><td port='" + strconv.Itoa(Contador2+1) + "'>" + strconv.Itoa(int(j.B_inodo)) + "</td></tr>\n"
					Contador2++
				}
				Dot += "</table>>];\n"
				Contador2 = 0
				for _, j := range carpeta.B_content {
					if j.B_inodo != -1 {
						if j.B_name[0] != '.' {
							//Leer el inodo
							Dot += "Bloque" + strconv.Itoa(int(i)) + ":" + strconv.Itoa(Contador2+1) + " -> Inodo" + strconv.Itoa(int(j.B_inodo)) + ":0;\n"
							//Buscar el inodo siguiente
							DesplazamientoInodo := int(sb.S_inode_start) + (int(j.B_inodo) * binary.Size(Inodes{}))
							inodoSiguiente := NewInodes()
							archivo.Seek(int64(DesplazamientoInodo), 0)
							binary.Read(archivo, binary.LittleEndian, &inodoSiguiente)
							Dot += RecursivoTree(inodoSiguiente, sb, archivo, int(j.B_inodo))
						}
					}
					Contador2++
				}
			} else {
				file := Fileblock{}
				archivo.Seek(int64(DesplazamientoBloque), 0)
				binary.Read(archivo, binary.LittleEndian, &file)
				Dot += "<tr><td colspan=\"1\" port='0'>Bloque" + strconv.Itoa(int(i)) + "</td></tr>\n"
				Dot += "<tr><td port='1'>" + strings.TrimRight(string(file.B_content[:]), string(rune(0))) + "</td></tr>\n"
				Dot += "</table>>];\n"
			}
		}
		Contador++
	}

	return Dot
}

func ReporteSB(idValor string, pathValor string) {
	//Abrir el disco A
	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		fmt.Println("Error al crear el directorio: ", err)
		return
	}

	//Buscar la particion montada con el ID
	Particion := VerificarParticionMontada(idValor)
	if Particion == -1 {
		fmt.Println("No se encontro la particion montada con el ID: ", idValor)
		return
	}
	MountActual := particionesMontadas[Particion]
	//Abrir el disco
	archivo, err := os.OpenFile("Discos/"+MountActual.LetterValor+".dsk", os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	//Leer el superbloque
	sb := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sb)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}
	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte SuperBlock \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "a0[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">SuperBlock</TD><TD></TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_filesystem_type</TD><TD>" + strconv.Itoa(int(sb.S_filesystem_type)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inodes_count</TD><TD>" + strconv.Itoa(int(sb.S_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_blocks_count</TD><TD>" + strconv.Itoa(int(sb.S_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_free_blocks_count</TD><TD>" + strconv.Itoa(int(sb.S_free_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_free_inodes_count</TD><TD>" + strconv.Itoa(int(sb.S_free_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_mtime</TD><TD>" + string(sb.S_mtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_umtime</TD><TD>" + string(sb.S_umtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_mnt_count</TD><TD>" + strconv.Itoa(int(sb.S_mnt_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_magic</TD><TD>" + strconv.Itoa(int(sb.S_magic)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inode_size</TD><TD>" + strconv.Itoa(int(sb.S_inode_size)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_block_size</TD><TD>" + strconv.Itoa(int(sb.S_block_size)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_first_ino</TD><TD>" + strconv.Itoa(int(sb.S_first_ino)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_first_blo</TD><TD>" + strconv.Itoa(int(sb.S_first_blo)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_bm_inode_start</TD><TD>" + strconv.Itoa(int(sb.S_bm_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_bm_block_start</TD><TD>" + strconv.Itoa(int(sb.S_bm_block_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inode_start</TD><TD>" + strconv.Itoa(int(sb.S_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_block_start</TD><TD>" + strconv.Itoa(int(sb.S_block_start)) + "</TD></TR>\n"
	Dot += "</TABLE>>];\n}"

	//Quitar la extension al archivo (pdf, etc, )
	extension := path.Ext(pathValor)
	//Archivo sin extension
	fileName = strings.TrimSuffix(fileName, extension)
	DotName := dirPath + fileName + ".dot"
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
	//Quitar el punto
	extensionSinPunto := strings.TrimPrefix(extension, ".")
	//Correr con todos los permisos
	cmd := exec.Command("dot", "-T", extensionSinPunto, DotName, "-o", dirPath+fileName+extension)
	fmt.Println("dot", "-T", extensionSinPunto, DotName, "-o", dirPath+fileName+extension)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte sb generado con exito")
}
