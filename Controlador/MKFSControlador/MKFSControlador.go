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
			ParticionLogica(mount.Mount_particion_logica, file, filepath.Base(mount.Mount_path));
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
	copy(SB.SB_date_creacion[:], fechaHora.Format("2006-01-02 15:04:05"))
	copy(SB.SB_date_ultimo_montaje[:], fechaHora.Format("2006-01-02 15:04:05"))
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
	copy(SBCOPY.SB_date_creacion[:], fechaHoras.Format("2006-01-02 15:04:05"))
	copy(SBCOPY.SB_date_ultimo_montaje[:], fechaHoras.Format("2006-01-02 15:04:05"))
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


	/**
	 * OBTENEMOS EL SUPERBOOT
	 */
	SBDIR := Modelo.SUPERBOOT{}
	var sbSize int64 = int64(unsafe.Sizeof(SBDIR))
	file.Seek(particion.Part_start, 0)
	data := leerBytes(file, sbSize)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &SBDIR)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	//CAMBIAMOS EL 0 POR UN 1 EN EL BITMAP AVD
	file.Seek(SBDIR.SB_ap_bitmap_arbol_directorio, 0);
	var unidad int8 = '1'
	s1 := &unidad
	var binarytemp1 bytes.Buffer
	binary.Write(&binarytemp1, binary.BigEndian, s1)
	escrituraBytes(file, binarytemp1.Bytes())

	//INGRESAMOS EL AVD DE LA RUTA /
	file.Seek(SBDIR.SB_ap_arbol_directorio, 0);
	AVD := Modelo.AVD{}
	var fechaAhora = time.Now()
	copy(AVD.Avd_fecha_creacion[:], fechaAhora.Format("2006-01-02 15:04:05"))
	copy(AVD.Avd_nombre_directorio[:], "/")
	AVD.Avd_ap_detalle_directorio = SBDIR.SB_ap_detalle_directorio
	AVD.Avd_ap_arbol_virtual_directorio = -1;
	AVD.Avd_proper = 1;
	for i := 0; i < 6; i++ {
		AVD.Avd_ap_array_subdirectorios[i] = -1;
	}
	avd := &AVD
	var binarytemp2 bytes.Buffer
	binary.Write(&binarytemp2, binary.BigEndian, avd)
	escrituraBytes(file, binarytemp2.Bytes())

	//INGRESAMOS BITMAP DD
	file.Seek(SBDIR.SB_ap_bitmap_detalle_directorio, 0)
	var unidad2 int8 = '1'
	s3 := &unidad2
	var binarytemp3 bytes.Buffer
	binary.Write(&binarytemp3, binary.BigEndian, s3)
	escrituraBytes(file, binarytemp3.Bytes())

	//INGRESAMOS DETALLE DE DIRECTORIO
	file.Seek(SBDIR.SB_ap_detalle_directorio, 0);

	DD := Modelo.DD{}
	DD.DD_ap_detalle_directorio = 0;
	DD.DD_array_files[0].DD_file_ap_inodo = SBDIR.SB_ap_bitmap_tabla_inodo;
	copy(DD.DD_array_files[0].DD_file_nombre[:], "users.txt")
	copy(DD.DD_array_files[0].DD_file_date_creacion[:], fechaAhora.Format("2006-01-02 15:04:05"))
	copy(DD.DD_array_files[0].DD_file_date_modificacion[:], fechaAhora.Format("2006-01-02 15:04:05"))
	DD1 := &DD
	var binaryTemp4 bytes.Buffer
	binary.Write(&binaryTemp4, binary.BigEndian, DD1)
	escrituraBytes(file, binaryTemp4.Bytes())

	//BITMAP INODO
	file.Seek(SBDIR.SB_ap_bitmap_tabla_inodo, 0)
	var unidad3 int8 = '1'
	s9 := &unidad3
	var binarytemp5 bytes.Buffer
	binary.Write(&binarytemp5, binary.BigEndian, s9)
	escrituraBytes(file, binarytemp5.Bytes())

	//INODO
	file.Seek(SBDIR.SB_ap_tabla_inodo, 0);
	INODO := Modelo.INODO{}
	INODO.I_count_inodo = 1;
	INODO.I_size_archivo = 33;
	INODO.I_ap_indirecto = -1;
	INODO.I_id_proper = 1;
	INODO.I_count_bloques_asignados = 2;
	INODO.I_array_bloques[0] = SBDIR.SB_ap_bloques
	INODO.I_array_bloques[1] = SBDIR.SB_ap_bloques + sizeSB;
	INODO.I_array_bloques[2] = 0
	INODO.I_array_bloques[3] = 0

	inodo1 := &INODO
	var binaryTemp6 bytes.Buffer
	binary.Write(&binaryTemp6, binary.BigEndian, inodo1)
	escrituraBytes(file, binaryTemp6.Bytes())

	//INICIO DE BITMAP BLOQUE
	file.Seek(SBDIR.SB_ap_bitmap_bloques, 0)
	var unidad4 int8 = '1'
	s4 := &unidad4
	var binarytemp7 bytes.Buffer
	binary.Write(&binarytemp7, binary.BigEndian, s4)
	escrituraBytes(file, binarytemp7.Bytes())
	file.Seek(SBDIR.SB_ap_bloques+1, 0)
	escrituraBytes(file, binarytemp7.Bytes())

	//BLOQUES
	file.Seek(SBDIR.SB_ap_bloques, 0);

	B1 := Modelo.BLOQUE{};
	B2 := Modelo.BLOQUE{};
	copy(B1.DB_data[:], "1,G,root\n1,U,root,root,20")
	copy(B2.DB_data[:], "1801237\n")
	bloque1 := &B1
	var binaryTemp8 bytes.Buffer
	binary.Write(&binaryTemp8, binary.BigEndian, bloque1)
	escrituraBytes(file, binaryTemp8.Bytes())

	file.Seek(SBDIR.SB_ap_bloques + int64(unsafe.Sizeof(Modelo.BLOQUE{})), 0);
	bloque2 := &B2
	var binaryTemp9 bytes.Buffer
	binary.Write(&binaryTemp9, binary.BigEndian, bloque2)
	escrituraBytes(file, binaryTemp9.Bytes())

}

