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
		//color.Red(string(nombre[:]))
		if (bytes.Compare(nombre[:], m.Mbr_partition_1.Part_name[:]) == 0) {
			if m.Mbr_partition_1.Part_type == 'P' {
				AgregarMontajeDDP(path, name, m.Mbr_partition_1);
				return;
			}
		} else if (bytes.Compare(nombre[:], m.Mbr_partition_2.Part_name[:]) == 0) {
			if m.Mbr_partition_2.Part_type == 'P' {
				AgregarMontajeDDP(path, name, m.Mbr_partition_2);
				return;
			}
		} else if (bytes.Compare(nombre[:], m.Mbr_partition_3.Part_name[:]) == 0) {
			if m.Mbr_partition_3.Part_type == 'P' {
				AgregarMontajeDDP(path, name, m.Mbr_partition_3);
				return;
			}
		} else if (bytes.Compare(nombre[:], m.Mbr_partition_4.Part_name[:]) == 0) {
			if m.Mbr_partition_4.Part_type == 'P' {
				AgregarMontajeDDP(path, name, m.Mbr_partition_4);
				return;
			}
		}

		for i := 0; i < 50; i++ {
			//color.Red(string(nombre[:]) + " -> " + string(m.Mbr_partition_1.Part_EBR[i].Part_name[:]))
			if (bytes.Compare(nombre[:], m.Mbr_partition_1.Part_EBR[i].Part_name[:]) == 0) {
				AgregarMontajeDDPL(path, name, m.Mbr_partition_1.Part_EBR[i].Part_logica);
				return;
			}
		}

		for i := 0; i < 50; i++ {
			//color.Red(string(nombre[:]) + " -> " + string(m.Mbr_partition_2.Part_EBR[i].Part_name[:]))
			if (bytes.Compare(nombre[:], m.Mbr_partition_2.Part_EBR[i].Part_name[:]) == 0) {
				AgregarMontajeDDPL(path, name, m.Mbr_partition_2.Part_EBR[i].Part_logica);
				return;
			}
		}

		for i := 0; i < 50; i++ {
			//color.Red(string(nombre[:]) + " -> " + string(m.Mbr_partition_3.Part_EBR[i].Part_name[:]))
			if (bytes.Compare(nombre[:], m.Mbr_partition_3.Part_EBR[i].Part_name[:]) == 0) {
				AgregarMontajeDDPL(path, name, m.Mbr_partition_3.Part_EBR[i].Part_logica);
				return;
			}
		}

		for i := 0; i < 50; i++ {
			//color.Red(string(nombre[:]) + " -> " + string(m.Mbr_partition_4.Part_EBR[i].Part_name[:]))
			if (bytes.Compare(nombre[:], m.Mbr_partition_4.Part_EBR[i].Part_name[:]) == 0) {
				AgregarMontajeDDPL(path, name, m.Mbr_partition_4.Part_EBR[i].Part_logica);
				return;
			}
		}
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("            No existe dicha particion D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func AgregarMontajeDDP(path string, name string, particion Modelo.PARTICION)  {
	if BuscarParticion(path, name) {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("       La partición ya se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	} else {
		for j := 0; j < 26; j++ {
			if(montarDiscos[0][j].Mount_path == path) {
				//INCREMENTAR SOLO I PARA AGREGAR
				for i := 1; i < 100; i++ {
					if(montarDiscos[i][j] == Modelo.MOUNT{}) {
						s := strconv.Itoa(i+1)
						mount := Modelo.MOUNT{}
						mount.Mount_id = "vd" + CaracterASCII(97+j) + s
						mount.Mount_path = path
						mount.Mount_particion = name
						mount.Mount_partition = particion
						montarDiscos[i][j] = mount
						color.Green("╔══════════════════════════════════════════════════╗")
						color.Green("         Partición montada exitosamente :D")
						color.Green("╚══════════════════════════════════════════════════╝")
						return;
					}
				}
			} else {
				if(montarDiscos[0][j] == Modelo.MOUNT{}) {
					mount := Modelo.MOUNT{}
					mount.Mount_id = "vd" + CaracterASCII(97+j) + "1"
					mount.Mount_path = path
					mount.Mount_particion = name
					mount.Mount_partition = particion
					montarDiscos[0][j] = mount
					color.Green("╔══════════════════════════════════════════════════╗")
					color.Green("         Partición montada exitosamente :D")
					color.Green("╚══════════════════════════════════════════════════╝")
					return;
				}
			}
		}
	}
}

func AgregarMontajeDDPL(path string, name string, particion Modelo.PARTICION_LOGICA)  {
	if BuscarParticion(path, name) {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("       La partición ya se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	} else {
		for j := 0; j < 26; j++ {
			if(montarDiscos[0][j].Mount_path == path) {
				//INCREMENTAR SOLO I PARA AGREGAR
				for i := 1; i < 100; i++ {
					if(montarDiscos[i][j] == Modelo.MOUNT{}) {
						s := strconv.Itoa(i+1)
						mount := Modelo.MOUNT{}
						mount.Mount_id = "vd" + CaracterASCII(97+j) + s
						mount.Mount_path = path
						mount.Mount_particion = name
						mount.Mount_particion_logica = particion
						montarDiscos[i][j] = mount
						color.Green("╔══════════════════════════════════════════════════╗")
						color.Green("     Partición logica montada exitosamente :D")
						color.Green("╚══════════════════════════════════════════════════╝")
						return;
					}
				}
			} else {
				if(montarDiscos[0][j] == Modelo.MOUNT{}) {
					mount := Modelo.MOUNT{}
					mount.Mount_id = "vd" + CaracterASCII(97+j) + "1"
					mount.Mount_path = path
					mount.Mount_particion = name
					mount.Mount_particion_logica = particion
					montarDiscos[0][j] = mount
					color.Green("╔══════════════════════════════════════════════════╗")
					color.Green("      Partición logica montada exitosamente :D")
					color.Green("╚══════════════════════════════════════════════════╝")
					return;
				}
			}
		}
	}
}

func verificarPath(path string) bool {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if(montarDiscos[0][j].Mount_path == path) {
				return true;
			}
		}
	}
	return false;
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

func BuscarParticion(path string, name string) bool {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if(montarDiscos[i][j].Mount_particion == name && montarDiscos[i][j].Mount_path == path) {
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

func BuscarMOUNT(id string) Modelo.MOUNT {
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if(montarDiscos[i][j].Mount_id == id) {
				return montarDiscos[i][j];
			}
		}
	}
	return Modelo.MOUNT{};
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