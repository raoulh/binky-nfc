package mdns

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

var (
	hostConn = ""
)

//DiscoverBinkyServer search for Binky server with mDNS
func DiscoverBinkyServer() (string, error) {
	log.Println("Discover Binky Server...")
	resolver, err := zeroconf.NewResolver(zeroconf.SelectIPTraffic(zeroconf.IPv4))
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
		time.Sleep(time.Second * 5)
		return "", err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			//log.Println(entry)
			hostConn = fmt.Sprintf("%s:%d", entry.AddrIPv4, entry.Port)
		}
		//log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	err = resolver.Browse(ctx, "_binky._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
		time.Sleep(time.Second * 2)
		return "", err
	}

	<-ctx.Done()

	if hostConn == "" {
		return "", fmt.Errorf("no binky server found")
	}

	return hostConn, nil
}
