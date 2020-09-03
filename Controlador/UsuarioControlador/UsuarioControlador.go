package UsuarioControlador

import (
	"github.com/fatih/color"
	"os"
	"bufio"
	"strings"

	Modelo "../../Modelo"
	MOUNTControlador "../MOUNTControlador"
)

var usuario Modelo.USUARIO = Modelo.USUARIO{};

func LOGOUT()  {
	if (usuario == Modelo.USUARIO{}) {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("             No hay sesión activa D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	} else {
		color.Green("╔══════════════════════════════════════════════════╗")
		color.Green("                Hasta la próxima!")
		color.Green("╚══════════════════════════════════════════════════╝")
		usuario = Modelo.USUARIO{}
	}
}

func LOGIN(username string, password string, id string)  {
	if username != "" && password != "" && id != "" {
		autenticar(username, password, id);
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("      Verificar parametros obligatorios D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}

func autenticar(username string, password string, id string)  {
	if (MOUNTControlador.Buscar(id) == true) {
		//mount := Modelo.MOUNT{};
		//mount = MOUNTControlador.BuscarMOUNT(id);
		url := "/home/jmorente/Documentos/CodeApp GT/Archivos/" + "user" + id + ".txt";
		archivo, err := os.OpenFile(url, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		defer archivo.Close()
		if err != nil {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("            No existe archivo TXT D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}

		scanner := bufio.NewScanner(archivo)
		for scanner.Scan() {
			var strUsuario = strings.Split(strings.Trim(scanner.Text(), " "), ",")
			if len(strUsuario) > 3 {
				if (strings.Trim(strUsuario[3], " ") == username && 
					strings.Trim(strUsuario[4], " ") == password) {

					usuario.Usuario_id = id;
					usuario.Usuario_username = strUsuario[3];
					usuario.Usuario_password = strUsuario[4];
					color.Green("╔══════════════════════════════════════════════════╗")
					color.Green("         Haz iniciado sesión " + usuario.Usuario_username+" :D")
					color.Green("╚══════════════════════════════════════════════════╝")
					return;
				}
			}
		}

		if (usuario == Modelo.USUARIO{}) {
			color.Red("╔══════════════════════════════════════════════════╗")
			color.Red("          Usuario o contraseña incorrecto D:")
			color.Red("╚══════════════════════════════════════════════════╝")
		}
	} else {
		color.Red("╔══════════════════════════════════════════════════╗")
		color.Red("     La partición no se encuentra montada D:")
		color.Red("╚══════════════════════════════════════════════════╝")
	}
}