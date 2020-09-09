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
func CrearArchivo(ubicacion string, strCadena string) {
	info, err := os.Stat(ubicacion)
	if os.IsNotExist(err) {
		file, err := os.Create(ubicacion)
		defer file.Close()
		if err != nil {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("    Ha ocurrido un error creando el archivo D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}
		l, err := file.WriteString(strCadena)
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

/**
 * CREAR ARCHIVO
 */
func CrearArchivoRep(ubicacion string, strCadena string) {
	info, err := os.Stat(ubicacion)
	if os.IsNotExist(err) {
		file, err := os.Create(ubicacion)
		defer file.Close()
		if err != nil {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("    Ha ocurrido un error creando el archivo D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}
		l, err := file.WriteString(strCadena)
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
		updatedFile(ubicacion, strCadena)
	}
}

func updatedFile(path string, strText string) {
    // Open file using READ & WRITE permission.
    var file, err = os.OpenFile(path, os.O_RDWR, 0644)
    if isError(err) {
        return
    }
    defer file.Close()

    // Write some text line-by-line to file.
    _, err = file.WriteString(strText)
    if isError(err) {
        return
    }

    // Save file changes.
    err = file.Sync()
    if isError(err) {
        return
    }

    //fmt.Println("File Updated Successfully.")
}


func isError(err error) bool {
    if err != nil {
        fmt.Println(err.Error())
    }

    return (err != nil)
}