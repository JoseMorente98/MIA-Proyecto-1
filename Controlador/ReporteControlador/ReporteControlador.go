package ReporteControlador

import (
	"strings"
	"github.com/fatih/color"
	"runtime"
	"path"
	"path/filepath"
	"os/exec"
	"strconv"
	"log"
	"os"
	"unsafe"
	"fmt"
	"encoding/binary"
	"bytes"

	Modelo "../../Modelo"
	Ubicacion "../../Utilidades/Path"
	MOUNTControlador "../MOUNTControlador"
)

/**
 * FUNCION REP
 */
func REP(nombre string, ubicacion string, id string, ruta string) {
	if nombre != "" && ubicacion != "" && id != "" {
		reportes(nombre, ubicacion, id, ruta);

	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("       Verificar parametros obligatorios D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

/**
 * VALIDAR PATH
 */
func reportes(nombre string, ubicacion string, id string, ruta string)  {
	var arregloPath = strings.Split(ubicacion, "/")
	var strPath = "";
	for i := 0; i < len(arregloPath) -1 ; i++ {
		color.Green(arregloPath[i])
		strPath += arregloPath[i] + "/"
	}
	color.Cyan(strPath)
	Ubicacion.CrearDirectorio(strPath);
	//executeCommand(ubicacion)

	//VERIFICA QUE ESTE MONTADO EL DISCO
	if (MOUNTControlador.Buscar(id) == true) {
		mount := MOUNTControlador.BuscarMOUNT(id);
		TipoReporte(nombre, ubicacion, id, ruta, mount.Mount_path)
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("     La partición no se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}

}

/**
 * TIPO DE REPORTE
 */
func TipoReporte(nombre string, ubicacionFinal string, id string, ruta string, pathDisco string)  {
	switch strings.ToLower(nombre) {
		case "mbr":
			MBR(pathDisco, ubicacionFinal);
		case "disk":
			DISK(pathDisco, ubicacionFinal);
		case "sb":
			SUPERBOOT(id, ubicacionFinal);
		case "bm_arbdir":
			BM_AVD(id, ubicacionFinal);
		case "bm_detdir":
			BM_DD(id, ubicacionFinal);
		case "bm_inode":
			BM_INODO(id, ubicacionFinal);
		case "bm_block":
			BM_BLOQUE(id, ubicacionFinal);
		case "bitacora":
		case "directorio":
		case "tree_file":
		case "tree_directorio":
		case "tree_complete":
		case "ls":
	}
}

/**
 * EJECUTAR COMANDO
 */
func executeCommand(strFile string, strPath string)  {
	cmd := exec.Command("dot", "-Tpng", strFile, "-o", strPath)
	//color.Cyan(cmd)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("        Error al generar grafico comando D:")
		log.Fatal(err)
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

/**
 * MBR
 */
func MBR(ubicacion string, strPath string)  {
	file, err := os.OpenFile(ubicacion, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close() 
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
	data := lecturaBytes(file, size)
	//DECODIFICACION EN BINARIO
	buffer := bytes.NewBuffer(data)
	
	//GUARDAR VARIABLE EN M
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	var strReporte = "digraph {\n"+
	"tbl [ \n"+
	"shape=plaintext \n"+
	"label=<\n"+

	"<table border='1' cellborder='1'>\n"+
	"<tr><td>Name</td><td>Valor</td></tr>\n"

	//LLENADO MBR
	strReporte += "<tr><td>MBR Tamano</td><td>"+ strconv.FormatInt(m.Mbr_size, 10) +"</td></tr>\n"
	strReporte += "<tr><td>MBR Fecha Creacion</td><td>"+ string(m.Mbr_date[:]) +"</td></tr>\n"
	strReporte += "<tr><td>MBR Disk Asignature</td><td>"+  strconv.FormatInt(m.Mbr_disk_signature, 10) +"</td></tr>\n"

	//REPORTE PARTICION 1
	s1 := BytesToString(m.Mbr_partition_1.Part_name[:])
	s2 := BytesToString(m.Mbr_partition_2.Part_name[:])
	s3 := BytesToString(m.Mbr_partition_3.Part_name[:])
	s4 := BytesToString(m.Mbr_partition_4.Part_name[:])
	if (m.Mbr_partition_1 != Modelo.PARTICION{}) {
		strReporte += "<tr><td>MBR Particion No. 1</td><td></td></tr>\n"
		strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_1.Part_status) +"</td></tr>\n"
		strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_1.Part_type) +"</td></tr>\n"
		strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_1.Part_fit[:]) +"</td></tr>\n"
		strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_1.Part_start, 10) +"</td></tr>\n"
		strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_1.Part_size, 10) +"</td></tr>\n"
		strReporte += "<tr><td>Name</td><td>"+ s1 +"</td></tr>\n"
		if (m.Mbr_partition_1.Part_type == 'E' || m.Mbr_partition_1.Part_type == 'e') {
			for i := 0; i < 25; i++ {
				if (m.Mbr_partition_1.Part_EBR[i] != Modelo.EBR{}) {
					if(m.Mbr_partition_1.Part_EBR[i].Part_logica != Modelo.PARTICION_LOGICA{}) {
						strReporte += "<tr><td>EBR Particion</td><td></td></tr>\n"
						strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_1.Part_EBR[i].Part_logica.Part_status) +"</td></tr>\n"
						strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_1.Part_EBR[i].Part_logica.Part_type) +"</td></tr>\n"
						strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_1.Part_EBR[i].Part_logica.Part_fit[:]) +"</td></tr>\n"
						strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_1.Part_EBR[i].Part_logica.Part_start, 10) +"</td></tr>\n"
						strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_1.Part_EBR[i].Part_logica.Part_size, 10) +"</td></tr>\n"
						strReporte += "<tr><td>Name</td><td>"+ BytesToString(m.Mbr_partition_1.Part_EBR[i].Part_logica.Part_name[:]) +"</td></tr>\n"
					}
				}
			}
		}
	}
	if (m.Mbr_partition_2 != Modelo.PARTICION{}) {
		//REPORTE PARTICION 2
		strReporte += "<tr><td>MBR Particion No. 2</td><td></td></tr>\n"
		strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_2.Part_status) +"</td></tr>\n"
		strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_2.Part_type) +"</td></tr>\n"
		strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_2.Part_fit[:]) +"</td></tr>\n"
		strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_2.Part_start, 10) +"</td></tr>\n"
		strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_2.Part_size, 10) +"</td></tr>\n"
		strReporte += "<tr><td>Name</td><td>"+ s2 +"</td></tr>\n"
		if (m.Mbr_partition_2.Part_type == 'E' || m.Mbr_partition_2.Part_type == 'e') {
			for i := 0; i < 25; i++ {
				if (m.Mbr_partition_2.Part_EBR[i] != Modelo.EBR{}) {
					if(m.Mbr_partition_2.Part_EBR[i].Part_logica != Modelo.PARTICION_LOGICA{}) {
						strReporte += "<tr><td>EBR Particion</td><td></td></tr>\n"
						strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_2.Part_EBR[i].Part_logica.Part_status) +"</td></tr>\n"
						strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_2.Part_EBR[i].Part_logica.Part_type) +"</td></tr>\n"
						strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_2.Part_EBR[i].Part_logica.Part_fit[:]) +"</td></tr>\n"
						strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_2.Part_EBR[i].Part_logica.Part_start, 10) +"</td></tr>\n"
						strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_2.Part_EBR[i].Part_logica.Part_size, 10) +"</td></tr>\n"
						strReporte += "<tr><td>Name</td><td>"+ BytesToString(m.Mbr_partition_2.Part_EBR[i].Part_logica.Part_name[:]) +"</td></tr>\n"
					}
				}
			}
		}
	}
	if (m.Mbr_partition_3 != Modelo.PARTICION{}) {
		//REPORTE PARTICION 3
		strReporte += "<tr><td>MBR Particion No. 3</td><td></td></tr>\n"
		strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_3.Part_status) +"</td></tr>\n"
		strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_3.Part_type) +"</td></tr>\n"
		strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_3.Part_fit[:]) +"</td></tr>\n"
		strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_3.Part_start, 10) +"</td></tr>\n"
		strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_3.Part_size, 10) +"</td></tr>\n"
		strReporte += "<tr><td>Name</td><td>"+ s3 +"</td></tr>\n"
		if (m.Mbr_partition_3.Part_type == 'E' || m.Mbr_partition_3.Part_type == 'e') {
			for i := 0; i < 25; i++ {
				if (m.Mbr_partition_3.Part_EBR[i] != Modelo.EBR{}) {
					if(m.Mbr_partition_3.Part_EBR[i].Part_logica != Modelo.PARTICION_LOGICA{}) {
						strReporte += "<tr><td>EBR Particion</td><td></td></tr>\n"
						strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_3.Part_EBR[i].Part_logica.Part_status) +"</td></tr>\n"
						strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_3.Part_EBR[i].Part_logica.Part_type) +"</td></tr>\n"
						strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_3.Part_EBR[i].Part_logica.Part_fit[:]) +"</td></tr>\n"
						strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_3.Part_EBR[i].Part_logica.Part_start, 10) +"</td></tr>\n"
						strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_3.Part_EBR[i].Part_logica.Part_size, 10) +"</td></tr>\n"
						strReporte += "<tr><td>Name</td><td>"+ BytesToString(m.Mbr_partition_3.Part_EBR[i].Part_logica.Part_name[:]) +"</td></tr>\n"
					}
				}
			}
		}
	}
	if (m.Mbr_partition_4 != Modelo.PARTICION{}) {
		//REPORTE PARTICION 4
		strReporte += "<tr><td>MBR Particion No. 4</td><td></td></tr>\n"
		strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_4.Part_status) +"</td></tr>\n"
		strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_4.Part_type) +"</td></tr>\n"
		strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_4.Part_fit[:]) +"</td></tr>\n"
		strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_4.Part_start, 10) +"</td></tr>\n"
		strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_4.Part_size, 10) +"</td></tr>\n"
		strReporte += "<tr><td>Name</td><td>"+ s4 +"</td></tr>\n"
		if (m.Mbr_partition_4.Part_type == 'E' || m.Mbr_partition_4.Part_type == 'e') {
			for i := 0; i < 25; i++ {
				if (m.Mbr_partition_4.Part_EBR[i] != Modelo.EBR{}) {
					if(m.Mbr_partition_4.Part_EBR[i].Part_logica != Modelo.PARTICION_LOGICA{}) {
						strReporte += "<tr><td>EBR Particion</td><td></td></tr>\n"
						strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_4.Part_EBR[i].Part_logica.Part_status) +"</td></tr>\n"
						strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_4.Part_EBR[i].Part_logica.Part_type) +"</td></tr>\n"
						strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_4.Part_EBR[i].Part_logica.Part_fit[:]) +"</td></tr>\n"
						strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_4.Part_EBR[i].Part_logica.Part_start, 10) +"</td></tr>\n"
						strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_4.Part_EBR[i].Part_logica.Part_size, 10) +"</td></tr>\n"
						strReporte += "<tr><td>Name</td><td>"+ BytesToString(m.Mbr_partition_4.Part_EBR[i].Part_logica.Part_name[:]) +"</td></tr>\n"
					}
				}
			}
		}
	}

	strReporte += "</table>\n"


	strReporte = strReporte +">];}"
	_, b, _, _ := runtime.Caller(0)

	d := path.Join(path.Dir(b))
	Ubicacion.CrearArchivoRep(filepath.Dir(d) + "/ReporteControlador/mbr.dot", strReporte)

	
	var strUrl = filepath.Dir(d) + "/ReporteControlador/mbr.dot"
	var strUrl2 = strPath

	executeCommand(strUrl, strUrl2);

	color.Blue("╔══════════════════════════════════════════════════╗")
	color.Blue("               REPORTE MBR GENERADO :D")
	color.Blue("  PATH: " + strPath)
	color.Blue("╚══════════════════════════════════════════════════╝")
}

