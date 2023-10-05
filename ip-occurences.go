package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
)

func resolveDomainIP(domain string, ipCounter map[string]int, verbose bool, wg *sync.WaitGroup) {
	defer wg.Done()

	ip, err := net.LookupIP(domain)
	if err == nil {
		for _, addr := range ip {
			ipStr := addr.String()
			ipCounter[ipStr]++
			if verbose {
				fmt.Printf("Resolved %s to IP address: %s\n", domain, ipStr)
			}
		}
	} else if verbose {
		fmt.Printf("Failed to resolve IP address for %s\n", domain)
	}
}

func findCommonIPs(domains []string, numThreads int, verbose bool) map[string]int {
	ipCounter := make(map[string]int)
	var wg sync.WaitGroup

	for _, domain := range domains {
		wg.Add(1)
		go resolveDomainIP(domain, ipCounter, verbose, &wg)
	}

	wg.Wait()

	return ipCounter
}

func main() {
	domainListFile := flag.String("dl", "", "File containing a list of domains")
	numThreads := flag.Int("t", 1, "Number of threads to use for resolving IPs")
	verbose := flag.Bool("v", false, "Enable verbose mode")
	flag.Parse()

	var domains []string

	if *domainListFile != "" {
		file, err := os.Open(*domainListFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			domains = append(domains, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Enter domains (one per line) and press Ctrl+D (Ctrl+Z on Windows) when done:")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			domain := scanner.Text()
			domains = append(domains, domain)
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	ipCounter := findCommonIPs(domains, *numThreads, *verbose)

	if len(ipCounter) > 0 {
		fmt.Println("Most common IP addresses:")
		for ip, count := range ipCounter {
			fmt.Printf("%s: %d occurrences\n", ip, count)
		}
	} else {
		fmt.Println("No IP addresses were resolved.")
	}
}

