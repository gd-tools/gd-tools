package agent

const (
	PrefixAdmidio      = "ad"
	PrefixBookStack    = "bk"
	PrefixDovecot      = "dc" // system module
	PrefixFirefly      = "fi"
	PrefixImmich       = "im"
	PrefixMediaWiki    = "mw"
	PrefixMinecraft    = "mc"
	PrefixNextcloud    = "nc"
	PrefixOCIS         = "oc"
	PrefixOpenDKIM     = "od" // system module
	PrefixPaperlessNGX = "pl"
	PrefixPostfix      = "pf" // system module
	PrefixRoundcube    = "rc"
	PrefixRustDesk     = "rd"
	PrefixUptimeKuma   = "uk"
	PrefixVaultwarden  = "vw"
	PrefixWordPress    = "wp"
)

const (
	NamingAdmidio      = "admidio"
	NamingBookStack    = "bookstack"
	NamingFirefly      = "firefly"
	NamingImmich       = "immich"
	NamingMediaWiki    = "mediawiki"
	NamingMinecraft    = "minecraft"
	NamingNextcloud    = "nextcloud"
	NamingOCIS         = "ocis"
	NamingPaperlessNGX = "paperless"
	NamingRoundcube    = "roundcube"
	NamingRustDesk     = "rustdesk"
	NamingUptimeKuma   = "uptimekuma"
	NamingVaultwarden  = "vaultwarden"
	NamingWordPress    = "wordpress"
)

type NamingScheme struct {
	Short string
	Name  string
}

var namingSchemes = map[string]*NamingScheme{
	NamingAdmidio:      {Short: PrefixAdmidio, Name: NamingAdmidio},
	NamingBookStack:    {Short: PrefixBookStack, Name: NamingBookStack},
	NamingFirefly:      {Short: PrefixFirefly, Name: NamingFirefly},
	NamingImmich:       {Short: PrefixImmich, Name: NamingImmich},
	NamingMediaWiki:    {Short: PrefixMediaWiki, Name: NamingMediaWiki},
	NamingMinecraft:    {Short: PrefixMinecraft, Name: NamingMinecraft},
	NamingNextcloud:    {Short: PrefixNextcloud, Name: NamingNextcloud},
	NamingOCIS:         {Short: PrefixOCIS, Name: NamingOCIS},
	NamingPaperlessNGX: {Short: PrefixPaperlessNGX, Name: NamingPaperlessNGX},
	NamingRoundcube:    {Short: PrefixRoundcube, Name: NamingRoundcube},
	NamingRustDesk:     {Short: PrefixRustDesk, Name: NamingRustDesk},
	NamingUptimeKuma:   {Short: PrefixUptimeKuma, Name: NamingUptimeKuma},
	NamingVaultwarden:  {Short: PrefixVaultwarden, Name: NamingVaultwarden},
	NamingWordPress:    {Short: PrefixWordPress, Name: NamingWordPress},
}

// GetNamingScheme returns the configured naming scheme.
func GetNamingScheme(name string) *NamingScheme {
	ns, ok := namingSchemes[name]
	if !ok {
		panic("unknown naming scheme: " + name)
	}
	return ns
}
