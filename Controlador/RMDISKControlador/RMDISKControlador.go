package RMDISKControlador

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/fatih/color"
)

func RMDISK(ubicacion string) {
	if _, err := os.Stat(ubicacion); os.IsNotExist(err) {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("  No existe el archivo o directorio D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	} else {
		var arregloComando []string = strings.Split(ubicacion, "/")
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("  ¿Deseas Eliminar el archivo "+ arregloComando[len(arregloComando)-1] +"?")
		color.Red("╚══════════════════════════════════════════════════╝")
		fmt.Print("ESCRIBA S/N ")
		reader := bufio.NewReader(os.Stdin)
		comando, _ := reader.ReadString('\n')

		switch strings.TrimRight(comando, "\n") {
		case "S":
			err := os.Remove(ubicacion)

			if err != nil {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("  Ha ocurrido un error D:")
				color.Red("╚══════════════════════════════════════════════════╝")
				return
			}
			color.Green("╔══════════════════════════════════════════════════╗")
			color.Green("  Disco "+ arregloComando[len(arregloComando)-1]+ " eliminado exitosamente D:")
			color.Green("╚══════════════════════════════════════════════════╝")
		case "s":
			err := os.Remove(ubicacion)

			if err != nil {
				color.Red("╔══════════════════════════════════════════════════╗")
				color.Red("  Ha ocurrido un error D:")
				color.Red("╚══════════════════════════════════════════════════╝")
				return
			}
			color.Green("╔══════════════════════════════════════════════════╗")
			color.Green("  Disco "+ arregloComando[len(arregloComando)-1]+ " eliminado exitosamente D:")
			color.Green("╚══════════════════════════════════════════════════╝")
		default: 
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("  Eliminación de disco no autorizado D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}
	}
}
