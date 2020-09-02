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
	Part_size   		int64
	Part_name   		[16]byte
	Part_next   		int64
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
}

/**
 * MASTER BOOT RECORD
 */
type MBR struct {
	Mbr_size				int64
	Mbr_date 				[19]byte
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
