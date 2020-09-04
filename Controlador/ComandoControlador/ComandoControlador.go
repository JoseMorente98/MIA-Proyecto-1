package ComandoControlador

import (
	"fmt"
	"strings"

	MBRControlador "../MBRControlador"
	RMDISKControlador "../RMDISKControlador"
	MOUNTControlador "../MOUNTControlador"
	FDISKControlador "../FDISKControlador"
	MKFSControlador "../MKFSControlador"
	UsuarioControlador "../UsuarioControlador"
	ReporteControlador "../ReporteControlador"
)

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

	MBRControlador.MKDISK(size, path, name, unit);
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
	RMDISKControlador.RMDISK(path);
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

	FDISKControlador.FDISK(size, path, name, unit, types, fit, add, delete);
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

	MOUNTControlador.MOUNT(path, name);
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
			MOUNTControlador.UNMOUNT(arregloComando[1]);
		}
	}
}

/**
 * FUNCION COMANDO MKFS
 */
 func ComandoMKFS(strComando []string) {
	id := "";
	unit := "";
	types := "";
	add := "";
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		switch strParametro {
		case "-id":
			id = arregloComando[1]
		case "-unit":
			unit = arregloComando[1]
		case "-type":
			types = arregloComando[1]
		case "-add":
			add = arregloComando[1]
		}
	}

	fmt.Println("ID " + id);
	fmt.Println("UNIT " + unit);
	fmt.Println("TYPE " + types);
	fmt.Println("ADD " + add);

	MKFSControlador.MKFS(id, types, add, unit);
}

/**
 * FUNCION COMANDO LOGIN
 */
 func ComandoLOGIN(strComando []string) {
	id := "";
	user := "";
	password := "";
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		switch strParametro {
		case "-id":
			id = arregloComando[1]
		case "-usr":
			user = arregloComando[1]
		case "-pwd":
			password = arregloComando[1]
		}
	}

	fmt.Println("ID " + id);
	fmt.Println("USER " + user);
	fmt.Println("PASSWORD " + password);
	//FDISKControlador.FDISK(size, path, name, unit, types, fit, add, delete);
	UsuarioControlador.LOGIN(user, password, id);
}

/**
 * FUNCION COMANDO MKGRP
 */
 func ComandoMKGRP(strComando []string) {
	id := "";
	name := "";
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		switch strParametro {
		case "-id":
			id = arregloComando[1]
		case "-name":
			name = arregloComando[1]
		}
	}

	fmt.Println("ID " + id);
	fmt.Println("NAME " + name);
	//FDISKControlador.FDISK(size, path, name, unit, types, fit, add, delete);
}

/**
 * FUNCION COMANDO RMGRP
 */
 func ComandoRMGRP(strComando []string) {
	id := "";
	name := "";
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		switch strParametro {
		case "-id":
			id = arregloComando[1]
		case "-name":
			name = arregloComando[1]
		}
	}

	fmt.Println("ID " + id);
	fmt.Println("NAME " + name);
	//FDISKControlador.FDISK(size, path, name, unit, types, fit, add, delete);
}

/**
 * FUNCION COMANDO MKUSR
 */
 func ComandoMKUSR(strComando []string) {
	id := "";
	user := "";
	password := "";
	group := "";
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		switch strParametro {
		case "-id":
			id = arregloComando[1]
		case "-usr":
			user = arregloComando[1]
		case "-pwd":
			password = arregloComando[1]
		case "-grp":
			group = arregloComando[1]
		}
	}

	fmt.Println("ID " + id);
	fmt.Println("USER " + user);
	fmt.Println("PASSWORD " + password);
	fmt.Println("GROUP " + group);
	//FDISKControlador.FDISK(size, path, name, unit, types, fit, add, delete);
}


/**
 * FUNCION COMANDO RMUSR
 */
 func ComandoRMUSR(strComando []string) {
	id := "";
	user := "";
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		switch strParametro {
		case "-id":
			id = arregloComando[1]
		case "-usr":
			user = arregloComando[1]
		}
	}

	fmt.Println("ID " + id);
	fmt.Println("USER " + user);
	//FDISKControlador.FDISK(size, path, name, unit, types, fit, add, delete);
}

/**
 * FUNCION COMANDO RMUSR
 */
 func ComandoREP(strComando []string) {
	nombre := "";
	path := "";
	id := "";
	ruta := "";
	for i := 1; i < len(strComando); i++ {
		var arregloComando = strings.Split(strComando[i], "->")
		var strParametro string = strings.ToLower(arregloComando[0])
		switch strParametro {
		case "-id":
			id = arregloComando[1]
		case "-ruta":
			ruta = arregloComando[1]
		case "-path":
			contadorComilla := strings.Count(arregloComando[1], "\"")
			contadorIgual := strings.Count(arregloComando[1], "=")
			removerComilla := strings.Replace(arregloComando[1], "\"", "", contadorComilla)
			removerIgual := strings.Replace(removerComilla, "=", " ", contadorIgual)
			path = removerIgual
		case "-nombre":
			nombre = arregloComando[1]
		}
	}
	ReporteControlador.REP(nombre, path, id, ruta)
	//FDISKControlador.FDISK(size, path, name, unit, types, fit, add, delete);
}