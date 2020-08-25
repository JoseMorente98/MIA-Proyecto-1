package InterpreteControlador

import (
	"fmt"
	"strings"
	"os"
	"bufio"
  	"log"
	
	ComandoControlador "../ComandoControlador"
	Espacio "../../Utilidades/Espacio"
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
				if strings.Contains(strComando, "\\*") {
					remover := strings.Replace(strComando, "\\*", "", 1)
					comando += strings.TrimRight(remover, "\n")
				} else {
					if comando != "" {
						comandoTrim := strings.TrimRight(strComando, "\n")
						comando += comandoTrim
						//fmt.Print("COMANDO: " + comando)
						LineaComando(comando)
						comando = ""
					} else {
						//fmt.Print("COMANDO NUEVO: " + strComando)
						comando := strings.TrimRight(strComando, "\n")
						LineaComando(comando)
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
	arregloComando = strings.Split(Espacio.RemoverEspacioPath(strEntrada), " ")

	var strComando = strings.ToLower(arregloComando[0])

	switch strings.TrimRight(strComando, "\n") {
	case "exec":
		fmt.Println("╔══════════════════════════════════════════════════╗")
		fmt.Println("                      EXEC")
		fmt.Println("╚══════════════════════════════════════════════════╝")
		ComandoEXEC(arregloComando)
	case "mkdisk":
		fmt.Println("╔══════════════════════════════════════════════════╗")
		fmt.Println("                      MKDISK")
		fmt.Println("╚══════════════════════════════════════════════════╝")
		ComandoControlador.ComandoMKDISK(arregloComando)
	case "rmdisk":
		fmt.Println("╔══════════════════════════════════════════════════╗")
		fmt.Println("                     RMDISK")
		fmt.Println("╚══════════════════════════════════════════════════╝")
		ComandoControlador.ComandoRMDISK(arregloComando)
	case "fdisk":
		fmt.Println("╔══════════════════════════════════════════════════╗")
		fmt.Println("                      FDISK")
		fmt.Println("╚══════════════════════════════════════════════════╝")
		ComandoControlador.ComandoFDISK(arregloComando)
	case "mount":
		fmt.Println("╔══════════════════════════════════════════════════╗")
		fmt.Println("                      MOUNT")
		fmt.Println("╚══════════════════════════════════════════════════╝")
		ComandoControlador.ComandoMOUNT(arregloComando)
	case "unmount":
		fmt.Println("╔══════════════════════════════════════════════════╗")
		fmt.Println("                     UNMOUNT")
		fmt.Println("╚══════════════════════════════════════════════════╝")
		ComandoControlador.ComandoUNMOUNT(arregloComando)
	case "pause":
		fmt.Println("╔══════════════════════════════════════════════════╗")
		fmt.Println("  PRESIONE 'ENTER' PARA CONTINUAR...")
		fmt.Println("╚══════════════════════════════════════════════════╝")
    	bufio.NewReader(os.Stdin).ReadBytes('\n')
	default:
		if !strings.Contains(strComando, "#") {
			fmt.Println("╔══════════════════════════════════════════════════╗")
			fmt.Println("  COMANDO NO SOPORTADO")
			fmt.Println("╚══════════════════════════════════════════════════╝")
		}
	}
}


/**
 * FUNCION COMANDO EXEC
 */
 func ComandoEXEC(strComando []string) {
	path := "";
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		switch strParametro {
		case "-path":
			contadorComilla := strings.Count(arregloComando[1], "\"")
			contadorIgual := strings.Count(arregloComando[1], "=")
			removerComilla := strings.Replace(arregloComando[1], "\"", "", contadorComilla)
			removerIgual := strings.Replace(removerComilla, "=", " ", contadorIgual)
			path = removerIgual
		}
	}
	fmt.Println("PATH " + path);
	//LEER ARCHIVO
	LeerArchivo(path);
}


func LeerArchivo(ubicacion string) {
	file, err := os.Open(ubicacion)
	if err != nil {
	  log.Fatal(err)
	  fmt.Println("HA ocurrido un error");
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	var comando string = "";
	for scanner.Scan() {
	  fmt.Println(scanner.Text())
	  strComando := scanner.Text()
	  	if strComando != "" {
			if strings.Contains(strComando, "\\*") {
				remover := strings.Replace(strComando, "\\*", "", 1)
				comando += strings.TrimRight(remover, "\n")
			} else {
				if comando != "" {
					comandoTrim := strings.TrimRight(strComando, "\n")
					comando += comandoTrim
					LineaComando(comando)
					comando = ""
				} else {
					comando := strings.TrimRight(strComando, "\n")
					LineaComando(comando)
				}
			}
		}
	}
}