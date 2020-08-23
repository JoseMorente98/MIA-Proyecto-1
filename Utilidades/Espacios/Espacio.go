package Espacio

import (
	"fmt"
	"strings"
)

func RemoverEspacioComillas(strComando string) string {
	if strings.ContainsAny(strComando, "\"") {
		var nuevaCadena string = ""

		for i := 0; i < len(strComando); i++ {
			nuevaCadena = nuevaCadena + string(strComando[i])

			if string(strComando[i]) == "\"" {

				for j := i + 1; j < len(strComando); j++ {

					if string(strComando[j]) == " " {
						nuevaCadena = nuevaCadena + "*"
					} else {
						nuevaCadena = nuevaCadena + string(strComando[j])
					}

					if string(strComando[j]) == "\"" {
						i = j
						break
					}
				}
			}
		}
		return nuevaCadena
	}
	return strComando
}