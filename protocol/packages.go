package protocol

type PackageList struct {
	Packages  []string `json:"packages,omitempty"`
	DoUpgrade bool     `json:"do_upgrade"`
	UbuntuPro string   `json:"ubuntu_pro"`
}

func (req *Request) HasPackageList() bool {
	if req == nil {
		return false
	}
	return len(req.Packages) > 0 || req.DoUpgrade || req.UbuntuPro != ""
}
