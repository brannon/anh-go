package anh

import "strings"

type installationFeed struct {
	Entries []installationEntry `xml:"entry"`
}

type installationEntry struct {
	Content installationContent `xml:"content"`
}

type installationContent struct {
	AppleRegistrationDescription *installationAppleRegistration `xml:"AppleRegistrationDescription"`
	GcmRegistrationDescription   *installationGcmRegistration   `xml:"GcmRegistrationDescription"`
}

type installationRegistration struct {
	Tags string `xml:"Tags"`
}

type installationAppleRegistration struct {
	installationRegistration
}

type installationGcmRegistration struct {
	installationRegistration
}

func (r *installationRegistration) GetInstallationId() (string, bool) {
	_, remain, found := strings.Cut(r.Tags, "$InstallationId:{")
	if found {
		id, _, found := strings.Cut(remain, "}")
		if found {
			return id, true
		}
	}
	return "", false
}

func (r *installationRegistration) IsInstallation() bool {
	_, found := r.GetInstallationId()
	return found
}
