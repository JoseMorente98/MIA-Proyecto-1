package Path

import (
	"os"
	"github.com/fatih/color"
)

/**
 * FUNCION LINEA DE COMANDO
 */
func CrearDirectorio(ubicacion string) string {
	if _, err := os.Stat(ubicacion); os.IsNotExist(err) {
		err := os.Mkdir(ubicacion, 0755)
		if err != nil {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("  Ha ocurrido un error creando el directorio D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		} else {
			color.Green("╔══════════════════════════════════════════════════╗")
			color.Green("  Ha creado un directorio exitosamente :D")
			color.Green("╚══════════════════════════════════════════════════╝")
		}
	}
	return ubicacion;
}