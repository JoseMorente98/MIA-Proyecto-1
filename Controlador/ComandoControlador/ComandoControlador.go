package ComandoControlador

import (
	"fmt"
	"strings"
)


func ComandoExec(strComando string) {
	var arregloComando []string
	arregloComando = strings.Split(strComando, "->")
	comando := ""
	ubicacion := ""
	for i := 0; i < len(arregloComando); i++ {
		switch i {
			case 0:
				comando = arregloComando[0]
			case 1:
				contadorComilla := strings.Count(arregloComando[1], "\"")
				contadorIgual := strings.Count(arregloComando[1], "=")
				removerComilla := strings.Replace(arregloComando[1], "\"", "", contadorComilla)
				removerIgual := strings.Replace(removerComilla, "=", " ", contadorIgual)
				ubicacion = removerIgual
		}
	}

	fmt.Println("Comando " + comando)
	fmt.Println("Ubicacion " + ubicacion)
}