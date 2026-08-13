package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spa-nsk/oneview2"
)

// Выгрузка всех серверов OneView в каталог JSON-файлов.
//
//	set ONEVIEW_OV_ENDPOINT=https://oneview.example.com
//	set ONEVIEW_OV_USER=Administrator
//	set ONEVIEW_OV_PASSWORD=secret
//	go run . ./out
func main() {
	dir := "server-exports"
	if len(os.Args) > 1 {
		dir = os.Args[1]
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

	exports, err := c.ExportServers(ctx, oneview.ListOptions{Count: -1, Sort: "name:asc"}, oneview.ExportOptions{})
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatal(err)
	}
	for _, exp := range exports {
		fmt.Print(exp.Summary())
		name := sanitize(exp.Identity.Name)
		if name == "" {
			name = sanitize(exp.Identity.SerialNumber)
		}
		path := filepath.Join(dir, name+".json")
		if err := oneview.SaveServerExportJSON(path, exp); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  -> %s\n", path)
	}
	fmt.Printf("exported %d servers to %s\n", len(exports), dir)
}

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	repl := strings.NewReplacer(`/`, `_`, `\`, `_`, `:`, `_`, `*`, `_`, `?`, `_`, `"`, `_`, `<`, `_`, `>`, `_`, `|`, `_`, `,`, `_`)
	s = repl.Replace(s)
	s = strings.Join(strings.Fields(s), "_")
	if s == "" {
		return "server"
	}
	return s
}
