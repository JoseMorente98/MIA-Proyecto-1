package Modelo

import (
	"time"
)

/**
 * PARTICION MBR
 */
type PARTICION struct {
	part_status				byte
	part_type 				byte
	part_fit				byte
	part_start				int64
	part_size				int64
	part_name				string
}

/**
 * MASTER BOOT RECORD
 */
type MBR struct {
	mbr_size				int64
	mbr_date 				time.Time
	mbr_disk_signature   	int64
	mbr_partition_1			PARTICION
	mbr_partition_2			PARTICION
	mbr_partition_3			PARTICION
	mbr_partition_4			PARTICION
}

