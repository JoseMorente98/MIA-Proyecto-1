package FDISKControlador

import (
	/*"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"github.com/fatih/color"
	"log"
	"os"
	"strconv"
	"time"
	"math/rand"*/
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
	"log"
	"fmt"
	"bytes"
	"unsafe"
	"encoding/binary"

	Modelo "../../Modelo"
	//Path "../../Utilidades/Path"
)

func FDISK(size string, path string, name string, unit string, types string, fit string, add string, delete string)  {
	if size != "" && path != "" && name != "" {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("         No se encuentra el disco D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		} else {
			s := size
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("      No se puede convertir string a int64 D:")
				color.Red("╚══════════════════════════════════════════════════╝")
			}
			if n >= 1 {
				if unit != "" {
					if strings.ToLower(unit) != "b" || strings.ToLower(unit) != "k" || strings.ToLower(unit) != "m" {
						color.Red("╔══════════════════════════════════════════════════╗")
						color.Red("           UNIT solo admite B, K o M D:")
						color.Red("╚══════════════════════════════════════════════════╝")
					}
				} else {
					unit = "k"
				}
				if types != "" {
					if strings.ToLower(types) != "p" || strings.ToLower(types) != "e" || strings.ToLower(types) != "l" {
						color.Red("╔══════════════════════════════════════════════════╗")
						color.Red("           TYPE solo admite P, E o L D:")
						color.Red("╚══════════════════════════════════════════════════╝")
					}
				} else {
					types = "p"
				}
				if fit != "" {
					if strings.ToLower(fit) != "bf" || strings.ToLower(fit) != "ff" || strings.ToLower(fit) != "wf" {
						color.Red("╔══════════════════════════════════════════════════╗")
						color.Red("          FIT Solo admite BF, FF o WF D:")
						color.Red("╚══════════════════════════════════════════════════╝")
					}
				} else {
					fit = "wf"
				}
				//CORRECTO
				Reescritura(name, unit, path, types, fit, n)
			} else {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("          El tamaño no es permitido D:")
				color.Red("╚══════════════════════════════════════════════════╝")
			}
		}
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("      Verificar parametros obligatorios D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func Reescritura(name string, unit string, path string, types string, fit string, size int64)  {
	var sizeUNIT int64;
	switch strings.ToLower(unit) {
	case "b":
		sizeUNIT = 1;
	case "k":
		sizeUNIT = 1024;
	case "m":
		sizeUNIT = 1024 * 1024;
	default:
		sizeUNIT = 1024;
	}
	size = sizeUNIT * size;

	readFile(name, path, types, fit, size)
}

/**
 * LECTURA DE ARCHIVO
 */
func readFile(name string, path string, types string, fit string, sizeDisk int64) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close() 
	if err != nil { //validar que no sea nulo.
		log.Fatal(err)
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("         No se encuentra el disco D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}

	//DECLARAR MBR Y OBTENER TAMAÑO
	m := Modelo.MBR{}
	var size int = int(unsafe.Sizeof(m))
	
	//LEER CANTIDAD DE BYTES
	data := leerBytes(file, size)
	//DECODIFICACION EN BINARIO
	buffer := bytes.NewBuffer(data)
	
	//GUARDAR VARIABLE EN M
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	fmt.Println(m)

	//===================================================================
	//                    OTRAS COMPROBACIONES
	//===================================================================
	if strings.ToLower(types) == "p" {
		if (m.Mbr_partition_1 == Modelo.PARTICION{}) {
			particion := Modelo.PARTICION{};
			copy(particion.Part_name[:], name)
			copy(particion.Part_fit[:], fit)
			particion.Part_size = sizeDisk;
			particion.Part_status = 'T';
			particion.Part_start = int64(size);
			particion.Part_end = particion.Part_start + sizeDisk;
			m.Mbr_partition_1 = particion;
			informacionParticion(name, fit, strconv.FormatInt(particion.Part_size, 10), 
			string(particion.Part_status), 
			strconv.FormatInt(particion.Part_start, 10),
			strconv.FormatInt(particion.Part_end, 10), "1");
		} else if (m.Mbr_partition_2 == Modelo.PARTICION{}) {
			particion := Modelo.PARTICION{};
			copy(particion.Part_name[:], name)
			copy(particion.Part_fit[:], fit)
			particion.Part_size = sizeDisk;
			particion.Part_status = 'T';
			particion.Part_start = m.Mbr_partition_1.Part_end;
			particion.Part_end = particion.Part_start + sizeDisk;
			if bytes.Compare(particion.Part_name[:], m.Mbr_partition_1.Part_name[:]) == 0 {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("   El nombre de partición ya existe en el disco D:")
				color.Red("╚══════════════════════════════════════════════════╝")
			} else {
				m.Mbr_partition_2 = particion;
				informacionParticion(name, fit, strconv.FormatInt(particion.Part_size, 10), 
				string(particion.Part_status), 
				strconv.FormatInt(particion.Part_start, 10),
				strconv.FormatInt(particion.Part_end, 10), "2");
			}		
		} else if (m.Mbr_partition_3 == Modelo.PARTICION{}) {
			particion := Modelo.PARTICION{};
			copy(particion.Part_name[:], name)
			copy(particion.Part_fit[:], fit)
			particion.Part_size = sizeDisk;
			particion.Part_status = 'T';
			particion.Part_start = m.Mbr_partition_2.Part_end;
			particion.Part_end = particion.Part_start + sizeDisk;
			if (
			bytes.Compare(particion.Part_name[:], m.Mbr_partition_1.Part_name[:]) == 0 ||
			bytes.Compare(particion.Part_name[:], m.Mbr_partition_2.Part_name[:]) == 0) {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("   El nombre de partición ya existe en el disco D:")
				color.Red("╚══════════════════════════════════════════════════╝")
			} else {
				m.Mbr_partition_3 = particion;
				informacionParticion(name, fit, strconv.FormatInt(particion.Part_size, 10), 
				string(particion.Part_status), 
				strconv.FormatInt(particion.Part_start, 10),
				strconv.FormatInt(particion.Part_end, 10), "3");
			}
		} else if (m.Mbr_partition_4 == Modelo.PARTICION{}) {
			particion := Modelo.PARTICION{};
			copy(particion.Part_name[:], name)
			copy(particion.Part_fit[:], fit)
			particion.Part_size = sizeDisk;
			particion.Part_status = 'T';
			particion.Part_start = m.Mbr_partition_3.Part_end;
			particion.Part_end = particion.Part_start + sizeDisk;
			
			if (bytes.Compare(particion.Part_name[:], m.Mbr_partition_1.Part_name[:]) == 0 ||
			bytes.Compare(particion.Part_name[:], m.Mbr_partition_2.Part_name[:]) == 0 ||
			bytes.Compare(particion.Part_name[:], m.Mbr_partition_3.Part_name[:]) == 0) {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("   El nombre de partición ya existe en el disco D:")
				color.Red("╚══════════════════════════════════════════════════╝")
			} else {
				m.Mbr_partition_4 = particion;
				informacionParticion(name, fit, strconv.FormatInt(particion.Part_size, 10), 
				string(particion.Part_status), 
				strconv.FormatInt(particion.Part_start, 10),
				strconv.FormatInt(particion.Part_end, 10), "4");
			}
		} else {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("      Ya no se puede crear mas particiones D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}
	} else {
		
	}
	fmt.Println(m)

	//Se situa en la posicion 0,0 del archivo
	file.Seek(0, 0)
	//Escribe el mbr con particiones en el archivo
	s1 := &m
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escrituraBytes(file, binario3.Bytes())
}

/**
 * LECTURA DE BYTES
 */
func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes
	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

/**
 * ESCRITURA DE BYTES
 */
func escrituraBytes(archivo *os.File, bytes []byte) {
	_, err := archivo.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

/**
 * INFORMACIÓN DE PARTICION 
 */
func informacionParticion(name string, fit string, size string, status string, inicio string, fin string, number string) {
	color.Yellow("╔══════════════════════════════════════════════════╗")
	color.Yellow("                    PARTICION " + number)
	color.Yellow("  Nombre: " + name)
	color.Yellow("  Fit: " + fit)
	color.Yellow("  Status: " + status)
	color.Yellow("  Tamano: " + size)
	color.Yellow("  Inicio: " + inicio)
	color.Yellow("  Fin: " + fin)
	color.Yellow("╚══════════════════════════════════════════════════╝")
}