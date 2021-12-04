package settings

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os"
	"os/exec"
	"strconv"
)

func FindCheckProcess() ([]*process.Process, error) {
	return process.Processes()
}

func a() []string {
	list_, _ := FindCheckProcess()
	agentNums := []string{}
	for _, v := range list_ {
		ss, _ := v.Name()
		if ss == "agentx_linux" {
			agentNums = append(agentNums, ss)
		}
	}
	return agentNums
}

func saveID(pid int) {
	file, err := os.Create(Config().Pid)
	if err != nil {
		fmt.Printf("没有pid file : %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	_, err = file.WriteString(strconv.Itoa(pid))
	if err != nil {
		fmt.Printf("没有pid file : %v\n", err)
		os.Exit(1)
	}
	file.Sync()
}

func StartHandle() {
	if _, err := os.Stat(Config().Pid); err == nil {
		fmt.Println("已经运行或pid文件已存在")
		os.Exit(1)
	}
	agentList := a()
	if len(agentList) > 1 {
		fmt.Println(len(agentList))
		os.Exit(1)
	}
	cmd := exec.Command(os.Args[0], "main")
	cmd.Start()
	fmt.Println(GetVersion())
	fmt.Println("进程已起动 PID is", cmd.Process.Pid)
	fmt.Println(IP())
	saveID((cmd.Process.Pid))
	os.Exit(0)
}
