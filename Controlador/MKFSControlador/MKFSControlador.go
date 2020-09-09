package MKFSControlador

import (
	"strings"
	"github.com/fatih/color"

	Modelo "../../Modelo"
	Path "../../Utilidades/Path"
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
		Path.CrearArchivo("/home/jmorente/Documentos/CodeApp GT/Archivos/" + "user" + string(mount.Mount_id[:]) + ".txt", 
		"1,G,root\n1,U,root,root,201801237\n");
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("     La partición no se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}

}