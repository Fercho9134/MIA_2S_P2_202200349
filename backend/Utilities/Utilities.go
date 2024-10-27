package Utilities

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"proyecto1/Structs"
	"strings"
	"strconv"
)

// Funcion para crear un archivo binario
func CreateFile(name string) error {
	//Se asegura que el archivo existe
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	// Crear archivo
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err CreateFile create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Funcion para abrir un archivo binario ead/write mode
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}

// Funcion para escribir un objecto en un archivo binario
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

// Funcion para leer un objeto de un archivo binario
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}

func DeleteFile (name string) error {
	err := os.Remove(name)
	if err != nil {
		fmt.Println("Err DeleteFile==", err)
		return err
	}
	return nil
}

func GenerateReportMBR(mbr Structs.MRB, ebrs []Structs.EBR, outputPath string, file *os.File)error{
	// Crear la carpeta si no existe
	reportsDir := filepath.Dir(outputPath)
	err := os.MkdirAll(reportsDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error al crear la carpeta de reportes: %v", err)
	}

	// Crear el archivo .dot donde se generará el reporte
	dotFilePath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".dot"
	fileDot, err := os.Create(dotFilePath)
	if err != nil {
		return fmt.Errorf("Error al crear el archivo .dot de reporte: %v", err)
	}
	defer fileDot.Close()

	// Iniciar el contenido del archivo en formato Graphviz (.dot)
	content := "digraph G {\n"
	content += "\tnode [shape=none, margin=0]\n"

	//Primera tabla
	content += "tabla1 [label=<\n"
	content += "<table border=\"0\" cellborder=\"1\" cellspacing=\"0\" cellpadding=\"10\" bgcolor=\"#f7f7f7\">\n"
	content += "<tr>\n"
	content += "<td bgcolor=\"#003366\" colspan=\"2\" align=\"center\">"
	content += "<font color=\"white\"><b>Reporte del MBR</b></font>"
	content += "</td>\n"
	content += "</tr>\n"

	// Subgrafo del MBR Comenzamos con la informacion del MBR
	content += "<tr>\n"
	content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
	content += "<font color=\"white\"><b>MBR Tamaño</b></font>"
	content += "</td>\n"
	content += "<td bgcolor=\"#87cefa\" align=\"left\">"
	content +=	"<font color=\"black\">"
	content += fmt.Sprintf("%d", mbr.MbrSize)
	content += "</font>"
	content += "</td>\n"
	content += "</tr>\n"

	content += "<tr>\n"
	content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
	content += "<font color=\"white\"><b>MBR Fecha de Creación</b></font>"
	content += "</td>\n"
	content += "<td bgcolor=\"#87cefa\" align=\"left\">"
	content +=	"<font color=\"black\">"
	content += string(mbr.CreationDate[:])
	content += "</font>"
	content += "</td>\n"
	content += "</tr>\n"

	content += "<tr>\n"
	content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
	content += "<font color=\"white\"><b>MBR Signature</b></font>"
	content += "</td>\n"
	content += "<td bgcolor=\"#87cefa\" align=\"left\">"
	content +=	"<font color=\"black\">"
	content += fmt.Sprintf("%d", mbr.Signature)
	content += "</font>"
	content += "</td>\n"
	content += "</tr>\n"
	
	// Recorrer las particiones del MBR en orden

	for i := 0; i < 4; i++ {
		part := mbr.Partitions[i]
		if part.Size > 0 { // Si la partición tiene un tamaño válido
			partName := strings.TrimRight(string(part.Name[:]), "\x00") // Limpiar el nombre de la partición

			content += "<tr>\n"
			content += "<td bgcolor=\"#003366\" colspan=\"2\" align=\"center\">"
			content += fmt.Sprintf("<font color=\"white\"><b>Partición %s</b></font>", partName)
			content += "</td>\n"
			content += "</tr>\n"

			//Status
			content += "<tr>\n"
			content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
			content += "<font color=\"white\"><b>Status</b></font>"
			content += "</td>\n"
			content += "<td bgcolor=\"#87cefa\" align=\"left\">"
			content +=	"<font color=\"black\">"
			content += string(part.Status[:])
			content += "</font>"
			content += "</td>\n"
			content += "</tr>\n"

			//Type
			content += "<tr>\n"
			content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
			content += "<font color=\"white\"><b>Type</b></font>"
			content += "</td>\n"
			content += "<td bgcolor=\"#87cefa\" align=\"left\">"
			content +=	"<font color=\"black\">"
			content += string(part.Type[:])
			content += "</font>"
			content += "</td>\n"
			content += "</tr>\n"

			//Fit
			content += "<tr>\n"
			content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
			content += "<font color=\"white\"><b>Fit</b></font>"
			content += "</td>\n"
			content += "<td bgcolor=\"#87cefa\" align=\"left\">"
			content +=	"<font color=\"black\">"
			content += string(part.Fit[:])
			content += "</font>"
			content += "</td>\n"
			content += "</tr>\n"

			//Start
			content += "<tr>\n"
			content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
			content += "<font color=\"white\"><b>Start</b></font>"
			content += "</td>\n"
			content += "<td bgcolor=\"#87cefa\" align=\"left\">"
			content +=	"<font color=\"black\">"
			content += fmt.Sprintf("%d", part.Start)
			content += "</font>"
			content += "</td>\n"
			content += "</tr>\n"

			//Size
			content += "<tr>\n"
			content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
			content += "<font color=\"white\"><b>Size</b></font>"
			content += "</td>\n"
			content += "<td bgcolor=\"#87cefa\" align=\"left\">"
			content +=	"<font color=\"black\">"
			content += fmt.Sprintf("%d", part.Size)
			content += "</font>"
			content += "</td>\n"
			content += "</tr>\n"

			//Name
			content += "<tr>\n"
			content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
			content += "<font color=\"white\"><b>Name</b></font>"
			content += "</td>\n"
			content += "<td bgcolor=\"#87cefa\" align=\"left\">"
			content +=	"<font color=\"black\">"
			content += partName
			content += "</font>"
			content += "</td>\n"
			content += "</tr>\n"
			

			// Si la partición es extendida, leer los EBRs
			if string(part.Type[:]) == "e" {
				// Recolectamos todos los EBRs en orden
				ebrPos := part.Start
				var ebrList []Structs.EBR
				for {
					var ebr Structs.EBR
					err := ReadObject(file, &ebr, int64(ebrPos)) // Asegúrate de que la función ReadObject proviene de Utilities
					if err != nil {
						fmt.Println("Error al leer EBR:", err)
						break
					}
					ebrList = append(ebrList, ebr)

					// Si no hay más EBRs, salir del bucle
					if ebr.PartNext == -1 {
						break
					}

					// Mover a la siguiente posición de EBR
					ebrPos = ebr.PartNext
				}

				// Ahora agregamos los EBRs en orden correcto

				for j, ebr := range ebrList {
					ebrName := strings.TrimRight(string(ebr.PartName[:]), "\x00") // Limpiar el nombre del EBR

					content += "<tr>\n"
					content += "<td bgcolor=\"#003366\" colspan=\"2\" align=\"center\">"
					content += fmt.Sprintf("<font color=\"white\"><b>Particion logica %d</b></font>", j+1)
					content += "</td>\n"
					content += "</tr>\n"

					//part_status
					content += "<tr>\n"
					content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
					content += "<font color=\"white\"><b>Part_status</b></font>"
					content += "</td>\n"
					content += "<td bgcolor=\"#87cefa\" align=\"left\">"
					content +=	"<font color=\"black\">"
					content += "0"
					content += "</font>"
					content += "</td>\n"
					content += "</tr>\n"

					//part_next
					content += "<tr>\n"
					content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
					content += "<font color=\"white\"><b>Part_next</b></font>"
					content += "</td>\n"
					content += "<td bgcolor=\"#87cefa\" align=\"left\">"
					content +=	"<font color=\"black\">"
					content += fmt.Sprintf("%d", ebr.PartNext)
					content += "</font>"
					content += "</td>\n"
					content += "</tr>\n"

					//part_fit
					content += "<tr>\n"
					content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
					content += "<font color=\"white\"><b>Part_fit</b></font>"
					content += "</td>\n"
					content += "<td bgcolor=\"#87cefa\" align=\"left\">"
					content +=	"<font color=\"black\">"
					content += string(ebr.PartFit)
					content += "</font>"
					content += "</td>\n"
					content += "</tr>\n"

					//part_start
					content += "<tr>\n"
					content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
					content += "<font color=\"white\"><b>Part_start</b></font>"
					content += "</td>\n"
					content += "<td bgcolor=\"#87cefa\" align=\"left\">"
					content +=	"<font color=\"black\">"
					content += fmt.Sprintf("%d", ebr.PartStart)
					content += "</font>"
					content += "</td>\n"
					content += "</tr>\n"

					//part_size
					content += "<tr>\n"
					content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
					content += "<font color=\"white\"><b>Part_size</b></font>"
					content += "</td>\n"
					content += "<td bgcolor=\"#87cefa\" align=\"left\">"
					content +=	"<font color=\"black\">"
					content += fmt.Sprintf("%d", ebr.PartSize)
					content += "</font>"
					content += "</td>\n"
					content += "</tr>\n"

					//part_name
					content += "<tr>\n"
					content += "<td bgcolor=\"#1e90ff\" align=\"left\">"
					content += "<font color=\"white\"><b>Part_name</b></font>"
					content += "</td>\n"
					content += "<td bgcolor=\"#87cefa\" align=\"left\">"
					content +=	"<font color=\"black\">"
					content += ebrName
					content += "</font>"
					content += "</td>\n"
					content += "</tr>\n"
					
				}
			}
		}
	}

	content += "</table>\n"

	// Cerrar la tabla 1
	content += ">];\n"

	// Cerrar el archivo .dot
	content += "}\n"

	// Escribir el contenido en el archivo .dot
	_, err = fileDot.WriteString(content)
	if err != nil {
		return fmt.Errorf("Error al escribir en el archivo .dot: %v", err)
	}

	fmt.Println("Reporte MBR generado exitosamente en:", dotFilePath)
	return nil
}

