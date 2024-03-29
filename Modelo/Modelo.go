package Modelo

import (
)

/**
 * PARTICION EBR 
 */
 type EBR struct {
	Part_status 		byte
	Part_fit 			[2]byte
	Part_start   		int64
	Part_end	   		int64
	Part_size   		int64
	Part_name   		[16]byte
	Part_next   		int8
	Part_logica			PARTICION_LOGICA
}

/**
 * PARTICION MBR
 */
type PARTICION struct {
	Part_status				byte
	Part_type 				byte
	Part_fit				[2]byte
	Part_start				int64
	Part_end				int64
	Part_size				int64
	Part_size_disponible	int64
	Part_name				[16]byte
	Part_EBR				[25]EBR
}

/**
 * PARTICION LOGICA
 */
 type PARTICION_LOGICA struct {
	Part_status				byte
	Part_type 				byte
	Part_fit				[2]byte
	Part_start				int64
	Part_end				int64
	Part_size				int64
	Part_name				[16]byte
}

/**
 * MASTER BOOT RECORD
 */
type MBR struct {
	Mbr_size				int64
	Mbr_size_disponible		int64
	Mbr_date 				[19]byte
	Mbr_disk_signature   	int64
	Mbr_Extendida	   		int8
	Mbr_partition_1			PARTICION
	Mbr_partition_2			PARTICION
	Mbr_partition_3			PARTICION
	Mbr_partition_4			PARTICION
}

/**
 * MOUNT 
 */
type MOUNT struct {
	Mount_id 				string
	Mount_path 				string
	Mount_particion   		string
	Mount_partition			PARTICION
	Mount_particion_logica 	PARTICION_LOGICA
}

/**
 * SUPERBOOT 
 */
type SUPERBOOT struct {
	SB_nombre_hd									[100]byte
	SB_arbol_virtual_count							int64
	SB_detalle_directorio_count						int64
	SB_inodos_count									int64
	SB_bloques_count								int64
	SB_arbol_virtual_free							int64
	SB_detalle_directorio_free						int64
	SB_inodos_free									int64
	SB_bloques_free									int64
	SB_date_creacion								[19]byte
	SB_date_ultimo_montaje							[19]byte
	SB_montaje_count								int64
	SB_ap_bitmap_arbol_directorio					int64
	SB_ap_arbol_directorio							int64
	SB_ap_bitmap_detalle_directorio					int64
	SB_ap_detalle_directorio						int64
	SB_ap_bitmap_tabla_inodo						int64
	SB_ap_tabla_inodo								int64
	SB_ap_bitmap_bloques							int64
	SB_ap_bloques									int64
	SB_ap_log										int64
	SB_size_struct_arbol_directorio					int64
	SB_size_struct_detalle_directorio				int64
	SB_size_struct_inodo							int64
	SB_size_struct_bloque							int64
	SB_free_first_bit_struct_arbol_directorio		int64
	SB_free_first_bit_struct_detalle_directorio		int64
	SB_free_first_bit_struct_inodo					int64
	SB_free_first_bit_struct_bloque					int64
	SB_magic_num									int64
}

/**
 * ARBOL VIRTUAL DE DIRECTORIO 
 */
type AVD struct {
	Avd_fecha_creacion						[19]byte
	Avd_nombre_directorio					[20]byte
	Avd_ap_array_subdirectorios				[6]int64
	Avd_ap_detalle_directorio				int64
	Avd_ap_arbol_virtual_directorio			int64
	Avd_proper								int64
}

/**
 * DETALLE DE DIRECTORIO
 */
type DD struct {
	DD_array_files 							[5]DDFILE
	DD_ap_detalle_directorio				int64
}

/**
 * DETALLE DE DIRECTORIO
 */
type DDFILE struct {
	DD_file_nombre            				[20]byte
	DD_file_ap_inodo          				int64
	DD_file_date_creacion     				[19]byte
	DD_file_date_modificacion 				[19]byte
}

/**
 * INODO
 */
type INODO struct {
	I_count_inodo             				int64
	I_size_archivo            				int64
	I_count_bloques_asignados 				int64
	I_array_bloques           				[4]int64
	I_ap_indirecto            				int64
	I_id_proper               				int64
}

/**
 * BLOQUE
 */
type BLOQUE struct {
	DB_data 								[25]byte
}

/**
 * BITACORA
 */
type BITACORA struct {
	Log_tipo_operacion 						int64
	log_tipo           						int8
	log_nombre         						[100]byte
	log_contenido      						[50]byte
	log_fecha          						[19]byte
}

/**
 * USUARIO
 */
type USUARIO struct {
	Usuario_id       						string
	Usuario_type     						string
	Usuario_group    						string
	Usuario_username 						string
	Usuario_password 						string
}