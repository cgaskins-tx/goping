// goping2.go
// author: Chris C. Gaskins
// version 1.0
//
// This program pings 255 addresses: 192.168.1.1 through 192.168.1.255 using Go Routines
//

package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// function that checks if IP address is online or not. Returns "Online" or "Offline"
func ipAddressStatus(address string, ch chan string) {

	// Use the ping command, with the -c 4 option to send 4 packets
	cmd := exec.Command("ping", "-c", "4", address)

	// Run the command and capture the output
	_, err := cmd.CombinedOutput()

	// write to channel based upon result of error check
	if err == nil {
		ch <- fmt.Sprintf("%s,%s", address, "online")
	} else {
		ch <- fmt.Sprintf("%s,%s", address, "offline")
	}

}

func main() {

	var addr string = ""
	var pingResultList [256]string

	// make the channel
	addressChannel := make(chan string)

	// print some init stuff to the screen
	fmt.Println("goping2 v1")
	fmt.Println("Pinging 255 addresses: 192.168.1.1 through 192.168.1.255 \n")

	// kick off the go routines to ping each of the 255 ip addresses
	fmt.Println("Status: Kicking off Go Routines to check IP address status.")
	for a := 1; a < 256; a++ {
		addr = fmt.Sprintf("%s%d", "192.168.1.", a)
		go ipAddressStatus(addr, addressChannel)
	}

	// status update ot the screen
	fmt.Print("Status: Waiting for results and parsing the data")
	pingResult := ""

	// read from the channel 255 times to get all the responses
	for i := 1; i < 256; i++ {

		// print a dot as a progress indicator every increment of 5
		if i%5 == 0 {
			fmt.Print(".")
		}

		// read from the channel
		pingResult = <-addressChannel

		// split the string read from the channel in the ss[0]=ip address and ss[1]=status of ping of ip address
		ss := strings.FieldsFunc(pingResult, func(r rune) bool {
			if r == ',' {
				return true
			}
			return false
		})
		ipAddress := ss[0] // ip address
		ipStatus := ss[1]  // status of ping of ip address

		// get the 4th field from the IPV4 ip address
		endingip := strings.FieldsFunc(ipAddress, func(r rune) bool {
			if r == '.' {
				return true
			}
			return false
		})
		// convert the 4th string field of the ip address, convert it to integer store in variable n to use as an index into the array
		n, err := strconv.Atoi(endingip[3])
		if err != nil {
			fmt.Println(err)
		} else {
			// store the status result read from the channel for the array indexed to n in the array
			pingResultList[n] = ipStatus
		}
	}

	// initialize a couple of variables
	onlineIPs := 0
	totalIPs := 0
	var xvar string
	// print to the screen
	fmt.Println("\n\nAddresses that are online:")
	// iterate through the array and print out only IP addresses that are online
	for j := 1; j < 256; j++ {

		totalIPs++
		if pingResultList[j] == "online" {
			// format the string to print with a properly formatted IP address and the ping result from the array
			xvar = fmt.Sprintf("192.168.1.%-3d : %s", j, pingResultList[j])
			// print to the screen
			fmt.Println(xvar)
			onlineIPs++
		}
	}
	// print some statistics to the screen
	fmt.Println("Online IP Addresses: ", onlineIPs)
	fmt.Println("Offline IP Addresses:", totalIPs-onlineIPs)
	fmt.Println("Done.")

}
