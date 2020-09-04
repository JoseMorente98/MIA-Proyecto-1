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
func TipoReporte(nombre string, ubicacion string, id string, ruta string, pathDisco string)  {
	switch strings.ToLower(nombre) {
		case "mbr":
			MBR(pathDisco, ubicacion);
		case "disk":
		case "sb":
		case "bm_arbdir":
		case "bm_detdir":
		case "bm_inode":
		case "bm_block":
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
	data := leerBytes(file, size)
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
	"label=<"+

	"<table border='1' cellborder='1'>"+
	"<tr><td>Name</td><td>Valor</td></tr>"

	//LLENADO MBR
	strReporte += "<tr><td>MBR Tamano</td><td>"+ strconv.FormatInt(m.Mbr_size, 10) +"</td></tr>"
	strReporte += "<tr><td>MBR Fecha Creacion</td><td>"+ string(m.Mbr_date[:]) +"</td></tr>"
	strReporte += "<tr><td>MBR Disk Asignature</td><td>"+  strconv.FormatInt(m.Mbr_disk_signature, 10) +"</td></tr>"

	//REPORTE PARTICION 1
	s1 := BytesToString(m.Mbr_partition_1.Part_name[:])
	s2 := BytesToString(m.Mbr_partition_2.Part_name[:])
	s3 := BytesToString(m.Mbr_partition_3.Part_name[:])
	s4 := BytesToString(m.Mbr_partition_4.Part_name[:])
	strReporte += "<tr><td>MBR Particion No. 1</td><td></td></tr>"
	strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_1.Part_status) +"</td></tr>"
	strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_1.Part_type) +"</td></tr>"
	strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_1.Part_fit[:]) +"</td></tr>"
	strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_1.Part_start, 10) +"</td></tr>"
	strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_1.Part_size, 10) +"</td></tr>"
	strReporte += "<tr><td>Name</td><td>"+ s1 +"</td></tr>"
	//REPORTE PARTICION 2
	strReporte += "<tr><td>MBR Particion No. 2</td><td></td></tr>"
	strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_2.Part_status) +"</td></tr>"
	strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_2.Part_type) +"</td></tr>"
	strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_2.Part_fit[:]) +"</td></tr>"
	strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_2.Part_start, 10) +"</td></tr>"
	strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_2.Part_size, 10) +"</td></tr>"
	strReporte += "<tr><td>Name</td><td>"+ s2 +"</td></tr>"
	//REPORTE PARTICION 3
	strReporte += "<tr><td>MBR Particion No. 3</td><td></td></tr>"
	strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_3.Part_status) +"</td></tr>"
	strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_3.Part_type) +"</td></tr>"
	strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_3.Part_fit[:]) +"</td></tr>"
	strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_3.Part_start, 10) +"</td></tr>"
	strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_3.Part_size, 10) +"</td></tr>"
	strReporte += "<tr><td>Name</td><td>"+ s3 +"</td></tr>"
	//REPORTE PARTICION 4
	strReporte += "<tr><td>MBR Particion No. 4</td><td></td></tr>"
	strReporte += "<tr><td>Status</td><td>"+ string(m.Mbr_partition_4.Part_status) +"</td></tr>"
	strReporte += "<tr><td>Type</td><td>"+ string(m.Mbr_partition_4.Part_type) +"</td></tr>"
	strReporte += "<tr><td>Fit</td><td>"+ string(m.Mbr_partition_4.Part_fit[:]) +"</td></tr>"
	strReporte += "<tr><td>Start</td><td>"+ strconv.FormatInt(m.Mbr_partition_4.Part_start, 10) +"</td></tr>"
	strReporte += "<tr><td>Size</td><td>"+ strconv.FormatInt(m.Mbr_partition_4.Part_size, 10) +"</td></tr>"
	strReporte += "<tr><td>Name</td><td>"+ s4 +"</td></tr>"

	strReporte += "</table>"


	strReporte = strReporte +">];}"
	_, b, _, _ := runtime.Caller(0)

	d := path.Join(path.Dir(b))
	Ubicacion.CrearArchivo(filepath.Dir(d) + "/ReporteControlador/mbr.dot", strReporte)

	
	var strUrl = filepath.Dir(d) + "/ReporteControlador/mbr.dot"
	var strUrl2 = strPath

	executeCommand(strUrl, strUrl2);
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

func BytesToString(byteArray []byte) string {
	n := bytes.Index(byteArray, []byte{0})
	return string(byteArray[:n])
}