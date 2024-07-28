package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

type Device struct {
	IPv4         string
	IPv6         string
	HardwareAddr string
	DefaultRoute string
}

func main() {
	localIP, err := getLocalIP()
	if err != nil {
		log.Fatalf("Error getting local IP: %v", err)
	}
	fmt.Printf("Local IP: %s\n", localIP)

	devices, err := scanNetwork(localIP)
	if err != nil {
		log.Fatalf("Error scanning network: %v", err)
	}

	displayResults(devices)
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no local IP address found")
}

func scanNetwork(localIP string) ([]Device, error) {
	parts := strings.Split(localIP, ".")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid IP address format")
	}
	networkPrefix := strings.Join(parts[:3], ".")

	fmt.Println("Scanning network with nmap...")
	cmd := exec.Command("nmap", "-sn", networkPrefix+".0/24")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running nmap: %v", err)
	}

	return parseNmapOutput(string(output))
}

func parseNmapOutput(output string) ([]Device, error) {
	devices := []Device{}
	lines := strings.Split(output, "\n")

	ipRegex := regexp.MustCompile(`Nmap scan report for (\d+\.\d+\.\d+\.\d+)`)

	for _, line := range lines {
		if ipMatch := ipRegex.FindStringSubmatch(line); len(ipMatch) == 2 {
			device, err := getDeviceDetails(ipMatch[1])
			if err != nil {
				log.Printf("Error getting device details for %s: %v", ipMatch[1], err)
				continue
			}
			devices = append(devices, *device)
		}
	}

	return devices, nil
}

func getDeviceDetails(ip string) (*Device, error) {
	addrs, err := net.LookupIP(ip)
	if err != nil {
		return nil, err
	}

	device := &Device{
		IPv4: ip,
	}

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			device.IPv4 = ipv4.String()
		} else if ipv6 := addr.To16(); ipv6 != nil {
			device.IPv6 = ipv6.String()
		}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.Equal(net.ParseIP(ip)) {
					device.HardwareAddr = iface.HardwareAddr.String()
					break
				}
			}
		}
	}

	routes, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		if ipnet, ok := route.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				device.DefaultRoute = ipnet.IP.String()
				break
			}
		}
	}

	return device, nil
}

func displayResults(devices []Device) {
	if len(devices) == 0 {
		fmt.Println("No devices found.")
		return
	}

	fmt.Println("Connected devices:")
	for i, device := range devices {
		fmt.Printf("%d. IPv4: %s, IPv6: %s, Hardware Address: %s, Default Route: %s\n",
			i+1, device.IPv4, device.IPv6, device.HardwareAddr, device.DefaultRoute)
	}
}
