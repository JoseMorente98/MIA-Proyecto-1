package Path

import (
	"os"
	"fmt"
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
			color.Red("   Ha ocurrido un error creando el directorio D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		} else {
			color.Green("╔══════════════════════════════════════════════════╗")
			color.Green("     Ha creado un directorio exitosamente :D")
			color.Green("╚══════════════════════════════════════════════════╝")
		}
	}
	return ubicacion;
}

/**
 * CREAR ARCHIVO
 */
func CrearArchivo(path string, text string) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		defer file.Close()
		if err != nil {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("    Ha ocurrido un error creando el archivo D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}
		l, err := file.WriteString(text)
		if err != nil {
			fmt.Println(err)
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("   Ha ocurrido un error escribiendo el archivo D:")
			color.Red("╚══════════════════════════════════════════════════╝")
			file.Close()
			return
		} else {
			fmt.Println(l)
		}
	} else if !info.IsDir() {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("             El archivo ya existe D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}