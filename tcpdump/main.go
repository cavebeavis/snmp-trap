package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// https://rmoff.net/2019/11/29/using-tcpdump-with-docker/
	// sudo tcpdump --interface any -vvv -A 'port 1162'
	cmd := exec.Command("tcpdump", "--interface", "any", "-vv")

	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stdout

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	for err == nil {
		if strings.Contains(line, ".162") || strings.Contains(line, ".1162") {
			log.Println(line)
		}
		
		line, err = reader.ReadString('\n')
	}

	// scanner := bufio.NewScanner(stdout)
	// scanner.Split(bufio.ScanLines)
	// for scanner.Scan() {
	// 	m := scanner.Text()
	// 	log.Println(m)
	// }

	cmd.Wait()
}