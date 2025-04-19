package blocklist

import (
	"bufio"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"
)

var BlocklistReady = false
var BlocklistMutex sync.RWMutex
var Blocklist []string

var BlocklistRangesMutex sync.RWMutex
var BlocklistRanges []net.IPNet

// Fetches the latest blocklist from the remote servers and updates the local blocklist
func RefreshBlocklists() {
	BlocklistReady = false
	var wg sync.WaitGroup

	// Singular IPs
	BlocklistMutex.Lock()
	Blocklist = nil
	BlocklistMutex.Unlock()

	go appendLatestIpbans(&wg)

	// Ranges
	BlocklistRangesMutex.Lock()
	BlocklistRanges = nil
	BlocklistRangesMutex.Unlock()

	go appendLatestDatacenterIps(&wg)
	go appendLatestVpnIps(&wg)

	wg.Wait()
	BlocklistReady = true
}

func appendLatestIpbans(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	slog.Info("Fetching ip block list")
	res, err := http.Get("https://raw.githubusercontent.com/bitwire-it/ipblocklist/refs/heads/main/ip-list.txt")

	if err != nil {
		slog.Error("Got an error while fetching blocklist", "error", err)
		return
	}

	slog.Info("Processing blocklist IPs")

	ipcount := 0
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		// Try parsing ip, if fail just warn & ignore
		if net.ParseIP(line) == nil {
			slog.Warn("Unable to parse IP!", "ip", line)
			continue
		}

		ipcount++
		// God forbid im locking & unlocking too fast :3
		BlocklistMutex.Lock()
		Blocklist = append(Blocklist, line)
		BlocklistMutex.Unlock()
	}

	slog.Info("Done processing IPs", "ipCount", ipcount)
}

func appendLatestDatacenterIps(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	slog.Info("Fetching datacenter ip list")
	res, err := http.Get("https://github.com/X4BNet/lists_vpn/raw/refs/heads/main/output/datacenter/ipv4.txt")
	if err != nil {
		slog.Error("Got an error while fetching datacenter ip list", "error", err)
		return
	}

	slog.Info("Processing datacenter IPs")
	ipcount := 0
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		ipcount++
		_, ipnet, err := net.ParseCIDR(line)
		if err != nil {
			slog.Warn("Unable to parse IPNet!", "ip", line, "err", err)
			continue
		}

		BlocklistRangesMutex.Lock()
		BlocklistRanges = append(BlocklistRanges, *ipnet)
		BlocklistRangesMutex.Unlock()
	}
	slog.Info("Done processing datacenter IPs", "ipCount", ipcount)
}

func appendLatestVpnIps(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	slog.Info("Fetching VPNs ip list")
	res, err := http.Get("https://raw.githubusercontent.com/X4BNet/lists_vpn/refs/heads/main/output/vpn/ipv4.txt")
	if err != nil {
		slog.Error("Got an error while fetching VPNs ip list", "error", err)
		return
	}

	slog.Info("Processing VPNs IPs")
	ipcount := 0
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		ipcount++
		_, ipnet, err := net.ParseCIDR(line)
		if err != nil {
			slog.Warn("Unable to parse IPNet!", "ip", line, "err", err)
			continue
		}

		BlocklistRangesMutex.Lock()
		BlocklistRanges = append(BlocklistRanges, *ipnet)
		BlocklistRangesMutex.Unlock()
	}
	slog.Info("Done processing VPNs IPs", "ipCount", ipcount)
}

// Is provided IP in the blocklist
func Contains(ip net.IP) bool {
	BlocklistMutex.RLock()
	isBlocked := slices.Contains(Blocklist, ip.String())
	BlocklistMutex.RUnlock()

	if isBlocked {
		return true
	}

	BlocklistRangesMutex.RLock()
	defer BlocklistRangesMutex.RUnlock()
	for _, ipnet := range BlocklistRanges {
		if ipnet.Contains(ip) {
			isBlocked = true
			break
		}
	}
	return isBlocked
}

// Schedules the blocklist refresh every 6 hours
func ScheduleRefreshBlocklists(quitSignal chan struct{}) {
	ticker := time.NewTicker(6 * time.Hour)
	slog.Info("Refreshing IP blocklist every 6 hours")

	for {
		select {
		case <-ticker.C:
			slog.Info("Refreshing blocklists")
			RefreshBlocklists()
		case <-quitSignal:
			slog.Info("Received quit signal, stopping blocklist refresh")
			ticker.Stop()
		}
	}
}
