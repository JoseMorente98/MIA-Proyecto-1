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
	//"fmt"
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
					if strings.ToLower(unit) == "b" || strings.ToLower(unit) == "k" || strings.ToLower(unit) == "m" {
						unit = strings.ToUpper(unit);
					} else {
						color.Red("╔══════════════════════════════════════════════════╗")
						color.Red("           UNIT solo admite B, K o M D:")
						color.Red("╚══════════════════════════════════════════════════╝")
					}
				} else {
					unit = "K"
				}
				if types != "" {
					if strings.ToLower(types) == "p" || strings.ToLower(types) == "e" || strings.ToLower(types) == "l" {
						types = strings.ToUpper(types);						
					} else {
						color.Red("╔══════════════════════════════════════════════════╗")
						color.Red("           TYPE solo admite P, E o L D:")
						color.Red("╚══════════════════════════════════════════════════╝")
					}
				} else {
					types = "P"
				}
				if fit != "" {
					if strings.ToLower(fit) == "bf" || strings.ToLower(fit) == "ff" || strings.ToLower(fit) == "wf" {
						fit = strings.ToUpper(fit);		
					} else {
						color.Red("╔══════════════════════════════════════════════════╗")
						color.Red("          FIT Solo admite BF, FF o WF D:")
						color.Red("╚══════════════════════════════════════════════════╝")
					}
				} else {
					fit = "WF"
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
	size = sizeUNIT * size - 1;

	readFile(name, path, types, fit, size)
}

/**
 * LECTURA DE ARCHIVO
 */
