package anh

type atomFeed struct {
	Entries []*atomEntry `xml:"entry"`
}

type atomEntry struct {
	Content *atomEntryContent `xml:"content"`
}

func (e *atomEntry) GetRegistration() Registration {
	if e.Content != nil {
		if e.Content.AppleRegistrationDescription != nil {
			return e.Content.AppleRegistrationDescription
		} else if e.Content.GcmRegistrationDescription != nil {
			return e.Content.GcmRegistrationDescription
		}
	}

	return nil
}

type atomEntryContent struct {
	AppleRegistrationDescription *AppleRegistrationDescription `xml:"AppleRegistrationDescription"`
	GcmRegistrationDescription   *GcmRegistrationDescription   `xml:"GcmRegistrationDescription"`
}
