package osstat

import (
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"time"
)

func GetOSStat() (OSStatReport, error) {
	mem, err := memory.Get()
	if err != nil {
		return OSStatReport{}, err
	}

	beforeCPU, err := cpu.Get()
	if err != nil {
		return OSStatReport{}, err
	}
	//beforeNet, err := network.Get()
	//if err != nil {
	//	return OSStatReport{}, err
	//}
	//beforeDisk, err := disk.Get()
	//if err != nil {
	//	return OSStatReport{}, err
	//}

	time.Sleep(time.Second)

	afterCPU, err := cpu.Get()
	if err != nil {
		return OSStatReport{}, err
	}
	//afterNet, err := network.Get()
	//if err != nil {
	//	return OSStatReport{}, err
	//}
	//afterDisk, err := disk.Get()
	//if err != nil {
	//	return OSStatReport{}, err
	//}

	totalCPU := float64(afterCPU.Total - beforeCPU.Total)

	//var (
	//	totalNetIn, totalNetOut, totalDiskIn, totalDiskOut uint64
	//)
	//
	//for _, v := range afterNet {
	//	totalNetIn += v.RxBytes
	//	totalNetOut += v.TxBytes
	//}
	//for _, v := range beforeNet {
	//	totalNetIn -= v.RxBytes
	//	totalNetOut -= v.TxBytes
	//}
	//for _, v := range afterDisk {
	//	totalDiskIn += v.WritesCompleted
	//	totalDiskOut += v.ReadsCompleted
	//}
	//for _, v := range beforeDisk {
	//	totalDiskIn -= v.WritesCompleted
	//	totalDiskOut -= v.ReadsCompleted
	//}

	return OSStatReport{
		MemoryUsage: (1 - float64(mem.Free+mem.Buffers+mem.Cached)/float64(mem.Total)) * 100,
		CpuUsage:    float64(afterCPU.User-beforeCPU.User+afterCPU.System-beforeCPU.System) / totalCPU * 100,
		//DiskIn:      totalDiskIn,
		//DiskOut:     totalDiskOut,
		//NetIn:       totalNetIn,
		//NetOut:      totalNetOut,
	}, nil
}
