package ComandoControlador

import (
	"fmt"
	"strings"
)

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
}

/**
 * FUNCION COMANDO MKDISK
 */
func ComandoMKDISK(strComando []string) {
	size := "";
	path := "";
	name := "";
	unit := "";
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		switch strParametro {
		case "-size":
			size = arregloComando[1]
		case "-path":
			contadorComilla := strings.Count(arregloComando[1], "\"")
			contadorIgual := strings.Count(arregloComando[1], "=")
			removerComilla := strings.Replace(arregloComando[1], "\"", "", contadorComilla)
			removerIgual := strings.Replace(removerComilla, "=", " ", contadorIgual)
			path = removerIgual
		case "-name":
			name = arregloComando[1]
		case "-unit":
			unit = arregloComando[1]
		}
	}

	fmt.Println("SIZE " + size);
	fmt.Println("PATH " + path);
	fmt.Println("NAME " +name);
	fmt.Println("UNIT " + unit);
}

/**
 * FUNCION COMANDO RMDISK
 */
func ComandoRMDISK(strComando []string) {
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
}

/**
 * FUNCION COMANDO FDISK
 */
 func ComandoFDISK(strComando []string) {
	size := "";
	path := "";
	name := "";
	unit := "";
	types := "";
	fit := "";
	delete := "";
	add := "";
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		switch strParametro {
		case "-size":
			size = arregloComando[1]
		case "-path":
			contadorComilla := strings.Count(arregloComando[1], "\"")
			contadorIgual := strings.Count(arregloComando[1], "=")
			removerComilla := strings.Replace(arregloComando[1], "\"", "", contadorComilla)
			removerIgual := strings.Replace(removerComilla, "=", " ", contadorIgual)
			path = removerIgual
		case "-name":
			name = arregloComando[1]
		case "-unit":
			unit = arregloComando[1]
		case "-type":
			types = arregloComando[1]
		case "-fit":
			fit = arregloComando[1]
		case "-delete":
			delete = arregloComando[1]
		case "-add":
			add = arregloComando[1]
		}
	}

	fmt.Println("SIZE " + size);
	fmt.Println("PATH " + path);
	fmt.Println("NAME " +name);
	fmt.Println("UNIT " + unit);
	fmt.Println("TYPE " + types);
	fmt.Println("FIT " + fit);
	fmt.Println("DELETE " + delete);
	fmt.Println("ADD " + add);
}

/**
 * FUNCION COMANDO MOUNT
 */
 func ComandoMOUNT(strComando []string) {
	path := "";
	name := "";
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
		case "-name":
			name = arregloComando[1]
		}
	}

	fmt.Println("PATH " + path);
	fmt.Println("NAME " +name);
}

/**
 * FUNCION COMANDO UNMOUNT
 */
 func ComandoUNMOUNT(strComando []string) {
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		if strings.Contains(strParametro, "-id") {
			fmt.Println("ID: " + arregloComando[1]);
		}
	}
}