func GenerateReportDisk(mbr Structs.MRB, ebrs []Structs.EBR, outputPath string, file *os.File, totalDiskSize int32, fileName string) error {
    // Crear la carpeta si no existe
    reportsDir := filepath.Dir(outputPath)
    err := os.MkdirAll(reportsDir, os.ModePerm)
    if err != nil {
        return fmt.Errorf("Error al crear la carpeta de reportes: %v", err)
    }

    // Crear el archivo .dot donde se generara el reporte
    dotFilePath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".dot"
    fileDot, err := os.Create(dotFilePath)
    if err != nil {
        return fmt.Errorf("Error al crear el archivo .dot de reporte: %v", err)
    }
    defer fileDot.Close()

    // Iniciar el contenido del archivo en formato Graphviz (.dot)
    content := "digraph G {\n"
    content += "\tnode [shape=none];\n"
    content += "\tgraph [splines=false];\n"
    content += "\tsubgraph cluster_disk {\n"
    content += "\t\tlabel=<<b>" + fileName + "</b>>;\n" // Título en negrita
    content += "\t\tstyle=rounded;\n"
    content += "\t\tcolor=black;\n"

    // Iniciar tabla para las particiones
    content += "\t\ttable [label=<\n\t\t\t<TABLE BORDER=\"1\" CELLBORDER=\"2\" CELLSPACING=\"0\" CELLPADDING=\"15\" BGCOLOR=\"#F7F7F7\">\n" // Cambiar el espaciado y color de fondo
    content += "\t\t\t<TR>\n"
    content += "\t\t\t<TD BGCOLOR=\"#ADD8E6\"><b>MBR (159 bytes)</b></TD>\n" // Color para el MBR

    // Variables para el porcentaje y espacio libre
    var usedSpace int32 = 159 // Tamaño del MBR en bytes
    var freeSpace int32 = totalDiskSize - usedSpace

    for i := 0; i < 4; i++ {
        part := mbr.Partitions[i]
        if part.Size > 0 { // Si la partición tiene un tamaño valido
            percentage := float64(part.Size) / float64(totalDiskSize) * 100
            partName := strings.TrimRight(string(part.Name[:]), "\x00") // Limpiar el nombre de la partición

            if string(part.Type[:]) == "p" { // Partición primaria
                content += fmt.Sprintf("\t\t\t<TD BGCOLOR=\"#90EE90\">Primaria<br/><b>%s</b><br/>%.2f%% del disco</TD>\n", partName, percentage) // Color verde para primaria
                usedSpace += part.Size
            } else if string(part.Type[:]) == "e" { // Partición extendida
                content += "\t\t\t<TD BGCOLOR=\"#FFD700\">\n" // Color dorado para extendida
                content += "\t\t\t\t<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"10\">\n"
                content += ("\t\t\t\t<TR><TD COLSPAN=\"")
				content += strconv.Itoa(len(ebrs) * 2 + 1)
				content += "\"><b>Extendida</b></TD></TR>\n"
				espacioExtendida := part.Size
				espacioAcumuladoExtendida := int32(0)
				espacioLibreExtendida := int32(0)

                // Leer los EBRs y agregar las particiones lógicas
                content += "\t\t\t\t<TR>\n"
                for _, ebr := range ebrs {
                    logicalPercentage := float64(ebr.PartSize) / float64(totalDiskSize) * 100
					logicalPercentage_particion := float64(ebr.PartSize) / float64(espacioExtendida) * 100
                    content += fmt.Sprintf("\t\t\t\t<TD BGCOLOR=\"#FFB6C1\">EBR (32 bytes)</TD>\n<TD BGCOLOR=\"#FFB6C1\">Lógica<br/>%.2f%% del disco<br/>%.2f%% de la particion</TD>\n", logicalPercentage, logicalPercentage_particion) // Color rosado para lógica
                    usedSpace += ebr.PartSize + 32 // Añadir el tamaño de la partición lógica y el EBR
					espacioAcumuladoExtendida += ebr.PartSize + 32
                }
				espacioLibreExtendida = espacioExtendida - espacioAcumuladoExtendida
				usedSpace += espacioLibreExtendida
				percentageExtendida := float64(espacioLibreExtendida) / float64(totalDiskSize) * 100
				percentageExtendida_particion := float64(espacioLibreExtendida) / float64(espacioExtendida) * 100
				content += fmt.Sprintf("\t\t\t\t<TD BGCOLOR=\"#D3D3D3\">Libre<br/>%.2f%% del disco<br/>%.2f%% de la particion</TD>\n", percentageExtendida, percentageExtendida_particion) // Color gris para el espacio libre
                content += "\t\t\t\t</TR>\n"
                content += "\t\t\t\t</TABLE>\n"
                content += "\t\t\t</TD>\n"
            }
        }
    }

    // Recalcular el espacio libre
    freeSpace = totalDiskSize - usedSpace
    freePercentage := float64(freeSpace) / float64(totalDiskSize) * 100

    // Agregar el espacio libre restante
    content += fmt.Sprintf("\t\t\t<TD BGCOLOR=\"#D3D3D3\"><b>Libre</b><br/>%.2f%% del disco</TD>\n", freePercentage) // Color gris para espacio libre
    content += "\t\t\t</TR>\n"
    content += "\t\t\t</TABLE>\n>];\n"
    content += "\t}\n"
    content += "}\n"

    // Escribir el contenido en el archivo .dot
    _, err = fileDot.WriteString(content)
    if err != nil {
        return fmt.Errorf("Error al escribir en el archivo .dot: %v", err)
    }

    fmt.Println("Reporte DISK generado exitosamente en:", dotFilePath)
    return nil
}

func FillWithZeros(file *os.File, start int32, size int32) error {
	// Posiciona el archivo al inicio del área que debe ser llenada
	file.Seek(int64(start), 0)

	// Crear un buffer lleno de ceros
	buffer := make([]byte, size)

	// Escribir los ceros en el archivo
	_, err := file.Write(buffer)
	if err != nil {
		fmt.Println("Error al llenar el espacio con ceros:", err)
		return err
	}

	fmt.Println("Espacio llenado con ceros desde el byte", start, "por", size, "bytes.")
	return nil
}

func VerifyZeros(file *os.File, start int32, size int32) {
	zeros := make([]byte, size)
	_, err := file.ReadAt(zeros, int64(start))
	if err != nil {
		fmt.Println("Error al leer la sección eliminada:", err)
		return
	}

	// Verificar si todos los bytes leídos son ceros
	isZeroFilled := true
	for _, b := range zeros {
		if b != 0 {
			isZeroFilled = false
			break
		}
	}

	if isZeroFilled {
		fmt.Println("La partición eliminada está completamente llena de ceros.")
	} else {
		fmt.Println("Advertencia: La partición eliminada no está completamente llena de ceros.")
	}
}