package protocol

// PackageList contains Ubuntu package names that should be installed.
//
// The package selection is defined by gdt. Applications do not extend this
// list dynamically through the protocol.
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
