package MKFSControlador

import (
	"strings"
	"github.com/fatih/color"
	"os"
	"bytes"
	"unsafe"
	"strconv"
	"log"
	"time"
	"encoding/binary"
	"path/filepath"

	Modelo "../../Modelo"
	//Path "../../Utilidades/Path"
	MOUNTControlador "../MOUNTControlador"
)

/**
 * FUNCION MKFS
 */
func MKFS(id string, types string, add string, unit string) {
	if id != ""  {
		if types != "" {
			if strings.ToLower(types) == "full" || strings.ToLower(types) == "fast" {
				types = strings.ToUpper(types);
			} else {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("       Solo puedes utilizar FULL o FAST D:")
				color.Red("╚══════════════════════════════════════════════════╝")
			}
		} else {
			types = "FULL";
		}

		if unit != "" {
			if strings.ToLower(unit) == "b" || strings.ToLower(unit) == "k" || strings.ToLower(unit) == "m"  {
				unit = strings.ToUpper(unit);
			} else {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("         Solo puedes utilizar B, K o M D:")
				color.Red("╚══════════════════════════════════════════════════╝")
			}
		} else {
			unit = "K";
		}
		formatearParticion(unit, types, id, add);
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("       Verificar parametros obligatorios D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func formatearParticion(unit string, types string, id string, add string)  {
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
	color.Red(string(sizeUNIT));

	
	if (MOUNTControlador.Buscar(id) == true) {
		mount := Modelo.MOUNT{};
		mount = MOUNTControlador.BuscarMOUNT(id);
		//Path.CrearArchivo("/home/jmorente/Documentos/CodeApp GT/Archivos/" + "user" + string(mount.Mount_id[:]) + ".txt", 
		//"1,G,root\n1,U,root,root,201801237\n");
		FormatearParticion(mount);
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("     La partición no se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func FormatearParticion(mount Modelo.MOUNT)  {
	file, err := os.OpenFile(mount.Mount_path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close()
	if err != nil {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("     Ha ocurrido un error D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	} else {
		//PARTICION 
		if (mount.Mount_particion_logica == Modelo.PARTICION_LOGICA{}) {
			Particion(mount.Mount_partition, file, filepath.Base(mount.Mount_path));
		//PARTICION LOGICA
		} else if (mount.Mount_partition == Modelo.PARTICION{}) {
			
		}
	}
}

func Particion(particion Modelo.PARTICION, file *os.File, nombre string)  {
	/*for i := particion.Part_start; i < particion.Part_end; i++ {
		var init int8 = 'p'
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escrituraBytes(file, binarioTemp.Bytes())
	}*/

	var sizeSB int64 = int64(unsafe.Sizeof(Modelo.SUPERBOOT{}))
	var sizeAVD int64 = int64(unsafe.Sizeof(Modelo.AVD{}))
	var sizeDD int64 = int64(unsafe.Sizeof(Modelo.DD{}))
	var sizeInodos int64 = int64(unsafe.Sizeof(Modelo.INODO{}))
	var sizeBloques int64 = int64(unsafe.Sizeof(Modelo.BLOQUE{}))
	var sizeBitacora int64 = int64(unsafe.Sizeof(Modelo.BITACORA{}))

	var cantidadEstructuras int64 = (particion.Part_size - (2 * int64(sizeSB))) / (27 + int64(sizeAVD) + int64(sizeDD) + (5*int64(sizeInodos) + (20 * int64(sizeBloques)) + int64(sizeBitacora)))
	var cantidadAVD int64 = cantidadEstructuras
	var cantidadDD int64 = cantidadEstructuras
	var cantidadInodos int64 = 5 * cantidadEstructuras
	var cantidadBloques int64 = 20 * cantidadEstructuras
	var cantidadBitacoras int64 = cantidadEstructuras

	Inicio_bitmapAVD := particion.Part_start + sizeSB
	Inicio_AVD := Inicio_bitmapAVD + cantidadAVD
	Inicio_bitmapDD := Inicio_AVD + (sizeAVD * cantidadAVD)
	Inicio_DD := Inicio_bitmapDD + cantidadDD
	Inicio_bitmapInodo := Inicio_DD + (sizeDD * cantidadDD)
	Inicio_Inodos := Inicio_bitmapInodo + cantidadInodos
	Inicio_bitmapBloque := Inicio_Inodos + (sizeInodos * cantidadInodos)
	Inicio_Bloque := Inicio_bitmapBloque + cantidadBloques
	Inicio_Bitacora := Inicio_Bloque + (sizeBloques * cantidadBloques)
	Inicio_SBCopia := Inicio_Bitacora + (sizeBitacora * cantidadBitacoras)

	/**
	* ESCRITURA DE SUPERBLOQUE
	*/
	file.Seek(particion.Part_start, 0)
	var SB = Modelo.SUPERBOOT{};
	fechaHora := time.Now();
	copy(SB.SB_nombre_hd[:], nombre);
	copy(SB.SB_date_creacion[:], fechaHora.Format("2000-01-01 00:00:01"))
	copy(SB.SB_date_ultimo_montaje[:], fechaHora.Format("2000-01-01 00:00:01"))
	SB.SB_arbol_virtual_count = cantidadEstructuras;
	SB.SB_detalle_directorio_count = cantidadDD;
	SB.SB_inodos_count = cantidadInodos;
	SB.SB_bloques_count = cantidadBloques;
	SB.SB_arbol_virtual_free = cantidadEstructuras;
	SB.SB_detalle_directorio_free = cantidadDD;
	SB.SB_inodos_free = cantidadInodos;
	SB.SB_bloques_free = cantidadBloques;
	SB.SB_montaje_count	= 1;
	SB.SB_ap_bitmap_arbol_directorio = Inicio_bitmapAVD;
	SB.SB_ap_arbol_directorio = Inicio_AVD;
	SB.SB_ap_bitmap_detalle_directorio = Inicio_bitmapDD;
	SB.SB_ap_detalle_directorio = Inicio_DD;
	SB.SB_ap_bitmap_tabla_inodo = Inicio_bitmapInodo;
	SB.SB_ap_tabla_inodo = Inicio_Inodos;
	SB.SB_ap_bitmap_bloques = Inicio_bitmapBloque;
	SB.SB_ap_bloques = Inicio_Bloque;
	SB.SB_ap_log = Inicio_Bitacora;
	SB.SB_size_struct_arbol_directorio = sizeAVD;
	SB.SB_size_struct_detalle_directorio = sizeDD;
	SB.SB_size_struct_inodo = sizeInodos;
	SB.SB_size_struct_bloque = sizeBloques;
	SB.SB_free_first_bit_struct_arbol_directorio = SB.SB_ap_bitmap_arbol_directorio;
	SB.SB_free_first_bit_struct_detalle_directorio = SB.SB_ap_bitmap_detalle_directorio;
	SB.SB_free_first_bit_struct_inodo = SB.SB_ap_bitmap_tabla_inodo;
	SB.SB_free_first_bit_struct_bloque = SB.SB_ap_bitmap_bloques;
	SB.SB_magic_num = 201801237;

	color.Magenta("╔══════════════════════════════════════════════════╗")
	color.Magenta("                    SUPERBLOQUE")
	color.Magenta("mAGIC NUM: 201801237")
	color.Magenta("CANTIDAD DE ESTRUCTURAS AVD: "+ strconv.FormatInt(cantidadAVD, 10))
	color.Magenta("CANTIDAD DE ESTRUCTURAS DD: "+ strconv.FormatInt(cantidadDD, 10))
	color.Magenta("CANTIDAD DE ESTRUCTURAS INODO: "+ strconv.FormatInt(cantidadInodos, 10))
	color.Magenta("CANTIDAD DE ESTRUCTURAS BLOQUE: "+ strconv.FormatInt(cantidadBloques, 10))
	color.Magenta("CANTIDAD DE ESTRUCTURAS DISPONIBLES AVD: "+ strconv.FormatInt(cantidadAVD, 10))
	color.Magenta("CANTIDAD DE ESTRUCTURAS DISPONIBLES DD: "+ strconv.FormatInt(cantidadDD, 10))
	color.Magenta("CANTIDAD DE ESTRUCTURAS DISPONIBLES INODO: "+ strconv.FormatInt(cantidadInodos, 10))
	color.Magenta("CANTIDAD DE ESTRUCTURAS DISPONIBLES BLOQUE: "+ strconv.FormatInt(cantidadBloques, 10))
	color.Magenta("INICIO BITMAP AVD: "+ strconv.FormatInt(Inicio_bitmapAVD, 10))
	color.Magenta("AVD: "+ strconv.FormatInt(Inicio_AVD, 10))
	color.Magenta("INICIO BITMAP DD: "+ strconv.FormatInt(Inicio_bitmapDD, 10))
	color.Magenta("DD: "+ strconv.FormatInt(Inicio_DD, 10))
	color.Magenta("INICIO BITMAP INODO: "+ strconv.FormatInt(Inicio_bitmapInodo, 10))
	color.Magenta("INODO: "+ strconv.FormatInt(Inicio_Inodos, 10))
	color.Magenta("INICIO BITMAP BLOQUE: "+ strconv.FormatInt(Inicio_bitmapBloque, 10))
	color.Magenta("BLOQUE: "+ strconv.FormatInt(Inicio_Bloque, 10))
	color.Magenta("BITACORA: "+ strconv.FormatInt(Inicio_Bitacora, 10))
	color.Magenta("SB COPIA: "+ strconv.FormatInt(Inicio_SBCopia, 10))
	color.Magenta("╚══════════════════════════════════════════════════╝")

	var p = SB;
	a := &p;
	var binarioP bytes.Buffer;
	binary.Write(&binarioP, binary.BigEndian, a);
	file.Write(binarioP.Bytes());

	/**
	* ESCRITURA BITMAP AVD
	*/
	for i := Inicio_bitmapAVD; i < Inicio_AVD; i++ {
		var init int8 = '0'
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escrituraBytes(file, binarioTemp.Bytes())
	}

	/**
	* ESCRITURA AVD
	*/
	for i := Inicio_AVD; i < Inicio_bitmapDD; i++ {
		var init int8 = '0'
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escrituraBytes(file, binarioTemp.Bytes())
	}

	/**
	* ESCRITURA BITMAP DD
	*/
	for i := Inicio_bitmapDD; i < Inicio_DD; i++ {
		var init int8 = '0'
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escrituraBytes(file, binarioTemp.Bytes())
	}

	/**
	* ESCRITURA DD
	*/
	for i := Inicio_DD; i < Inicio_bitmapInodo; i++ {
		var init int8 = '0'
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escrituraBytes(file, binarioTemp.Bytes())
	}

	/**
	* ESCRITURA BITMAP INODO
	*/
	for i := Inicio_bitmapInodo; i < Inicio_Inodos; i++ {
		var init int8 = '0'
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escrituraBytes(file, binarioTemp.Bytes())
	}

	/**
	* ESCRITURA INODO
	*/
	for i := Inicio_Inodos; i < Inicio_bitmapBloque; i++ {
		var init int8 = '0'
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escrituraBytes(file, binarioTemp.Bytes())
	}

	/**
	* ESCRITURA BITMAP BLOQUE
	*/
	for i := Inicio_bitmapBloque; i < Inicio_Bloque; i++ {
		var init int8 = '0'
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escrituraBytes(file, binarioTemp.Bytes())
	}

	/**
	* ESCRITURA BLOQUE
	*/
	for i := Inicio_Bloque; i < Inicio_Bitacora; i++ {
		var init int8 = '0'
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escrituraBytes(file, binarioTemp.Bytes())
	}

	/**
	* ESCRITURA BITACORA
	*/
	for i := Inicio_Bitacora; i < Inicio_SBCopia; i++ {
		var init int8 = '0'
		o := &init
		file.Seek(i, 0)
		var binarioTemp bytes.Buffer
		binary.Write(&binarioTemp, binary.BigEndian, o)
		escrituraBytes(file, binarioTemp.Bytes())
	}

	file.Seek(Inicio_SBCopia, 0)
	var SBCOPY = Modelo.SUPERBOOT{};
	fechaHoras := time.Now();
	copy(SBCOPY.SB_nombre_hd[:], nombre);
	copy(SBCOPY.SB_date_creacion[:], fechaHoras.Format("2000-01-01 00:00:01"))
	copy(SBCOPY.SB_date_ultimo_montaje[:], fechaHoras.Format("2000-01-01 00:00:01"))
	SBCOPY.SB_arbol_virtual_count = cantidadEstructuras;
	SBCOPY.SB_detalle_directorio_count = cantidadDD;
	SBCOPY.SB_inodos_count = cantidadInodos;
	SBCOPY.SB_bloques_count = cantidadBloques;
	SBCOPY.SB_arbol_virtual_free = cantidadEstructuras;
	SBCOPY.SB_detalle_directorio_free = cantidadDD;
	SBCOPY.SB_inodos_free = cantidadInodos;
	SBCOPY.SB_bloques_free = cantidadBloques;
	SBCOPY.SB_montaje_count	= 1;
	SBCOPY.SB_ap_bitmap_arbol_directorio = Inicio_bitmapAVD;
	SBCOPY.SB_ap_arbol_directorio = Inicio_AVD;
	SBCOPY.SB_ap_bitmap_detalle_directorio = Inicio_bitmapDD;
	SBCOPY.SB_ap_detalle_directorio = Inicio_DD;
	SBCOPY.SB_ap_bitmap_tabla_inodo = Inicio_bitmapInodo;
	SBCOPY.SB_ap_tabla_inodo = Inicio_Inodos;
	SBCOPY.SB_ap_bitmap_bloques = Inicio_bitmapBloque;
	SBCOPY.SB_ap_bloques = Inicio_Bloque;
	SBCOPY.SB_ap_log = Inicio_Bitacora;
	SBCOPY.SB_size_struct_arbol_directorio = sizeAVD;
	SBCOPY.SB_size_struct_detalle_directorio = sizeDD;
	SBCOPY.SB_size_struct_inodo = sizeInodos;
	SBCOPY.SB_size_struct_bloque = sizeBloques;
	SBCOPY.SB_free_first_bit_struct_arbol_directorio = SBCOPY.SB_ap_bitmap_arbol_directorio;
	SBCOPY.SB_free_first_bit_struct_detalle_directorio = SBCOPY.SB_ap_bitmap_detalle_directorio;
	SBCOPY.SB_free_first_bit_struct_inodo = SBCOPY.SB_ap_bitmap_tabla_inodo;
	SBCOPY.SB_free_first_bit_struct_bloque = SBCOPY.SB_ap_bitmap_bloques;
	SBCOPY.SB_magic_num = 201801237;

	var q = SBCOPY
	b := &q
	var binarioPQ bytes.Buffer
	binary.Write(&binarioPQ, binary.BigEndian, b)
	file.Write(binarioPQ.Bytes())
	
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