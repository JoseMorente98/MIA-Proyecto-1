package Path

import (
	"fmt"
	"os"
)

/**
 * FUNCION LINEA DE COMANDO
 */
func CrearDirectorio(ubicacion string) {
	if _, err := os.Stat(ubicacion); os.IsNotExist(err) {
		err := os.Mkdir(ubicacion, 0755)
		if err != nil {
			fmt.Println("Ha ocurrido un error D:", err)
		} else {
			fmt.Println("Ha creado un directorio exitosamente :D")
		}
	}
}