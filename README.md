# Network Device Scanner

This Go application scans the local network and retrieves detailed information about connected devices. It uses the `nmap` tool to scan the network and the `net` package to gather additional details such as IPv4, IPv6, hardware address, and default route.

## Features

- Scans the local network for connected devices.
- Retrieves detailed information about each device, including:
  - IPv4 address
  - IPv6 address
  - Hardware address (MAC address)
  - Default route

## Prerequisites

- Go (Golang) installed on your system.
- `nmap` installed and available in your system's PATH.

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/network-device-scanner.git
    cd network-device-scanner
    ```

2. Install the required dependencies:

    ```sh
    go mod tidy
    ```

## Usage

1. Build the application:

    ```sh
    go build -o network-device-scanner
    ```

2. Run the application:

    ```sh
    ./network-device-scanner
    ```

## Code Overview

### Main Function

The `main` function initializes the scanning process by retrieving the local IP address and then scanning the network.

```go
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
