package Espacio

import (
	"strings"
)

func RemoverEspacioPath(strComando string) string {
	if strings.ContainsAny(strComando, "\"") {
		strCadena := ""

		for x := 0; x < len(strComando); x++ {
			strCadena += strComando[x:x+1]

			if strComando[x:x+1] == "\"" {

				for y := x + 1; y < len(strComando); y++ {

					if strComando[y:y+1] == " " {
						strCadena += "="
					} else {
						strCadena += strComando[y:y+1]
					}

					if strComando[y:y+1] == "\"" {
						x = y
						break
					}
				}
			}
		}
		return strCadena
	}
	return strComando
}