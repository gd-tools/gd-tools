package setup

import (
	"flag"
	"testing"

	"github.com/gd-tools/gd-tools/platform"
	"github.com/gd-tools/gd-tools/utils"
	"github.com/urfave/cli/v2"
)

func TestBuildConfig(t *testing.T) {
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.Bool("verbose", false, "")
	set.String("swap-size", "0", "")
	set.String("dmarc", "", "")
	set.String("sysadmin", "", "")
	set.String("company", "", "")
	set.String("domain", "", "")
	set.String("help-url", "", "")
	set.String("spambarrier", "", "")
	set.String("ubuntu-pro", "", "")
	set.String("hetzner-dns", "", "")
	set.String("ionos-dns", "", "")
	set.String("cloudflare-dns", "", "")
	set.String("hetzner-volume", "", "")

	err := set.Parse([]string{
		"--swap-size=4G",
		"--hetzner-volume=4711",
		"--company=Example GmbH",
	})
	if err != nil {
		t.Fatal(err)
	}

	c := cli.NewContext(nil, set, nil)

	pf, err := platform.LoadPlatform(platform.DefaultBaseline, nil)
	if err != nil {
		t.Fatal(err)
	}

	id := &utils.Identity{
		Company:  "Example GmbH",
		Domain:   "example.org",
		SysAdmin: "admin@example.org",
		HelpURL:  "https://support.example.org/",
		TimeZone: "Europe/Berlin",
		Language: "de",
		Region:   "DE",
		RegTTL:   3600,
		DMARC:    "v=DMARC1; p=quarantine; pct=100; adkim=s; aspf=s",
	}

	cfg := buildConfig(c, pf, id, "host00", "gd-tools.de")

	if cfg.BaselineName != platform.DefaultBaseline {
		t.Fatalf("unexpected baseline: %q", cfg.BaselineName)
	}

	if cfg.HostName != "host00" {
		t.Fatalf("unexpected host name: %q", cfg.HostName)
	}

	if cfg.DomainName != "gd-tools.de" {
		t.Fatalf("unexpected domain name: %q", cfg.DomainName)
	}

	if cfg.Company != "Example GmbH" {
		t.Fatalf("unexpected company: %q", cfg.Company)
	}

	if cfg.Domain != "example.org" {
		t.Fatalf("unexpected identity domain: %q", cfg.Domain)
	}

	if cfg.SwapSize != "4G" {
		t.Fatalf("unexpected swap size: %q", cfg.SwapSize)
	}

	if cfg.FQDN() != "host00.gd-tools.de" {
		t.Fatalf("unexpected fqdn: %q", cfg.FQDN())
	}

	if len(cfg.UsedFQDNs) != 1 || cfg.UsedFQDNs[0] != "host00.gd-tools.de" {
		t.Fatalf("unexpected used fqdns: %#v", cfg.UsedFQDNs)
	}

	if len(cfg.Mounts) != 1 {
		t.Fatalf("unexpected mount count: %d", len(cfg.Mounts))
	}

	if cfg.Mounts[0].Provider != "hetzner" {
		t.Fatalf("unexpected mount provider: %q", cfg.Mounts[0].Provider)
	}
}
