package main

import (
	"strconv"
	"strings"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

// ByteSize returns a human-readable byte string of the form 10M, 12.5K, and so forth.  The following units are available:
//	E: Exabyte
//	P: Petabyte
//	T: Terabyte
//	G: Gigabyte
//	M: Megabyte
//	K: Kilobyte
//	B: Byte
// The unit that results in the smallest number greater than or equal to 1 is always chosen.
func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= EXABYTE:
		unit = "EiB"
		value = value / EXABYTE
	case bytes >= PETABYTE:
		unit = "PiB"
		value = value / PETABYTE
	case bytes >= TERABYTE:
		unit = "TiB"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "GiB"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "MiB"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "KiB"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0B"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}
