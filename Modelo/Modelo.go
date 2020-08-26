package Modelo

import (
)

/**
 * PARTICION MBR
 */
type PARTICION struct {
	Particion_status			byte
	Particion_type 				byte
	Particion_fit				byte
	Particion_start				int64
	Particion_size				int64
	Particion_name				string
}

/**
 * MASTER BOOT RECORD
 */
type MBR struct {
	Mbr_size				int64
	Mbr_date 				string
	Mbr_disk_signature   	int64
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
}