package anh

type atomFeed struct {
	Entries []*atomEntry `xml:"entry"`
}

type atomEntry struct {
	Content *atomEntryContent `xml:"content"`
}

type atomEntryContent struct {
	AppleRegistrationDescription *AppleRegistrationDescription `xml:"AppleRegistrationDescription"`
	GcmRegistrationDescription   *GcmRegistrationDescription   `xml:"GcmRegistrationDescription"`
}
