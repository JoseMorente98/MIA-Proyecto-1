package MOUNTControlador

import (
	"fmt"
	"strconv"
	"github.com/fatih/color"
	"os"
	"bytes"
	"unsafe"
	"encoding/binary"
	"log"

	Modelo "../../Modelo"
)

var montarDiscos [100][26]Modelo.MOUNT

func MOUNT(path string, name string)  {
	if path != "" && name != "" {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer file.Close() 
		if err != nil {
			//log.Fatal(err)
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("         No existe el archivo o directorio D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		} else {
			AgregarMontaje(path, name);
		}
	} else {
		Mostrar();
	}
}

func UNMOUNT(id string) {
	EliminarMontaje(id);
}

func Mostrar()  {
	color.Blue("╔══════════════════════════════════════════════════╗")
	color.Blue("            LISTA DE PARTICIONES MONTADAS")
	color.Blue("╚══════════════════════════════════════════════════╝")
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if(montarDiscos[i][j] != Modelo.MOUNT{}) {
				fmt.Print("ID->" + montarDiscos[i][j].Mount_id + 
				" -PATH->" + montarDiscos[i][j].Mount_path + 
				" -NAME->" + montarDiscos[i][j].Mount_particion)
				fmt.Println("")
			}
		}
	}
}

func AgregarMontaje(path string, name string)  {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	defer file.Close()
	if err != nil {
		//log.Fatal(err)
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("         No existe el archivo o directorio D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	} else {
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

		var nombre [16]byte;
		copy(nombre[:], name)

		if (bytes.Compare(nombre[:], m.Mbr_partition_1.Part_name[:]) == 0) {
			AgregarMontajeDD(path, name);
		} else if (bytes.Compare(nombre[:], m.Mbr_partition_2.Part_name[:]) == 0) {
			AgregarMontajeDD(path, name);
		} else if (bytes.Compare(nombre[:], m.Mbr_partition_3.Part_name[:]) == 0) {
			AgregarMontajeDD(path, name);
		} else if (bytes.Compare(nombre[:], m.Mbr_partition_4.Part_name[:]) == 0) {
			AgregarMontajeDD(path, name);
		} else {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("            No existe dicha particion D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}
	}
}

func AgregarMontajeDD(path string, name string)  {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if(montarDiscos[i][j] == Modelo.MOUNT{}) {
				s := strconv.Itoa(i+1)
				mount := Modelo.MOUNT{}
				mount.Mount_id = "vd" + CaracterASCII(97+i) + j
				mount.Mount_path = path
				mount.Mount_particion = name
				montarDiscos[i][j] = mount
				color.Green("╔══════════════════════════════════════════════════╗")
				color.Green("         Partición montada exitosamente :D")
				color.Green("╚══════════════════════════════════════════════════╝")
				return;
			}
		}
	}
}

func EliminarMontaje(id string)  {
	if BuscarMontaje(id) {
		color.Green("╔══════════════════════════════════════════════════╗")
		color.Green("       Partición desmontada exitosamente :D")
		color.Green("╚══════════════════════════════════════════════════╝")
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("             Montaje no encontrado D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func BuscarMontaje(id string) bool {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if(montarDiscos[i][j].Mount_id == id) {
				montarDiscos[i][j] = Modelo.MOUNT{}
				return true;
			}
		}
	}
	return false;
}

func CaracterASCII(numero int) string {
	ascii := numero
	caracter := string(ascii)
	return caracter;
}

func Buscar(id string) bool {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if(montarDiscos[i][j].Mount_id == id) {
				return true;
			}
		}
	}
	return false;
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