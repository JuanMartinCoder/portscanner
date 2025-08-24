package utils

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func Ulimit() int64 {
	cmd := exec.Command("bash", "-c", "ulimit -n")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	s := strings.TrimSpace(string(out))
	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		panic(err)
	}
	return i
}
