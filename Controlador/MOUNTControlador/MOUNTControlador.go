package MOUNTControlador

import (
	"fmt"
	"strconv"
	"github.com/fatih/color"
	"os"

	Modelo "../../Modelo"
)

var montarDiscos [100][26]Modelo.MOUNT

func MOUNT(path string, name string)  {
	if path != "" && name != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("  No existe el archivo o directorio D:")
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
	for i := 0; i < 100; i++ {
		for j := 0; j < 26; j++ {
			if(montarDiscos[i][j] == Modelo.MOUNT{}) {
				s := strconv.Itoa(i+1)
				mount := Modelo.MOUNT{}
				mount.Mount_id = "vd" + CaracterASCII(97+j) + s
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