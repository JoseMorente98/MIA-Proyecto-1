package ControladorMBR

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"github.com/fatih/color"
	"log"
	"os"
	"strconv"
	"time"
	"math/rand"

	Modelo "../../Modelo"
	Path "../../Utilidades/Path"
)

/**
 * FUNCION MKDISK
 */
func MKDISK(size string, path string, name string, unit string) {
	if size != "" && path != "" && name != "" {
		if unit != "" {
			if strings.ToLower(unit) == "k" || strings.ToLower(unit) == "m" {
				s := size
				n, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					color.Red("╔══════════════════════════════════════════════════╗")
					color.Red("  No se puede convertir string a int64 D:")
					color.Red("╚══════════════════════════════════════════════════╝")
				}
				CrearArchivo(n, path, name, unit);
			} else {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("  Solo puedes utilizar m o k D:")
				color.Red("╚══════════════════════════════════════════════════╝")
			}
		} else {
			s := size
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("  No se puede convertir string a int64 D:")
				color.Red("╚══════════════════════════════════════════════════╝")
			}
			CrearArchivo(n, path, name, "m");
		}		
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("  Verificar parametros obligatorios D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

/**
 * CREAR ARCHIVO BINARIO
 */
func CrearArchivo(size int64, path string, name string, unit string) {
	var sizeUNIT int64;
	if size >= 1 {
		if strings.Contains(name, ".dsk") {
			switch strings.ToLower(unit) {
			case "k":
				sizeUNIT = 1024;
			case "m":
				sizeUNIT = 1024 * 1024;
			default:
				sizeUNIT = 1024 * 1024;
			}

			Path.CrearDirectorio(path);
			archivo, err := os.Create(path + name)
			defer archivo.Close()
			if err != nil {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("  Ha ocurrido un error D:")
				color.Red("╚══════════════════════════════════════════════════╝")
			}
			var other int8 = 0

			s := &other

			//ESCRIBIR CERO AL INICIO DEL ARCHIVO
			var binaryTemp bytes.Buffer
			binary.Write(&binaryTemp, binary.BigEndian, s)
			escrituraBytes(archivo, binaryTemp.Bytes())

			//TAMAÑO DE ARCHIVO SEGUN UNIT
			archivo.Seek(sizeUNIT*size-1, 0)
		
			//ESCRIBIR 0 AL FINAL DEL ARCHIVO
			var binaryTemp2 bytes.Buffer
			binary.Write(&binaryTemp2, binary.BigEndian, s)
			escrituraBytes(archivo, binaryTemp2.Bytes())
			archivo.Seek(0, 0)
		
			//ASIGNAR VALORES A STRUCT
			fechaHora := time.Now();
			disco := Modelo.MBR{}
			disco.Mbr_size = size*sizeUNIT-1;
			copy(disco.Mbr_date[:], fechaHora.Format("2000-01-01 00:00:01"))
			disco.Mbr_disk_signature = rand.Int63();
			disco.Mbr_partition_1 = Modelo.PARTICION{};
			disco.Mbr_partition_2 = Modelo.PARTICION{};
			disco.Mbr_partition_3 = Modelo.PARTICION{};
			disco.Mbr_partition_4 = Modelo.PARTICION{};
			fmt.Println(disco)
			s1 := &disco

		
			//ESCRITURA DEL STRUCT
			var binaryTemp3 bytes.Buffer
			binary.Write(&binaryTemp3, binary.BigEndian, s1)
			escrituraBytes(archivo, binaryTemp3.Bytes())
			color.Green("╔══════════════════════════════════════════════════╗")
			color.Green("  Disco creado exitosamente :D")
			color.Green("╚══════════════════════════════════════════════════╝")
		} else {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("  La extensión del archivo debe de ser .dsk")
			color.Red("╚══════════════════════════════════════════════════╝")
		}
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("  El tamaño no es permitido D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
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