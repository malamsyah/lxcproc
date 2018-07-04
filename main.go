package main

import ("fmt"
    "os/exec"
    "log")

func main() {
	cmd := exec.Command("lxc", "info", "something-ubuntu")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}