func ParticionLogica(particion Modelo.PARTICION_LOGICA, file *os.File, nombre string)  {
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
	copy(SB.SB_date_creacion[:], fechaHora.Format("2006-01-02 15:04:05"))
	copy(SB.SB_date_ultimo_montaje[:], fechaHora.Format("2006-01-02 15:04:05"))
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
	copy(SBCOPY.SB_date_creacion[:], fechaHoras.Format("2006-01-02 15:04:05"))
	copy(SBCOPY.SB_date_ultimo_montaje[:], fechaHoras.Format("2006-01-02 15:04:05"))
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


	/**
	 * OBTENEMOS EL SUPERBOOT
	 */
	SBDIR := Modelo.SUPERBOOT{}
	var sbSize int64 = int64(unsafe.Sizeof(SBDIR))
	file.Seek(particion.Part_start, 0)
	data := leerBytes(file, sbSize)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &SBDIR)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	//CAMBIAMOS EL 0 POR UN 1 EN EL BITMAP AVD
	file.Seek(SBDIR.SB_ap_bitmap_arbol_directorio, 0);
	var unidad int8 = '1'
	s1 := &unidad
	var binarytemp1 bytes.Buffer
	binary.Write(&binarytemp1, binary.BigEndian, s1)
	escrituraBytes(file, binarytemp1.Bytes())

	//INGRESAMOS EL AVD DE LA RUTA /
	file.Seek(SBDIR.SB_ap_arbol_directorio, 0);
	AVD := Modelo.AVD{}
	var fechaAhora = time.Now()
	copy(AVD.Avd_fecha_creacion[:], fechaAhora.Format("2006-01-02 15:04:05"))
	copy(AVD.Avd_nombre_directorio[:], "/")
	AVD.Avd_ap_detalle_directorio = SBDIR.SB_ap_detalle_directorio
	AVD.Avd_ap_arbol_virtual_directorio = -1;
	AVD.Avd_proper = 1;
	for i := 0; i < 6; i++ {
		AVD.Avd_ap_array_subdirectorios[i] = -1;
	}
	avd := &AVD
	var binarytemp2 bytes.Buffer
	binary.Write(&binarytemp2, binary.BigEndian, avd)
	escrituraBytes(file, binarytemp2.Bytes())

	//INGRESAMOS BITMAP DD
	file.Seek(SBDIR.SB_ap_bitmap_detalle_directorio, 0)
	var unidad2 int8 = '1'
	s3 := &unidad2
	var binarytemp3 bytes.Buffer
	binary.Write(&binarytemp3, binary.BigEndian, s3)
	escrituraBytes(file, binarytemp3.Bytes())

	//INGRESAMOS DETALLE DE DIRECTORIO
	file.Seek(SBDIR.SB_ap_detalle_directorio, 0);

	DD := Modelo.DD{}
	DD.DD_ap_detalle_directorio = 0;
	DD.DD_array_files[0].DD_file_ap_inodo = SBDIR.SB_ap_bitmap_tabla_inodo;
	copy(DD.DD_array_files[0].DD_file_nombre[:], "users.txt")
	copy(DD.DD_array_files[0].DD_file_date_creacion[:], fechaAhora.Format("2006-01-02 15:04:05"))
	copy(DD.DD_array_files[0].DD_file_date_modificacion[:], fechaAhora.Format("2006-01-02 15:04:05"))
	DD1 := &DD
	var binaryTemp4 bytes.Buffer
	binary.Write(&binaryTemp4, binary.BigEndian, DD1)
	escrituraBytes(file, binaryTemp4.Bytes())

	//BITMAP INODO
	file.Seek(SBDIR.SB_ap_bitmap_tabla_inodo, 0)
	var unidad3 int8 = '1'
	s9 := &unidad3
	var binarytemp5 bytes.Buffer
	binary.Write(&binarytemp5, binary.BigEndian, s9)
	escrituraBytes(file, binarytemp5.Bytes())

	//INODO
	file.Seek(SBDIR.SB_ap_tabla_inodo, 0);
	INODO := Modelo.INODO{}
	INODO.I_count_inodo = 1;
	INODO.I_size_archivo = 33;
	INODO.I_ap_indirecto = -1;
	INODO.I_id_proper = 1;
	INODO.I_count_bloques_asignados = 2;
	INODO.I_array_bloques[0] = SBDIR.SB_ap_bloques
	INODO.I_array_bloques[1] = SBDIR.SB_ap_bloques + sizeSB;
	INODO.I_array_bloques[2] = 0
	INODO.I_array_bloques[3] = 0

	inodo1 := &INODO
	var binaryTemp6 bytes.Buffer
	binary.Write(&binaryTemp6, binary.BigEndian, inodo1)
	escrituraBytes(file, binaryTemp6.Bytes())

	//INICIO DE BITMAP BLOQUE
	file.Seek(SBDIR.SB_ap_bitmap_bloques, 0)
	var unidad4 int8 = '1'
	s4 := &unidad4
	var binarytemp7 bytes.Buffer
	binary.Write(&binarytemp7, binary.BigEndian, s4)
	escrituraBytes(file, binarytemp7.Bytes())
	file.Seek(SBDIR.SB_ap_bloques+1, 0)
	escrituraBytes(file, binarytemp7.Bytes())

	//BLOQUES
	file.Seek(SBDIR.SB_ap_bloques, 0);

	B1 := Modelo.BLOQUE{};
	B2 := Modelo.BLOQUE{};
	copy(B1.DB_data[:], "1,G,root\n1,U,root,root,20")
	copy(B2.DB_data[:], "1801237\n")
	bloque1 := &B1
	var binaryTemp8 bytes.Buffer
	binary.Write(&binaryTemp8, binary.BigEndian, bloque1)
	escrituraBytes(file, binaryTemp8.Bytes())

	file.Seek(SBDIR.SB_ap_bloques + int64(unsafe.Sizeof(Modelo.BLOQUE{})), 0);
	bloque2 := &B2
	var binaryTemp9 bytes.Buffer
	binary.Write(&binaryTemp9, binary.BigEndian, bloque2)
	escrituraBytes(file, binaryTemp9.Bytes())

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
 * LECTURA DE BYTES
 */
 func leerBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number) //array de bytes
	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
