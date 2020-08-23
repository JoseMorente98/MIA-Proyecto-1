package InterpreteControlador

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	
	//ComandoControlador "../ComandoControlador"
)

/**
 * FUNCION INTERPRETE
 */
func Interprete() {
	finalizar := 0
	inicio := "╔══════════════════════════════════════════════════╗\n"
	inicio += "  UNIVERSIDAD DE SAN CARLOS DE GUATEMALA\n"
	inicio += "  MANEJO E IMPLEMENTACIÓN DE ARCHIVOS A-\n"
	inicio += "  JOSÉ RAFAEL MORENTE GONZÁLEZ\n"
	inicio += "  201801237\n"
	inicio += "╠══════════════════════════════════════════════════╣\n"
	inicio += "  Escriba SALIR para finalizar.\n"
	inicio += "╚══════════════════════════════════════════════════╝"
	fmt.Println(inicio)
	var comando string = "";

	for finalizar != 1 {
		lecturaBuffer := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		strComando, _ := lecturaBuffer.ReadString('\n')

		if strings.TrimRight(strComando, "\n") == "SALIR" {
			finalizar = 1
		} else {
			if strComando != "" {
				if strings.Contains(strComando, "/*") {
					remover := strings.Replace(strComando, "/*", "", 1) 
					comando += strings.TrimRight(remover, "\n")
				} else {
					if comando != "" {
						comando += strComando
						//fmt.Print("COMANDO: " + comando)
						LineaComando(comando)
						comando = ""
					} else {
						//fmt.Print("COMANDO NUEVO: " + strComando)
						LineaComando(strComando)
					}
				}

			}
		}
	}
}

/**
 * FUNCION LINEA DE COMANDO
 */
func LineaComando(strEntrada string) {
	var arregloComando []string
	arregloComando = strings.Split(strEntrada, " ")

	var strComando = strings.ToLower(arregloComando[0])

	if(strings.TrimRight(strComando, "\n") == "exec") {
		ComandoControlador.ComandoExec(arregloComando);
	}
	if(strings.TrimRight(strComando, "\n") == "mkdisk") {
		fmt.Println(strComando)
	}
	if(strings.TrimRight(strComando, "\n") == "rmdisk") {
		fmt.Println(strComando)
	}
	if(strings.TrimRight(strComando, "\n") == "fdisk") {
		fmt.Println(strComando)
	}
	if(strings.TrimRight(strComando, "\n") == "mount") {
		fmt.Println(strComando)
	}
	if(strings.TrimRight(strComando, "\n") == "unmount") {
		fmt.Println(strComando)
	}
}
