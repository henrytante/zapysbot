package clear

import (
	"os"
	"os/exec"
	"runtime"
)


func Clear()  {
	var system string
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin"{
		system = "clear"
	}else{
		system = "cls"
	}
	cmd := exec.Command(system)
	cmd.Stdout = os.Stdout
	cmd.Run()
}