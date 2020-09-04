package main

import (
	InterpreteControlador "./Controlador/InterpreteControlador"
	//MOUNTControlador "./Controlador/MOUNTControlador"
	///ReporteControlador "./Controlador/ReporteControlador"
)

func main() {
	//MOUNTControlador.AgregarMontaje("/home/jmorente/Documentos/Manejo e Implementación de Archivos/DiscoTest2.dsk", "Particion2");

	InterpreteControlador.Interprete()
	/*MOUNTControlador.AgregarMontaje("/home/jmorente/Documentos/Manejo e Implementación de Archivos/DiscoTest.dsk", "Particion1");
	MOUNTControlador.AgregarMontaje("/home/jmorente/Documentos/Manejo e Implementación de Archivos/DiscoTest.dsk", "Particion2");
	MOUNTControlador.AgregarMontaje("/home/jmorente/Documentos/Manejo e Implementación de Archivos/DiscoTest.dsk", "Particion3");
	MOUNTControlador.AgregarMontaje("/home/jmorente/Documentos/Manejo e Implementación de Archivos/DiscoTest2.dsk", "Particion1");
	MOUNTControlador.AgregarMontaje("/home/jmorente/Documentos/Manejo e Implementación de Archivos/DiscoTest2.dsk", "Particion2");
	MOUNTControlador.AgregarMontaje("/home/jmorente/Documentos/Manejo e Implementación de Archivos/DiscoTest2.dsk", "Particion1");
	MOUNTControlador.AgregarMontaje("/home/jmorente/Documentos/Manejo e Implementación de Archivos/DiscoTest2.dsk", "Particion2");
	MOUNTControlador.AgregarMontaje("/home/jmorente/Documentos/Manejo e Implementación de Archivos/DiscoTest.dsk", "Particion1");
	MOUNTControlador.AgregarMontaje("/home/jmorente/Documentos/Manejo e Implementación de Archivos/DiscoTest.dsk", "Particion1");
	MOUNTControlador.Mostrar();*/
	//ReporteControlador.REP("mbr", "/home/jmorente/Documentos/MIA/imagen.png", "vda1", "abc");
}