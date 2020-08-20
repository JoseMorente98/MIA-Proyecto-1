package InterpreteControlador

import (
	"fmt"
	"strings"
	"os"
	"bufio"
)

func Interprete() {
	finalizar := 0
	inicio := "╔══════════════════════════════════════════════════╗\n"
	inicio += "  UNIVERSIDAD DE SAN CARLOS DE GUATEMALA\n"
	inicio += "  MANEJO E IMPLEMENTACIÓN DE ARCHIVOS A-\n"
	inicio += "  JOSÉ RAFAEL MORENTE GONZÁLEZ\n"
	inicio += "  201801237\n"
	inicio += "╚══════════════════════════════════════════════════╝"
	fmt.Println(inicio)

	for finalizar != 1 {
		lecturaBuffer := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		strComando, _ := lecturaBuffer.ReadString('\n')

		if strings.TrimRight(strComando, "\n") == "SALIR" {
			finalizar = 1
		} else {
			if strComando != "" {
				LineaComando(strComando)
			}
		}
	}
    
}

func LineaComando(strEntrada string) {
	var arregloComando []string
	arregloComando = strings.Split(strEntrada, " ")

	var strComando = strings.ToLower(arregloComando[0])

	fmt.Println("comando" + strComando)
	switch strings.TrimRight(strComando, "\n") {
		case "exec":
			fmt.Println(strComando)
		case "mkdisk":
			fmt.Println(strComando)
		case "rmdisk":
			fmt.Println(strComando)
		case "fdisk":
			fmt.Println(strComando)
		case "mount":
			fmt.Println(strComando)
		case "unmount":
			fmt.Println(strComando)
	}
}