func readFile(name string, path string, types string, fit string, sizeDisk int64) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close() 
	//VALIDAR QUE EXISTA EL DISCO
	if err != nil { 
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

	//===================================================================
	//                    OTRAS COMPROBACIONES
	//===================================================================
	if strings.ToLower(types) == "p" {
		if (m.Mbr_partition_1 == Modelo.PARTICION{}) {
			if (sizeDisk < m.Mbr_size_disponible) {
				m.Mbr_partition_1 = crearParticion(name, types, fit, sizeDisk, int64(size), "1");
				m.Mbr_size_disponible = m.Mbr_size_disponible - sizeDisk;
			} else {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("        Ya no queda espacio disponible D: ")
				color.Red("╚══════════════════════════════════════════════════╝")
				informacionMBR(strconv.FormatInt(m.Mbr_size, 10), strconv.FormatInt(m.Mbr_size_disponible, 10))
			}
			//CREAR PARTICION 1
		} else if (m.Mbr_partition_2 == Modelo.PARTICION{}) {
			particion := Modelo.PARTICION{};
			copy(particion.Part_name[:], name)
			if (sizeDisk < m.Mbr_size_disponible) {
				if bytes.Compare(particion.Part_name[:], m.Mbr_partition_1.Part_name[:]) == 0 {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("   El nombre de partición ya existe en el disco D:")
					color.Red("╚══════════════════════════════════════════════════╝")
				} else {
					//CREAR PARTICION 1
					m.Mbr_partition_2 = crearParticion(name, types, fit, sizeDisk, m.Mbr_partition_1.Part_end, "2");
					m.Mbr_size_disponible = m.Mbr_size_disponible - sizeDisk;
				}	
			} else {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("        Ya no queda espacio disponible D: ")
				color.Red("╚══════════════════════════════════════════════════╝")
				informacionMBR(strconv.FormatInt(m.Mbr_size, 10), strconv.FormatInt(m.Mbr_size_disponible, 10))
			}	
		} else if (m.Mbr_partition_3 == Modelo.PARTICION{}) {
			particion := Modelo.PARTICION{};
			copy(particion.Part_name[:], name)
			if (sizeDisk < m.Mbr_size_disponible) {
				if (
				bytes.Compare(particion.Part_name[:], m.Mbr_partition_1.Part_name[:]) == 0 ||
				bytes.Compare(particion.Part_name[:], m.Mbr_partition_2.Part_name[:]) == 0) {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("   El nombre de partición ya existe en el disco D:")
					color.Red("╚══════════════════════════════════════════════════╝")
				} else {
					//CREAR PARTICION 3
					m.Mbr_partition_3 = crearParticion(name, types, fit, sizeDisk, m.Mbr_partition_2.Part_end, "1");
					m.Mbr_size_disponible = m.Mbr_size_disponible - sizeDisk;
				}
			} else {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("        Ya no queda espacio disponible D: ")
				color.Red("╚══════════════════════════════════════════════════╝")
				informacionMBR(strconv.FormatInt(m.Mbr_size, 10), strconv.FormatInt(m.Mbr_size_disponible, 10))
			}	
		} else if (m.Mbr_partition_4 == Modelo.PARTICION{}) {
			particion := Modelo.PARTICION{};
			copy(particion.Part_name[:], name)
			if (sizeDisk < m.Mbr_size_disponible) {
			
				if (bytes.Compare(particion.Part_name[:], m.Mbr_partition_1.Part_name[:]) == 0 ||
				bytes.Compare(particion.Part_name[:], m.Mbr_partition_2.Part_name[:]) == 0 ||
				bytes.Compare(particion.Part_name[:], m.Mbr_partition_3.Part_name[:]) == 0) {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("   El nombre de partición ya existe en el disco D:")
					color.Red("╚══════════════════════════════════════════════════╝")
				} else {
					//CREAR PARTICION 4
					m.Mbr_partition_4 = crearParticion(name, types, fit, sizeDisk, m.Mbr_partition_3.Part_end, "1");
					m.Mbr_size_disponible = m.Mbr_size_disponible - sizeDisk;
				}
			} else {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("        Ya no queda espacio disponible D: ")
				color.Red("╚══════════════════════════════════════════════════╝")
				informacionMBR(strconv.FormatInt(m.Mbr_size, 10), strconv.FormatInt(m.Mbr_size_disponible, 10))
			}	
		} else {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("      Ya no se puede crear mas particiones D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}
	} else if strings.ToLower(types) == "e" {
		if(m.Mbr_Extendida == 0) {
			if (m.Mbr_partition_1 == Modelo.PARTICION{}) {
				if (sizeDisk < m.Mbr_size_disponible) {
					particion := Modelo.PARTICION{};
					copy(particion.Part_name[:], name)
					copy(particion.Part_fit[:], fit)
					particion.Part_size = sizeDisk;
					particion.Part_size_disponible = sizeDisk;
					particion.Part_status = 'T';
					particion.Part_type = types[0];
					particion.Part_start = int64(size);
					particion.Part_end = particion.Part_start + sizeDisk;
					
					/**
					* LLENADO EBR
					*/
					particionEBR := Modelo.EBR{};
					copy(particionEBR.Part_name[:], name)
					copy(particionEBR.Part_fit[:], fit)
					particionEBR.Part_size = sizeDisk;
					particionEBR.Part_status = 'T';
					particionEBR.Part_start = int64(size);
					particionEBR.Part_end = particion.Part_start + sizeDisk;
					particionEBR.Part_next = -1;
					particionEBR.Part_logica = Modelo.PARTICION_LOGICA{};
					particion.Part_EBR[0] = particionEBR;

					m.Mbr_partition_1 = particion;
					informacionParticion(name, fit, strconv.FormatInt(particion.Part_size, 10), 
					string(particion.Part_status), 
					strconv.FormatInt(particion.Part_start, 10),
					strconv.FormatInt(particion.Part_end, 10), "1", string(particion.Part_type));
					m.Mbr_Extendida = 1;
					m.Mbr_size_disponible = m.Mbr_size_disponible - sizeDisk;
				} else {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("        Ya no queda espacio disponible D: ")
					color.Red("╚══════════════════════════════════════════════════╝")
					informacionMBR(strconv.FormatInt(m.Mbr_size, 10), strconv.FormatInt(m.Mbr_size_disponible, 10))
				}	

			} else if (m.Mbr_partition_2 == Modelo.PARTICION{}) {
				if (sizeDisk < m.Mbr_size_disponible) {

					particion := Modelo.PARTICION{};
					copy(particion.Part_name[:], name)
					copy(particion.Part_fit[:], fit)
					particion.Part_size = sizeDisk;
					particion.Part_size_disponible = sizeDisk;
					particion.Part_status = 'T';
					particion.Part_type = types[0];
					particion.Part_start = m.Mbr_partition_1.Part_end;
					particion.Part_end = particion.Part_start + sizeDisk;
					if bytes.Compare(particion.Part_name[:], m.Mbr_partition_1.Part_name[:]) == 0 {
						color.Red("╔══════════════════════════════════════════════════╗")
						color.Red("   El nombre de partición ya existe en el disco D:")
						color.Red("╚══════════════════════════════════════════════════╝")
					} else {
						
						/**
						* LLENADO EBR
						*/
						particionEBR := Modelo.EBR{};
						copy(particionEBR.Part_name[:], name)
						copy(particionEBR.Part_fit[:], fit)
						particionEBR.Part_size = sizeDisk;
						particionEBR.Part_status = 'T';
						particionEBR.Part_start = int64(size);
						particionEBR.Part_end = particion.Part_start + sizeDisk;
						particionEBR.Part_next = -1;
						particionEBR.Part_logica = Modelo.PARTICION_LOGICA{};
						particion.Part_EBR[0] = particionEBR;

						m.Mbr_partition_2 = particion;
						informacionParticion(name, fit, strconv.FormatInt(particion.Part_size, 10), 
						string(particion.Part_status), 
						strconv.FormatInt(particion.Part_start, 10),
						strconv.FormatInt(particion.Part_end, 10), "2", string(particion.Part_type));
					}	
					m.Mbr_Extendida = 1;
					m.Mbr_size_disponible = m.Mbr_size_disponible - sizeDisk;
				} else {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("        Ya no queda espacio disponible D: ")
					color.Red("╚══════════════════════════════════════════════════╝")
					informacionMBR(strconv.FormatInt(m.Mbr_size, 10), strconv.FormatInt(m.Mbr_size_disponible, 10))
				}
			} else if (m.Mbr_partition_3 == Modelo.PARTICION{}) {
				if (sizeDisk < m.Mbr_size_disponible) {

					particion := Modelo.PARTICION{};
					copy(particion.Part_name[:], name)
					copy(particion.Part_fit[:], fit)
					particion.Part_size = sizeDisk;
					particion.Part_size_disponible = sizeDisk;
					particion.Part_status = 'T';
					particion.Part_type = types[0];
					particion.Part_start = m.Mbr_partition_2.Part_end;
					particion.Part_end = particion.Part_start + sizeDisk;
					if (
					bytes.Compare(particion.Part_name[:], m.Mbr_partition_1.Part_name[:]) == 0 ||
					bytes.Compare(particion.Part_name[:], m.Mbr_partition_2.Part_name[:]) == 0) {
						color.Red("╔══════════════════════════════════════════════════╗")
						color.Red("   El nombre de partición ya existe en el disco D:")
						color.Red("╚══════════════════════════════════════════════════╝")
					} else {
						
						/**
						* LLENADO EBR
						*/
						particionEBR := Modelo.EBR{};
						copy(particionEBR.Part_name[:], name)
						copy(particionEBR.Part_fit[:], fit)
						particionEBR.Part_size = sizeDisk;
						particionEBR.Part_status = 'T';
						particionEBR.Part_start = int64(size);
						particionEBR.Part_end = particion.Part_start + sizeDisk;
						particionEBR.Part_next = -1;
						particionEBR.Part_logica = Modelo.PARTICION_LOGICA{};
						particion.Part_EBR[0] = particionEBR;

						m.Mbr_partition_3 = particion;
						informacionParticion(name, fit, strconv.FormatInt(particion.Part_size, 10), 
						string(particion.Part_status), 
						strconv.FormatInt(particion.Part_start, 10),
						strconv.FormatInt(particion.Part_end, 10), "3", string(particion.Part_type));

					}
					m.Mbr_Extendida = 1;
					m.Mbr_size_disponible = m.Mbr_size_disponible - sizeDisk;
				} else {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("        Ya no queda espacio disponible D: ")
					color.Red("╚══════════════════════════════════════════════════╝")
					informacionMBR(strconv.FormatInt(m.Mbr_size, 10), strconv.FormatInt(m.Mbr_size_disponible, 10))
				}
			} else if (m.Mbr_partition_4 == Modelo.PARTICION{}) {
				if (sizeDisk < m.Mbr_size_disponible) {

					particion := Modelo.PARTICION{};
					copy(particion.Part_name[:], name)
					copy(particion.Part_fit[:], fit)
					particion.Part_size = sizeDisk;
					particion.Part_size_disponible = sizeDisk;
					particion.Part_status = 'T';
					particion.Part_type = types[0];
					particion.Part_start = m.Mbr_partition_3.Part_end;
					particion.Part_end = particion.Part_start + sizeDisk;
					
					if (bytes.Compare(particion.Part_name[:], m.Mbr_partition_1.Part_name[:]) == 0 ||
					bytes.Compare(particion.Part_name[:], m.Mbr_partition_2.Part_name[:]) == 0 ||
					bytes.Compare(particion.Part_name[:], m.Mbr_partition_3.Part_name[:]) == 0) {
						color.Red("╔══════════════════════════════════════════════════╗")
						color.Red("   El nombre de partición ya existe en el disco D:")
						color.Red("╚══════════════════════════════════════════════════╝")
					} else {
						
						/**
						* LLENADO EBR
						*/
						particionEBR := Modelo.EBR{};
						copy(particionEBR.Part_name[:], name)
						copy(particionEBR.Part_fit[:], fit)
						particionEBR.Part_size = sizeDisk;
						particionEBR.Part_status = 'T';
						particionEBR.Part_start = int64(size);
						particionEBR.Part_end = particion.Part_start + sizeDisk;
						particionEBR.Part_next = -1;
						particionEBR.Part_logica = Modelo.PARTICION_LOGICA{};
						particion.Part_EBR[0] = particionEBR;

						m.Mbr_partition_4 = particion;
						informacionParticion(name, fit, strconv.FormatInt(particion.Part_size, 10), 
						string(particion.Part_status), 
						strconv.FormatInt(particion.Part_start, 10),
						strconv.FormatInt(particion.Part_end, 10), "4", string(particion.Part_type));
					}
					m.Mbr_Extendida = 1;
					m.Mbr_size_disponible = m.Mbr_size_disponible - sizeDisk;
				} else {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("        Ya no queda espacio disponible D: ")
					color.Red("╚══════════════════════════════════════════════════╝")
					informacionMBR(strconv.FormatInt(m.Mbr_size, 10), strconv.FormatInt(m.Mbr_size_disponible, 10))
				}
			} else {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("      Ya no se puede crear mas particiones D:")
				color.Red("╚══════════════════════════════════════════════════╝")
				return;
			}
		} else {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("    Solo puede crear una partición extendida D:")
			color.Red("╚══════════════════════════════════════════════════╝")
			m.Mbr_Extendida = 0;
			return;
		}		
	} else if strings.ToLower(types) == "l" {
		if(m.Mbr_Extendida == 1) {
			var nombre [16]byte;
			copy(nombre[:], name)

			if (m.Mbr_partition_1.Part_type == 'E') {
				if (bytes.Compare(nombre[:], m.Mbr_partition_1.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_2.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_3.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_4.Part_name[:]) == 0) {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("   El nombre de partición ya existe en el disco D:")
					color.Red("╚══════════════════════════════════════════════════╝")
					return;
				} 
				for i := 0; i < 50; i++ {
					//color.Red(string(m.Mbr_partition_1.Part_EBR[i].Part_logica.Part_name[:]));
					//color.Red(string(nombre[:]));
					if bytes.Compare(nombre[:], m.Mbr_partition_1.Part_EBR[i].Part_logica.Part_name[:]) == 0 {
						color.Cyan("╔══════════════════════════════════════════════════╗")
						color.Cyan("   El nombre de partición ya existe en el disco D:")
						color.Cyan("╚══════════════════════════════════════════════════╝")
						return;
					}
				}

				for i := 0; i < 50; i++ {
					if(i == 0) {
						if(m.Mbr_partition_1.Part_EBR[0].Part_logica == Modelo.PARTICION_LOGICA{}) {
							//CREAR PARTICION LOGICA
							if (sizeDisk < m.Mbr_partition_1.Part_size_disponible) {
								m.Mbr_partition_1.Part_EBR[0].Part_logica = crearParticionLogica(name, types, fit, sizeDisk, m.Mbr_partition_1.Part_EBR[0].Part_start);
								m.Mbr_partition_1.Part_size_disponible = m.Mbr_partition_1.Part_size_disponible - sizeDisk;
								break;
							} else {
								color.Red("╔══════════════════════════════════════════════════╗")
								color.Red("        Ya no queda espacio disponible D: ")
								color.Red("╚══════════════════════════════════════════════════╝")
								informacionEXT(strconv.FormatInt(m.Mbr_partition_1.Part_size, 10), strconv.FormatInt(m.Mbr_partition_1.Part_size_disponible, 10))
								break;
							}
						}
					} else {
						if(m.Mbr_partition_1.Part_EBR[i] == Modelo.EBR{}) {
							if (sizeDisk < m.Mbr_partition_1.Part_size_disponible) {
								m.Mbr_partition_1.Part_EBR[i] = crearParticionEBR(name, types, fit, sizeDisk, m.Mbr_partition_1.Part_EBR[i-1].Part_end);
								m.Mbr_partition_1.Part_EBR[i].Part_logica = crearParticionLogica(name, types, fit, sizeDisk, m.Mbr_partition_1.Part_EBR[0].Part_end);
								m.Mbr_partition_1.Part_size_disponible = m.Mbr_partition_1.Part_size_disponible - sizeDisk;
								break;
							} else {
								color.Red("╔══════════════════════════════════════════════════╗")
								color.Red("        Ya no queda espacio disponible D: ")
								color.Red("╚══════════════════════════════════════════════════╝")
								informacionEXT(strconv.FormatInt(m.Mbr_partition_1.Part_size, 10), strconv.FormatInt(m.Mbr_partition_1.Part_size_disponible, 10))
								break;
							}
						}
					}
				}
			} else if (m.Mbr_partition_2.Part_type == 'E') {
				if (bytes.Compare(nombre[:], m.Mbr_partition_1.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_2.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_3.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_4.Part_name[:]) == 0) {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("   El nombre de partición ya existe en el disco D:")
					color.Red("╚══════════════════════════════════════════════════╝")
					return;
				} 
				for i := 0; i < 50; i++ {
					if bytes.Compare(nombre[:], m.Mbr_partition_2.Part_EBR[i].Part_logica.Part_name[:]) == 0 {
						color.Cyan("╔══════════════════════════════════════════════════╗")
						color.Cyan("   El nombre de partición ya existe en el disco D:")
						color.Cyan("╚══════════════════════════════════════════════════╝")
						return;
					}
				}

				for i := 0; i < 50; i++ {
					if(i == 0) {
						if(m.Mbr_partition_2.Part_EBR[0].Part_logica == Modelo.PARTICION_LOGICA{}) {
							//CREAR PARTICION LOGICA
							if (sizeDisk < m.Mbr_partition_2.Part_size_disponible) {
								m.Mbr_partition_2.Part_EBR[0].Part_logica = crearParticionLogica(name, types, fit, sizeDisk, m.Mbr_partition_2.Part_EBR[0].Part_start);
								m.Mbr_partition_2.Part_size_disponible = m.Mbr_partition_2.Part_size_disponible - sizeDisk;
								break;
							} else {
								color.Red("╔══════════════════════════════════════════════════╗")
								color.Red("        Ya no queda espacio disponible D: ")
								color.Red("╚══════════════════════════════════════════════════╝")
								informacionEXT(strconv.FormatInt(m.Mbr_partition_2.Part_size, 10), strconv.FormatInt(m.Mbr_partition_2.Part_size_disponible, 10))
								break;
							}
						}
					} else {
						if(m.Mbr_partition_2.Part_EBR[i] == Modelo.EBR{}) {
							if (sizeDisk < m.Mbr_partition_2.Part_size_disponible) {
								m.Mbr_partition_2.Part_EBR[i] = crearParticionEBR(name, types, fit, sizeDisk, m.Mbr_partition_2.Part_EBR[i-1].Part_end);
								m.Mbr_partition_2.Part_EBR[i].Part_logica = crearParticionLogica(name, types, fit, sizeDisk, m.Mbr_partition_2.Part_EBR[0].Part_end);
								m.Mbr_partition_2.Part_size_disponible = m.Mbr_partition_2.Part_size_disponible - sizeDisk;
								break;
							} else {
								color.Red("╔══════════════════════════════════════════════════╗")
								color.Red("        Ya no queda espacio disponible D: ")
								color.Red("╚══════════════════════════════════════════════════╝")
								informacionEXT(strconv.FormatInt(m.Mbr_partition_2.Part_size, 10), strconv.FormatInt(m.Mbr_partition_2.Part_size_disponible, 10))
								break;
							}
						}
					}
				}
			} else if (m.Mbr_partition_3.Part_type == 'E') {
				if (bytes.Compare(nombre[:], m.Mbr_partition_1.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_2.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_3.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_4.Part_name[:]) == 0) {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("   El nombre de partición ya existe en el disco D:")
					color.Red("╚══════════════════════════════════════════════════╝")
					return;
				} 
				for i := 0; i < 50; i++ {
					//color.Red(string(m.Mbr_partition_3.Part_EBR[i].Part_logica.Part_name[:]));
					//color.Red(string(nombre[:]));
					if bytes.Compare(nombre[:], m.Mbr_partition_3.Part_EBR[i].Part_logica.Part_name[:]) == 0 {
						color.Cyan("╔══════════════════════════════════════════════════╗")
						color.Cyan("   El nombre de partición ya existe en el disco D:")
						color.Cyan("╚══════════════════════════════════════════════════╝")
						return;
					}
				}

				for i := 0; i < 50; i++ {
					if(i == 0) {
						if(m.Mbr_partition_3.Part_EBR[0].Part_logica == Modelo.PARTICION_LOGICA{}) {
							//CREAR PARTICION LOGICA
							if (sizeDisk < m.Mbr_partition_3.Part_size_disponible) {
								m.Mbr_partition_3.Part_EBR[0].Part_logica = crearParticionLogica(name, types, fit, sizeDisk, m.Mbr_partition_3.Part_EBR[0].Part_start);
								m.Mbr_partition_3.Part_size_disponible = m.Mbr_partition_3.Part_size_disponible - sizeDisk;
								break;
							} else {
								color.Red("╔══════════════════════════════════════════════════╗")
								color.Red("        Ya no queda espacio disponible D: ")
								color.Red("╚══════════════════════════════════════════════════╝")
								informacionEXT(strconv.FormatInt(m.Mbr_partition_3.Part_size, 10), strconv.FormatInt(m.Mbr_partition_3.Part_size_disponible, 10))
								break;
							}
						}
					} else {
						if(m.Mbr_partition_3.Part_EBR[i] == Modelo.EBR{}) {
							if (sizeDisk < m.Mbr_partition_3.Part_size_disponible) {
								m.Mbr_partition_3.Part_EBR[i] = crearParticionEBR(name, types, fit, sizeDisk, m.Mbr_partition_3.Part_EBR[i-1].Part_end);
								m.Mbr_partition_3.Part_EBR[i].Part_logica = crearParticionLogica(name, types, fit, sizeDisk, m.Mbr_partition_3.Part_EBR[0].Part_end);
								m.Mbr_partition_3.Part_size_disponible = m.Mbr_partition_3.Part_size_disponible - sizeDisk;
								break;
							} else {
								color.Red("╔══════════════════════════════════════════════════╗")
								color.Red("        Ya no queda espacio disponible D: ")
								color.Red("╚══════════════════════════════════════════════════╝")
								informacionEXT(strconv.FormatInt(m.Mbr_partition_3.Part_size, 10), strconv.FormatInt(m.Mbr_partition_3.Part_size_disponible, 10))
								break;
							}
						}
					}
				}
			} else if (m.Mbr_partition_4.Part_type == 'E') {
				if (bytes.Compare(nombre[:], m.Mbr_partition_1.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_2.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_3.Part_name[:]) == 0 ||
				bytes.Compare(nombre[:], m.Mbr_partition_4.Part_name[:]) == 0) {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("   El nombre de partición ya existe en el disco D:")
					color.Red("╚══════════════════════════════════════════════════╝")
					return;
				} 
				for i := 0; i < 50; i++ {
					//color.Red(string(m.Mbr_partition_4.Part_EBR[i].Part_logica.Part_name[:]));
					//color.Red(string(nombre[:]));
					if bytes.Compare(nombre[:], m.Mbr_partition_4.Part_EBR[i].Part_logica.Part_name[:]) == 0 {
						color.Cyan("╔══════════════════════════════════════════════════╗")
						color.Cyan("   El nombre de partición ya existe en el disco D:")
						color.Cyan("╚══════════════════════════════════════════════════╝")
						return;
					}
				}

				for i := 0; i < 50; i++ {
					if(i == 0) {
						if(m.Mbr_partition_4.Part_EBR[0].Part_logica == Modelo.PARTICION_LOGICA{}) {
							//CREAR PARTICION LOGICA
							if (sizeDisk < m.Mbr_partition_4.Part_size_disponible) {
								m.Mbr_partition_4.Part_EBR[0].Part_logica = crearParticionLogica(name, types, fit, sizeDisk, m.Mbr_partition_4.Part_EBR[0].Part_start);
								m.Mbr_partition_4.Part_size_disponible = m.Mbr_partition_4.Part_size_disponible - sizeDisk;
								break;
							} else {
								color.Red("╔══════════════════════════════════════════════════╗")
								color.Red("        Ya no queda espacio disponible D: ")
								color.Red("╚══════════════════════════════════════════════════╝")
								informacionEXT(strconv.FormatInt(m.Mbr_partition_4.Part_size, 10), strconv.FormatInt(m.Mbr_partition_4.Part_size_disponible, 10))
								break;
							}
						}
					} else {
						if(m.Mbr_partition_4.Part_EBR[i] == Modelo.EBR{}) {
							if (sizeDisk < m.Mbr_partition_4.Part_size_disponible) {
								m.Mbr_partition_4.Part_EBR[i] = crearParticionEBR(name, types, fit, sizeDisk, m.Mbr_partition_4.Part_EBR[i-1].Part_end);
								m.Mbr_partition_4.Part_EBR[i].Part_logica = crearParticionLogica(name, types, fit, sizeDisk, m.Mbr_partition_4.Part_EBR[0].Part_end);
								m.Mbr_partition_4.Part_size_disponible = m.Mbr_partition_4.Part_size_disponible - sizeDisk;
								break;
							} else {
								color.Red("╔══════════════════════════════════════════════════╗")
								color.Red("        Ya no queda espacio disponible D: ")
								color.Red("╚══════════════════════════════════════════════════╝")
								informacionEXT(strconv.FormatInt(m.Mbr_partition_4.Part_size, 10), strconv.FormatInt(m.Mbr_partition_4.Part_size_disponible, 10))
								break;
							}
						}
					}
				}
			}
		} else {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("        Debe crear una partición extendida D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}
		
	}
	//fmt.Println(m)

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
 func informacionMBR(size string, sizeD string) {
	color.Magenta("╔══════════════════════════════════════════════════╗")
	color.Magenta("                       MBR ")
	color.Magenta("  Size: " + size)
	color.Magenta("  Size Disponible: " + sizeD)
	color.Magenta("╚══════════════════════════════════════════════════╝")
}

/**
 * INFORMACIÓN DE PARTICION 
 */
func informacionEXT(size string, sizeD string) {
	color.Magenta("╔══════════════════════════════════════════════════╗")
	color.Magenta("                 PARTICION EXTENDIDA ")
	color.Magenta("  Size: " + size)
	color.Magenta("  Size Disponible: " + sizeD)
	color.Magenta("╚══════════════════════════════════════════════════╝")
}

/**
 * INFORMACIÓN DE PARTICION 
 */
func informacionParticion(name string, fit string, size string, status string, inicio string, fin string, number string, types string) {
	color.Yellow("╔══════════════════════════════════════════════════╗")
	color.Yellow("                    PARTICION " + number)
	color.Yellow("  Nombre: " + name)
	color.Yellow("  Tipo: " + types)
	color.Yellow("  Fit: " + fit)
	color.Yellow("  Status: " + status)
	color.Yellow("  Tamano: " + size)
	color.Yellow("  Inicio: " + inicio)
	color.Yellow("  Fin: " + fin)
	color.Yellow("╚══════════════════════════════════════════════════╝")
}

/**
 * INFORMACIÓN DE PARTICION 
 */
 func informacionParticionLogica(name string, fit string, size string, status string, inicio string, fin string, number string, types string) {
	color.Yellow("╔══════════════════════════════════════════════════╗")
	color.Yellow("               PARTICION LOGICA")
	color.Yellow("  Nombre: " + name)
	color.Yellow("  Tipo: " + types)
	color.Yellow("  Fit: " + fit)
	color.Yellow("  Status: " + status)
	color.Yellow("  Tamano: " + size)
	color.Yellow("  Inicio: " + inicio)
	color.Yellow("  Fin: " + fin)
	color.Yellow("╚══════════════════════════════════════════════════╝")
}

/**
 * PARTICION LOGICA 
 */
func crearParticionLogica(name string, types string, fit string, sizeDisk int64, size int64) Modelo.PARTICION_LOGICA {
	particion_logica := Modelo.PARTICION_LOGICA{}
	copy(particion_logica.Part_name[:], name)
	copy(particion_logica.Part_fit[:], fit)
	particion_logica.Part_type = types[0];
	particion_logica.Part_size = sizeDisk;
	particion_logica.Part_status = 'T';
	particion_logica.Part_start = int64(size);
	particion_logica.Part_end = particion_logica.Part_start + sizeDisk;
	informacionParticionLogica(name, fit, strconv.FormatInt(particion_logica.Part_size, 10), 
		string(particion_logica.Part_status), 
		strconv.FormatInt(particion_logica.Part_start, 10),
		strconv.FormatInt(particion_logica.Part_end, 10), "2", string(particion_logica.Part_type));

	return particion_logica; 
}

/**
 * PARTICION EBR 
 */
func crearParticionEBR(name string, types string, fit string, sizeDisk int64, size int64) Modelo.EBR {
	particionEBR := Modelo.EBR{};
	copy(particionEBR.Part_name[:], name)
	copy(particionEBR.Part_fit[:], fit)
	particionEBR.Part_size = sizeDisk;
	particionEBR.Part_status = 'T';
	particionEBR.Part_start = int64(size);
	particionEBR.Part_end = particionEBR.Part_start + sizeDisk;
	particionEBR.Part_next = -1;
	particionEBR.Part_logica = Modelo.PARTICION_LOGICA{};
	return particionEBR;
}

/**
 * CREAR PARTICION 
 */
 func crearParticion(name string, types string, fit string, sizeDisk int64, size int64, numero string ) Modelo.PARTICION {
	particion := Modelo.PARTICION{};
	copy(particion.Part_name[:], name)
	copy(particion.Part_fit[:], fit)
	particion.Part_size = sizeDisk;
	particion.Part_type = types[0];
	particion.Part_status = 'T';
	particion.Part_start = int64(size);
	particion.Part_end = particion.Part_start + sizeDisk;
	informacionParticion(name, fit, strconv.FormatInt(particion.Part_size, 10), 
	string(particion.Part_status), 
	strconv.FormatInt(particion.Part_start, 10),
	strconv.FormatInt(particion.Part_end, 10), numero, string(particion.Part_type));

	return particion;
}