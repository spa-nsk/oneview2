package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spa-nsk/oneview2"
)

func main() {
	cfg := oneview.ConfigFromEnv()
	if cfg.Endpoint == "" {
		cfg = oneview.Config{
			Endpoint:    os.Getenv("OV_ENDPOINT"),
			Username:    os.Getenv("OV_USER"),
			Password:    os.Getenv("OV_PASSWORD"),
			Domain:      os.Getenv("OV_DOMAIN"),
			InsecureTLS: true,
			Timeout:     60 * time.Second,
		}
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

	min, cur := c.ApplianceVersion()
	fmt.Printf("product=%s api=%d (appliance range %d-%d)\n", c.Product(), c.APIVersion(), min, cur)

	hw, err := c.ListServerHardware(ctx, oneview.ListOptions{Count: -1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("server-hardware: %d / %d\n", hw.Count, hw.Total)
	for _, s := range hw.Members {
		fmt.Printf("  %s  power=%s  status=%s  model=%s\n", s.Name, s.PowerState, s.Status, s.ShortModel)
	}
}
