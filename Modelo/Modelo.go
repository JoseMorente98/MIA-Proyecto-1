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
	Part_name				[16]byte
	Part_EBR				[100]EBR
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
	SB_free_first_bit_struct_arbol_directorio		byte
	SB_free_first_bit_struct_detalle_directorio		byte
	SB_free_first_bit_struct_inodo					byte
	SB_free_first_bit_struct_bloque					byte
	SB_magic_num									byte
}

/**
 * ARBOL VIRTUAL DE DIRECTORIO 
 */
type AVD struct {
	Avd_fecha_creacion								[19]byte
	Avd_nombre_directorio							[100]byte
	Avd_ap_array_subdirectorios						byte
	Avd_ap_detalle_directorio						byte
	Avd_ap_arbol_virtual_directorio					byte
	Avd_proper										byte
}

type USUARIO struct {
	Usuario_id       		string
	Usuario_type     		string
	Usuario_group    		string
	Usuario_username 		string
	Usuario_password 		string
}