func DISK(ubicacion string, strPath string) {
	file, err := os.OpenFile(ubicacion, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close() 
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
	data := lecturaBytes(file, size)
	//DECODIFICACION EN BINARIO
	buffer := bytes.NewBuffer(data)
	
	//GUARDAR VARIABLE EN M
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	var strReporte = "digraph {"+
		"tbl [ "+
		"shape=plaintext "+
		"label=< "+

	"<table>\n";

	//REPORTE PARTICION 1
	//s1 := BytesToString(m.Mbr_partition_1.Part_name[:])
	//s2 := BytesToString(m.Mbr_partition_2.Part_name[:])
	//s3 := BytesToString(m.Mbr_partition_3.Part_name[:])
	//s4 := BytesToString(m.Mbr_partition_4.Part_name[:])
	strReporte += "<tr>"
	strReporte += "<td></td>"
	if (m.Mbr_partition_1 != Modelo.PARTICION{}) {
		if (m.Mbr_partition_1.Part_type == 'E' || m.Mbr_partition_1.Part_type == 'e') {
			strReporte += "<td>Particion Extendida  - "+strconv.FormatInt(m.Mbr_partition_1.Part_size, 10)+"</td>"
		} else {
			strReporte += "<td></td>"
		}
	} else {
		strReporte += "<td></td>"
	}
	if (m.Mbr_partition_2 != Modelo.PARTICION{}) {
		if (m.Mbr_partition_2.Part_type == 'E' || m.Mbr_partition_2.Part_type == 'e') {
			strReporte += "<td>Particion Extendida  - "+strconv.FormatInt(m.Mbr_partition_2.Part_size, 10)+"</td>"
		} else {
			strReporte += "<td></td>"
		}
	} else {
		strReporte += "<td></td>"
	}
	if (m.Mbr_partition_3 != Modelo.PARTICION{}) {
		if (m.Mbr_partition_3.Part_type == 'E' || m.Mbr_partition_3.Part_type == 'e') {
			strReporte += "<td>Particion Extendida  - "+strconv.FormatInt(m.Mbr_partition_3.Part_size, 10)+"</td>"
		} else {
			strReporte += "<td></td>"
		}
	} else {
		strReporte += "<td></td>"
	}
	if (m.Mbr_partition_4 != Modelo.PARTICION{}) {
		if (m.Mbr_partition_4.Part_type == 'E' || m.Mbr_partition_4.Part_type == 'e') {
			strReporte += "<td>Particion Extendida  - "+strconv.FormatInt(m.Mbr_partition_4.Part_size, 10)+"</td>"
		} else {
			strReporte += "<td></td>"
		}
	} else {
		strReporte += "<td></td>"
	}
	strReporte += "<td></td>"
	strReporte += "</tr>"

	strReporte += "<tr>"
	strReporte += "<td>MBR</td>"

	//PORCENTAJES
	if (m.Mbr_partition_1 != Modelo.PARTICION{}) {
		if (m.Mbr_partition_1.Part_type == 'E' || m.Mbr_partition_1.Part_type == 'e') {
			strReporte += "<td>"
			strReporte += "<table >"
			strReporte += "<tr>"
				//for
				for i := 0; i < 25; i++ {
					if (m.Mbr_partition_1.Part_EBR[i] != Modelo.EBR{}) {
						if(m.Mbr_partition_1.Part_EBR[i].Part_logica != Modelo.PARTICION_LOGICA{}) {
							strReporte += "<td>EBR</td>"
							strReporte += "<td>Partición: "+BytesToString(m.Mbr_partition_1.Part_EBR[i].Part_logica.Part_name[:])+"</td>"
						} else {
							strReporte += "<td>EBR</td>"
						}
					}
				}
			strReporte += "</tr>"
            strReporte += "</table>"
			strReporte += "</td>"
		} else {
			strReporte += "<td>Particion I - "+strconv.FormatInt(m.Mbr_partition_1.Part_size, 10)+"</td>"
		}
	} else {
		strReporte += "<td>Disponible</td>"
	}
	if (m.Mbr_partition_2 != Modelo.PARTICION{}) {
		if (m.Mbr_partition_2.Part_type == 'E' || m.Mbr_partition_2.Part_type == 'e') {
			strReporte += "<td>"
			strReporte += "<table >"
			strReporte += "<tr>"
				//for
				for i := 0; i < 25; i++ {
					if (m.Mbr_partition_2.Part_EBR[i] != Modelo.EBR{}) {
						if(m.Mbr_partition_2.Part_EBR[i].Part_logica != Modelo.PARTICION_LOGICA{}) {
							strReporte += "<td>EBR</td>"
							strReporte += "<td>Partición: "+BytesToString(m.Mbr_partition_2.Part_EBR[i].Part_logica.Part_name[:])+"</td>"
						} else {
							strReporte += "<td>EBR</td>"
						}
					}
				}
			strReporte += "</tr>"
            strReporte += "</table>"
			strReporte += "</td>"
		} else {
			strReporte += "<td>Particion II - "+strconv.FormatInt(m.Mbr_partition_2.Part_size, 10)+"</td>"
		}
	} else {
		strReporte += "<td>Disponible</td>"
	}
	if (m.Mbr_partition_3 != Modelo.PARTICION{}) {
		if (m.Mbr_partition_3.Part_type == 'E' || m.Mbr_partition_3.Part_type == 'e') {
			strReporte += "<td>"
			strReporte += "<table >"
			strReporte += "<tr>"
				//for
				for i := 0; i < 25; i++ {
					if (m.Mbr_partition_3.Part_EBR[i] != Modelo.EBR{}) {
						if(m.Mbr_partition_3.Part_EBR[i].Part_logica != Modelo.PARTICION_LOGICA{}) {
							strReporte += "<td>EBR</td>"
							strReporte += "<td>Partición: "+BytesToString(m.Mbr_partition_3.Part_EBR[i].Part_logica.Part_name[:])+"</td>"
						} else {
							strReporte += "<td>EBR</td>"
						}
					}
				}
			strReporte += "</tr>"
            strReporte += "</table>"
			strReporte += "</td>"
		} else {
			strReporte += "<td>Particion III - "+strconv.FormatInt(m.Mbr_partition_3.Part_size, 10)+"</td>"
		}
	} else {
		strReporte += "<td>Disponible</td>"
	}
	if (m.Mbr_partition_4 != Modelo.PARTICION{}) {
		if (m.Mbr_partition_4.Part_type == 'E' || m.Mbr_partition_4.Part_type == 'e') {
			strReporte += "<td>"
			strReporte += "<table >"
			strReporte += "<tr>"
				//for
				for i := 0; i < 25; i++ {
					if (m.Mbr_partition_4.Part_EBR[i] != Modelo.EBR{}) {
						if(m.Mbr_partition_4.Part_EBR[i].Part_logica != Modelo.PARTICION_LOGICA{}) {
							strReporte += "<td>EBR</td>"
							strReporte += "<td>Partición: "+BytesToString(m.Mbr_partition_4.Part_EBR[i].Part_logica.Part_name[:])+"</td>"
						} else {
							strReporte += "<td>EBR</td>"
						}
					}
				}
			strReporte += "</tr>"
            strReporte += "</table>"
			strReporte += "</td>"
		} else {
			strReporte += "<td>Particion IV - "+strconv.FormatInt(m.Mbr_partition_4.Part_size, 10)+"</td>"
		}
	} else {
		strReporte += "<td>Disponible</td>"
	}
	strReporte += "<td>Disponible - "+strconv.FormatInt(m.Mbr_size_disponible, 10)+"</td>"
	strReporte += "</tr>"

	strReporte += "</table>\n"


	strReporte = strReporte +">];}"
	_, b, _, _ := runtime.Caller(0)

	d := path.Join(path.Dir(b))
	Ubicacion.CrearArchivoRep(filepath.Dir(d) + "/ReporteControlador/disk.dot", strReporte)

	
	var strUrl = filepath.Dir(d) + "/ReporteControlador/disk.dot"
	var strUrl2 = strPath

	executeCommand(strUrl, strUrl2);

	color.Blue("╔══════════════════════════════════════════════════╗")
	color.Blue("               REPORTE DISK GENERADO :D")
	color.Blue("  PATH: " + strPath)
	color.Blue("╚══════════════════════════════════════════════════╝")
}

func SUPERBOOT(id string, strPath string)  {
	if (MOUNTControlador.Buscar(id) == true) {
		mount := Modelo.MOUNT{};
		mount = MOUNTControlador.BuscarMOUNT(id);
		//MOUNT
		//fmt.Println(mount)

		//SE ABRE EL ARCHIVO
		file, err := os.OpenFile(mount.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error, no se pudo abrir el disco duro")
		}
		sb := Modelo.SUPERBOOT{}
		//COMPARAR
		if (mount.Mount_particion_logica == Modelo.PARTICION_LOGICA{}) {
			//PARTICION NORMAL
			sb = ObtenerSuperBoot(mount.Mount_partition.Part_start, file)

		} else if (mount.Mount_partition == Modelo.PARTICION{}) {
			//PARTCION LOGICA
			sb = ObtenerSuperBoot(mount.Mount_particion_logica.Part_start, file)
		}
		//OBTENGO EL SUPER BOOT
		//fmt.Println("SUPERBOOT")
		//fmt.Println(sb)

		/**
		* REALIZAR REPORTE
 		*/
		strReporte := "digraph {"+
			"tbl [ "+
			"shape=plaintext "+
			"label=< "+
			"<table border='1' cellborder='1'> ";
		strReporte += "<tr><td>SB_nombre_hd</td><td>"+BytesToString(sb.SB_nombre_hd[:])+"</td></tr>"
		strReporte += "<tr><td>SB_arbol_virtual_count</td><td>"+ strconv.FormatInt(sb.SB_arbol_virtual_count, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_detalle_directorio_count</td><td>"+ strconv.FormatInt(sb.SB_detalle_directorio_count, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_inodos_count</td><td>"+ strconv.FormatInt(sb.SB_inodos_count, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_bloques_count</td><td>"+ strconv.FormatInt(sb.SB_bloques_count, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_arbol_virtual_free</td><td>"+ strconv.FormatInt(sb.SB_arbol_virtual_free, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_detalle_directorio_free</td><td>"+ strconv.FormatInt(sb.SB_detalle_directorio_free, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_inodos_free</td><td>"+ strconv.FormatInt(sb.SB_inodos_free, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_bloques_free</td><td>"+ strconv.FormatInt(sb.SB_bloques_free, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_date_creacion</td><td>"+string(sb.SB_date_creacion[:])+"</td></tr>"
		strReporte += "<tr><td>SB_date_ultimo_montaje</td><td>"+string(sb.SB_date_ultimo_montaje[:])+"</td></tr>"
		strReporte += "<tr><td>SB_montaje_count</td><td>"+ strconv.FormatInt(sb.SB_montaje_count, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_ap_bitmap_arbol_directorio</td><td>"+ strconv.FormatInt(sb.SB_ap_bitmap_arbol_directorio, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_ap_arbol_directorio</td><td>"+ strconv.FormatInt(sb.SB_ap_arbol_directorio, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_ap_bitmap_detalle_directorio</td><td>"+ strconv.FormatInt(sb.SB_ap_bitmap_detalle_directorio, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_ap_detalle_directorio</td><td>"+ strconv.FormatInt(sb.SB_ap_detalle_directorio, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_ap_bitmap_tabla_inodo</td><td>"+ strconv.FormatInt(sb.SB_ap_bitmap_tabla_inodo, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_ap_tabla_inodo</td><td>"+ strconv.FormatInt(sb.SB_ap_tabla_inodo, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_ap_bitmap_bloques</td><td>"+ strconv.FormatInt(sb.SB_ap_bitmap_bloques, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_ap_bloques</td><td>"+ strconv.FormatInt(sb.SB_ap_bloques, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_ap_log</td><td>"+ strconv.FormatInt(sb.SB_ap_log, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_size_struct_arbol_directorio</td><td>"+ strconv.FormatInt(sb.SB_size_struct_arbol_directorio, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_size_struct_detalle_directorio</td><td>"+ strconv.FormatInt(sb.SB_size_struct_detalle_directorio, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_size_struct_inodo</td><td>"+ strconv.FormatInt(sb.SB_size_struct_inodo, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_size_struct_bloque</td><td>"+ strconv.FormatInt(sb.SB_size_struct_bloque, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_free_first_bit_struct_arbol_directorio</td><td>"+ strconv.FormatInt(sb.SB_free_first_bit_struct_arbol_directorio, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_free_first_bit_struct_detalle_directorio</td><td>"+ strconv.FormatInt(sb.SB_free_first_bit_struct_detalle_directorio, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_free_first_bit_struct_inodo</td><td>"+ strconv.FormatInt(sb.SB_free_first_bit_struct_inodo, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_free_first_bit_struct_bloque</td><td>"+ strconv.FormatInt(sb.SB_free_first_bit_struct_bloque, 10) +"</td></tr>"
		strReporte += "<tr><td>SB_magic_num</td><td>201801237</td></tr>"
		strReporte += "</table>>];}"
	
		_, b, _, _ := runtime.Caller(0)

		d := path.Join(path.Dir(b))
		Ubicacion.CrearArchivoRep(filepath.Dir(d) + "/ReporteControlador/sb.dot", strReporte)

		
		var strUrl = filepath.Dir(d) + "/ReporteControlador/sb.dot"
		var strUrl2 = strPath

		executeCommand(strUrl, strUrl2);

		color.Blue("╔══════════════════════════════════════════════════╗")
		color.Blue("             REPORTE SUPERBOOT GENERADO :D")
		color.Blue("  PATH: " + strPath)
		color.Blue("╚══════════════════════════════════════════════════╝")
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("     La partición no se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func BM_AVD(id string, strPath string)  {
	if (MOUNTControlador.Buscar(id) == true) {
		mount := Modelo.MOUNT{};
		mount = MOUNTControlador.BuscarMOUNT(id);
		//MOUNT
		//fmt.Println(mount)

		//SE ABRE EL ARCHIVO
		file, err := os.OpenFile(mount.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error, no se pudo abrir el disco duro")
		}
		sb := Modelo.SUPERBOOT{}
		//COMPARAR
		if (mount.Mount_particion_logica == Modelo.PARTICION_LOGICA{}) {
			//PARTICION NORMAL
			sb = ObtenerSuperBoot(mount.Mount_partition.Part_start, file)

		} else if (mount.Mount_partition == Modelo.PARTICION{}) {
			//PARTCION LOGICA
			sb = ObtenerSuperBoot(mount.Mount_particion_logica.Part_start, file)
		}
		//OBTENGO EL SUPER BOOT
		//fmt.Println("SUPERBOOT")
		//fmt.Println(sb)

		file.Seek(sb.SB_ap_bitmap_arbol_directorio, 0)
		arregloBitmap := make([]byte, (sb.SB_ap_arbol_directorio - sb.SB_ap_bitmap_arbol_directorio))
		n1, err := file.Read(arregloBitmap)
		if err != nil {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("            Error lectura bitmaps D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}

		/**
		* REALIZAR REPORTE
 		*/
		strReporte := "digraph H {\n"+
			"graph[label=\"BITMAP Arbol de Directorio\", labelloc=t, fontsize=20, compound=true];"+
			"aHtmlTable [ \n"+
			"shape=plaintext \n"+
			"color=blue \n"+
			"label=< \n"+
			"<table border='1' cellborder='0'> \n";
			contador := 0;
			contador2 := 0;
			total := 0;
			for i := 0; i < len(string(arregloBitmap[:n1])); i++ {
				if contador == 0 {
					strReporte += "\t<tr><td>|</td>\n"
				}
				//fmt.Println(string(arregloBitmap[:n1][i]))
				contador++;
				contador2++;
				strReporte += "<td>"+string(arregloBitmap[:n1][i])+"</td><td>|</td>"
				if contador == 30 {
					strReporte += "</tr>\n"
					contador = 0;
				}
				total = len(string(arregloBitmap[:n1])) -  contador2;
				if (i == len(string(arregloBitmap[:n1])) -1) {
					if total < contador {
						strReporte += "</tr>\n"
						contador = 0;
					}
				}
			}

		strReporte += "\n</table>>];}"
	
		_, b, _, _ := runtime.Caller(0)

		d := path.Join(path.Dir(b))
		Ubicacion.CrearArchivoRep(filepath.Dir(d) + "/ReporteControlador/bmavd.dot", strReporte)

		
		var strUrl = filepath.Dir(d) + "/ReporteControlador/bmavd.dot"
		var strUrl2 = strPath

		executeCommand(strUrl, strUrl2);

		color.Blue("╔══════════════════════════════════════════════════╗")
		color.Blue("           REPORTE BITMAP AVD GENERADO :D")
		color.Blue("  PATH: " + strPath)
		color.Blue("╚══════════════════════════════════════════════════╝")
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("     La partición no se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func BM_DD(id string, strPath string)  {
	if (MOUNTControlador.Buscar(id) == true) {
		mount := Modelo.MOUNT{};
		mount = MOUNTControlador.BuscarMOUNT(id);
		//MOUNT
		//fmt.Println(mount)

		//SE ABRE EL ARCHIVO
		file, err := os.OpenFile(mount.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error, no se pudo abrir el disco duro")
		}
		sb := Modelo.SUPERBOOT{}
		//COMPARAR
		if (mount.Mount_particion_logica == Modelo.PARTICION_LOGICA{}) {
			//PARTICION NORMAL
			sb = ObtenerSuperBoot(mount.Mount_partition.Part_start, file)

		} else if (mount.Mount_partition == Modelo.PARTICION{}) {
			//PARTCION LOGICA
			sb = ObtenerSuperBoot(mount.Mount_particion_logica.Part_start, file)
		}
		//OBTENGO EL SUPER BOOT
		//fmt.Println("SUPERBOOT")
		//fmt.Println(sb)

		file.Seek(sb.SB_ap_bitmap_detalle_directorio, 0)
		arregloBitmap := make([]byte, (sb.SB_ap_detalle_directorio - sb.SB_ap_bitmap_detalle_directorio))
		n1, err := file.Read(arregloBitmap)
		if err != nil {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("            Error lectura bitmaps D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}

		/**
		* REALIZAR REPORTE
 		*/
		strReporte := "digraph H {\n"+
			"graph[label=\"BITMAP Detalle de Directorio\", labelloc=t, fontsize=20, compound=true];"+
			"aHtmlTable [ \n"+
			"shape=plaintext \n"+
			"color=blue \n"+
			"label=< \n"+
			"<table border='1' cellborder='0'> \n";
			contador := 0;
			contador2 := 0;
			total := 0;
			for i := 0; i < len(string(arregloBitmap[:n1])); i++ {
				if contador == 0 {
					strReporte += "\t<tr><td>|</td>\n"
				}
				//fmt.Println(string(arregloBitmap[:n1][i]))
				contador++;
				contador2++;
				strReporte += "<td>"+string(arregloBitmap[:n1][i])+"</td><td>|</td>"
				if contador == 30 {
					strReporte += "</tr>\n"
					contador = 0;
				}
				total = len(string(arregloBitmap[:n1])) -  contador2;
				if (i == len(string(arregloBitmap[:n1])) -1) {
					if total < contador {
						strReporte += "</tr>\n"
						contador = 0;
					}
				}
			}

		strReporte += "\n</table>>];}"
	
		_, b, _, _ := runtime.Caller(0)

		d := path.Join(path.Dir(b))
		Ubicacion.CrearArchivoRep(filepath.Dir(d) + "/ReporteControlador/bmdd.dot", strReporte)

		
		var strUrl = filepath.Dir(d) + "/ReporteControlador/bmdd.dot"
		var strUrl2 = strPath

		executeCommand(strUrl, strUrl2);

		color.Blue("╔══════════════════════════════════════════════════╗")
		color.Blue("           REPORTE BITMAP DD GENERADO :D")
		color.Blue("  PATH: " + strPath)
		color.Blue("╚══════════════════════════════════════════════════╝")
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("     La partición no se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func BM_INODO(id string, strPath string)  {
	if (MOUNTControlador.Buscar(id) == true) {
		mount := Modelo.MOUNT{};
		mount = MOUNTControlador.BuscarMOUNT(id);
		//MOUNT
		//fmt.Println(mount)

		//SE ABRE EL ARCHIVO
		file, err := os.OpenFile(mount.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error, no se pudo abrir el disco duro")
		}
		sb := Modelo.SUPERBOOT{}
		//COMPARAR
		if (mount.Mount_particion_logica == Modelo.PARTICION_LOGICA{}) {
			//PARTICION NORMAL
			sb = ObtenerSuperBoot(mount.Mount_partition.Part_start, file)

		} else if (mount.Mount_partition == Modelo.PARTICION{}) {
			//PARTCION LOGICA
			sb = ObtenerSuperBoot(mount.Mount_particion_logica.Part_start, file)
		}
		//OBTENGO EL SUPER BOOT
		//fmt.Println("SUPERBOOT")
		//fmt.Println(sb)

		file.Seek(sb.SB_ap_bitmap_tabla_inodo, 0)
		arregloBitmap := make([]byte, (sb.SB_ap_tabla_inodo - sb.SB_ap_bitmap_tabla_inodo))
		n1, err := file.Read(arregloBitmap)
		if err != nil {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("            Error lectura bitmaps D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}

		/**
		* REALIZAR REPORTE
 		*/
		strReporte := "digraph H {\n"+
			"graph[label=\"BITMAP Inodo\", labelloc=t, fontsize=20, compound=true];"+
			"aHtmlTable [ \n"+
			"shape=plaintext \n"+
			"color=blue \n"+
			"label=< \n"+
			"<table border='1' cellborder='0'> \n";
			contador := 0;
			contador2 := 0;
			total := 0;
			for i := 0; i < len(string(arregloBitmap[:n1])); i++ {
				if contador == 0 {
					strReporte += "\t<tr><td>|</td>\n"
				}
				//fmt.Println(string(arregloBitmap[:n1][i]))
				contador++;
				contador2++;
				strReporte += "<td>"+string(arregloBitmap[:n1][i])+"</td><td>|</td>"
				if contador == 30 {
					strReporte += "</tr>\n"
					contador = 0;
				}
				total = len(string(arregloBitmap[:n1])) -  contador2;
				if (i == len(string(arregloBitmap[:n1])) -1) {
					if total < contador {
						strReporte += "</tr>\n"
						contador = 0;
					}
				}
			}

		strReporte += "\n</table>>];}"
	
		_, b, _, _ := runtime.Caller(0)

		d := path.Join(path.Dir(b))
		Ubicacion.CrearArchivoRep(filepath.Dir(d) + "/ReporteControlador/bminodo.dot", strReporte)

		
		var strUrl = filepath.Dir(d) + "/ReporteControlador/bminodo.dot"
		var strUrl2 = strPath

		executeCommand(strUrl, strUrl2);

		color.Blue("╔══════════════════════════════════════════════════╗")
		color.Blue("           REPORTE BITMAP INODO GENERADO :D")
		color.Blue("  PATH: " + strPath)
		color.Blue("╚══════════════════════════════════════════════════╝")
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("     La partición no se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func BM_BLOQUE(id string, strPath string)  {
	if (MOUNTControlador.Buscar(id) == true) {
		mount := Modelo.MOUNT{};
		mount = MOUNTControlador.BuscarMOUNT(id);
		//MOUNT
		//fmt.Println(mount)

		//SE ABRE EL ARCHIVO
		file, err := os.OpenFile(mount.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close()
		if err != nil {
			fmt.Println("Hay un error, no se pudo abrir el disco duro")
		}
		sb := Modelo.SUPERBOOT{}
		//COMPARAR
		if (mount.Mount_particion_logica == Modelo.PARTICION_LOGICA{}) {
			//PARTICION NORMAL
			sb = ObtenerSuperBoot(mount.Mount_partition.Part_start, file)

		} else if (mount.Mount_partition == Modelo.PARTICION{}) {
			//PARTCION LOGICA
			sb = ObtenerSuperBoot(mount.Mount_particion_logica.Part_start, file)
		}
		//OBTENGO EL SUPER BOOT
		//fmt.Println("SUPERBOOT")
		//fmt.Println(sb)

		file.Seek(sb.SB_ap_bitmap_bloques, 0)
		arregloBitmap := make([]byte, (sb.SB_ap_bloques - sb.SB_ap_bitmap_bloques))
		n1, err := file.Read(arregloBitmap)
		if err != nil {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("            Error lectura bitmaps D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}

		/**
		* REALIZAR REPORTE
 		*/
		strReporte := "digraph H {\n"+
			"graph[label=\"BITMAP Bloque\", labelloc=t, fontsize=20, compound=true];"+
			"aHtmlTable [ \n"+
			"shape=plaintext \n"+
			"color=blue \n"+
			"label=< \n"+
			"<table border='1' cellborder='0'> \n";
			contador := 0;
			contador2 := 0;
			total := 0;
			for i := 0; i < len(string(arregloBitmap[:n1])); i++ {
				if contador == 0 {
					strReporte += "\t<tr><td>|</td>\n"
				}
				//fmt.Println(string(arregloBitmap[:n1][i]))
				contador++;
				contador2++;
				strReporte += "<td>"+string(arregloBitmap[:n1][i])+"</td><td>|</td>"
				if contador == 30 {
					strReporte += "</tr>\n"
					contador = 0;
				}
				total = len(string(arregloBitmap[:n1])) -  contador2;
				if (i == len(string(arregloBitmap[:n1])) -1) {
					if total < contador {
						strReporte += "</tr>\n"
						contador = 0;
					}
				}
			}

		strReporte += "\n</table>>];}"
	
		_, b, _, _ := runtime.Caller(0)

		d := path.Join(path.Dir(b))
		Ubicacion.CrearArchivoRep(filepath.Dir(d) + "/ReporteControlador/bmbloque.dot", strReporte)

		
		var strUrl = filepath.Dir(d) + "/ReporteControlador/bmbloque.dot"
		var strUrl2 = strPath

		executeCommand(strUrl, strUrl2);

		color.Blue("╔══════════════════════════════════════════════════╗")
		color.Blue("          REPORTE BITMAP BLOQUE GENERADO :D")
		color.Blue("  PATH: " + strPath)
		color.Blue("╚══════════════════════════════════════════════════╝")
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("     La partición no se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func ObtenerSuperBoot(inicioParticion int64, file *os.File) Modelo.SUPERBOOT {
	sb := Modelo.SUPERBOOT{}
	var sbSize int = int(unsafe.Sizeof(sb))

	file.Seek(inicioParticion, 0)
	datos := lecturaBytes(file, sbSize)
	buffer := bytes.NewBuffer(datos)
	err := binary.Read(buffer, binary.BigEndian, &sb)
	if err != nil {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("               Error de lectura D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
	return sb
}

/**
 * LECTURA DE BYTES
 */
 func lecturaBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes
	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func BytesToString(byteArray []byte) string {
	n := bytes.Index(byteArray, []byte{0})
	return string(byteArray[:n])
}