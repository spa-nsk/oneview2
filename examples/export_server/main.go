package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spa-nsk/oneview2"
)

// Выгрузка детальной конфигурации одного сервера (CPU, RAM, диски, PCI, NIC, firmware).
//
//	set ONEVIEW_OV_ENDPOINT=https://oneview.example.com
//	set ONEVIEW_OV_USER=Administrator
//	set ONEVIEW_OV_PASSWORD=secret
//	go run . "Encl1, bay 5"
//	go run . /rest/server-hardware/<uuid>
func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: export_server <name|uri|uuid|serial>")
	}
	cfg := oneview.ConfigFromEnv()
	if cfg.Endpoint == "" {
		cfg.Endpoint = os.Getenv("OV_ENDPOINT")
		cfg.Username = os.Getenv("OV_USER")
		cfg.Password = os.Getenv("OV_PASSWORD")
		cfg.Domain = os.Getenv("OV_DOMAIN")
		cfg.InsecureTLS = true
		cfg.Timeout = 90 * time.Second
	}
	if cfg.Endpoint == "" || cfg.Username == "" {
		log.Fatal("set ONEVIEW_OV_ENDPOINT / ONEVIEW_OV_USER / ONEVIEW_OV_PASSWORD")
	}

	c, err := oneview.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if err := c.Login(ctx); err != nil {
		log.Fatal(err)
	}
	defer c.Logout(ctx)

	exp, err := c.ExportServer(ctx, os.Args[1], oneview.ExportOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(exp.Summary())
	fmt.Println("--- processors ---")
	for _, p := range exp.Processors.Sockets {
		fmt.Printf("  socket=%s  %s  cores=%d  threads=%d  %d MHz  %s\n",
			p.Socket, p.Model, p.TotalCores, p.TotalThreads, p.MaxSpeedMHz, p.Health)
	}
	fmt.Println("--- memory DIMMs ---")
	for _, m := range exp.Memory.Modules {
		fmt.Printf("  %-16s  %d MiB  %s  %s  %s  %s\n",
			m.DeviceLocator, m.CapacityMiB, m.MemoryDeviceType, m.PartNumber, m.SerialNumber, m.DIMMStatus)
	}
	fmt.Println("--- physical disks ---")
	for _, d := range exp.Storage.Drives {
		fmt.Printf("  loc=%s  %s  %s  %d MiB  %s/%s  fw=%s  %s\n",
			d.Location, d.MediaType, d.InterfaceType, d.CapacityMiB, d.Model, d.SerialNumber, d.FirmwareVersion, d.Health)
	}
	fmt.Println("--- logical volumes ---")
	for _, v := range exp.Storage.Volumes {
		fmt.Printf("  #%d %s  RAID=%s  %d MiB  %s\n", v.Number, v.Name, v.RAID, v.CapacityMiB, v.Health)
	}
	fmt.Println("--- PCI / devices ---")
	for _, d := range exp.Devices {
		fmt.Printf("  %-22s  %-28s  loc=%s  fw=%s\n", d.DeviceType, d.Name, d.Location, d.FirmwareVersion)
	}
	fmt.Println("--- NIC ports ---")
	for _, p := range exp.NetworkPorts {
		fmt.Printf("  %s slot=%s port=%d  %s  mac=%s\n", p.DeviceName, p.DeviceSlot, p.PortNumber, p.Type, p.MAC)
	}

	out := "server-export.json"
	if err := oneview.SaveServerExportJSON(out, exp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nJSON written to %s\n", out)
	for _, w := range exp.Warnings {
		fmt.Printf("warning: %s\n", w)
	}
}
