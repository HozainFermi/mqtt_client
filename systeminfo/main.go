package systeminfo

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetVideocard() string {
	bytes, err := exec.Command("nvidia-smi").Output()
	if err != nil {
		return err.Error()
	}
	return string(bytes[:]) + "\n"
}

func GetPercentEvery() string {
	var ret string = "Usage per CPU:" + "\n"
	out, err := cpu.Percent(time.Second, true)
	if err != nil {
		return ret
	}
	for i := 0; i < len(out); i++ {
		ret += fmt.Sprintf("Usage[%d]=", i) + strconv.FormatFloat(out[i], 'f', 2, 64) + "\n"
	}
	return ret
}

func GetPercent() string {

	out, err := cpu.Percent(time.Second, false)
	if err != nil {
		return err.Error()
	}
	return "CPU usage (combined)=" + strconv.FormatFloat(out[0], 'f', 2, 64) + "%" + "\n"
}

func GetMemUsage() string {
	memUsage, err := mem.VirtualMemory()
	if err != nil {
		return err.Error()
	}
	ret := memUsage.UsedPercent

	return "Memory usage=" + strconv.FormatFloat(ret, 'g', -1, 64) + "%" + "\n"
}
