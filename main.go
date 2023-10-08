package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("Help!")

	}
}

func run() {
	fmt.Printf("Running %vas PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}
	must(cmd.Run())

}
func child() {
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())
	syscall.Sethostname([]byte("Zurich"))

	executableDir := filepath.Dir(".")
	ubuntuDir := filepath.Join(executableDir, "UBUNTU")

	// Now you can use ubuntuDir to chroot into
	log.Printf("Ubuntu Dir is: %s", ubuntuDir)

	must(syscall.Chroot(ubuntuDir))
	must(syscall.Chdir("/"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	must(cmd.Run())

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
