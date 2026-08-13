// Package oneview is a Go client for HPE OneView REST APIs.
//
// It covers two surfaces:
//
//   - HPE OneView appliance REST API versions 3800–8800 (OneView 6.60 and later).
//   - HPE OneView Global Dashboard REST API as described in swagger 300.json
//     (X-API-Version 2–300).
//
// Typical usage against an appliance:
//
//	c, err := oneview.New(oneview.Config{
//		Endpoint:    "https://oneview.example.com",
//		Username:    "Administrator",
//		Password:    "secret",
//		Domain:      "LOCAL",
//		InsecureTLS: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer c.Logout(context.Background())
//
//	if err := c.Login(context.Background()); err != nil {
//		log.Fatal(err)
//	}
//	hw, err := c.ListServerHardware(ctx, oneview.ListOptions{Count: -1})
//
// Detailed server dump (CPU, DIMMs, disks, PCI, firmware):
//
//	exp, err := c.ExportServer(ctx, "Encl1, bay 5", oneview.ExportOptions{})
//	fmt.Print(exp.Summary())
//
// Collect every server from several appliances / API versions (duplicates merged):
//
//	servers, err := oneview.CollectServers(ctx, []oneview.Config{cfgDashboard, cfgAppliance})
package oneview
