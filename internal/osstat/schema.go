package osstat

type OSStatReport struct {
	MemoryUsage, CpuUsage float64
	//DiskIn, DiskOut, NetIn, NetOut uint64
}
