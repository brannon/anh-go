package anh

import "strings"

type Platform string

const (
	PlatformApple Platform = "apns"
	PlatformGcm   Platform = "gcm"
)

type Registration interface {
	GetPlatform() Platform
	GetRegistrationId() string

	AddTag(tag string)
	GetTags() []string
	RemoveTag(tag string)
	SetTags(tags []string)
}

//
// RegistrationDescription
//

type RegistrationDescription struct {
	ETag           string `xml:"ETag"`
	ExpirationTime string `xml:"ExpirationTime"`
	RegistrationId string `xml:"RegistrationId"`
	Tags           string `xml:"Tags"`
}

func (desc *RegistrationDescription) GetRegistrationId() string {
	return desc.RegistrationId
}

func (desc *RegistrationDescription) AddTag(tag string) {
	tags := desc.parseTags()
	tags = append(tags, tag)
	desc.Tags = desc.joinTags(tags)
}

func (desc *RegistrationDescription) GetTags() []string {
	return desc.parseTags()
}

func (desc *RegistrationDescription) joinTags(tags []string) string {
	if len(tags) == 0 || (len(tags) == 1 && tags[0] == "") {
		return ""
	}
	return strings.Join(tags, ",")
}

func (desc *RegistrationDescription) parseTags() []string {
	t := strings.Split(desc.Tags, ",")
	if len(t) == 1 && t[0] == "" {
		return []string{}
	}
	return t
}

func (desc *RegistrationDescription) RemoveTag(tag string) {
	tags := desc.parseTags()
	for i, t := range tags {
		if t == tag {
			tags = append(tags[:i], tags[i+1:]...)
			break
		}
	}
	desc.Tags = desc.joinTags(tags)
}

func (desc *RegistrationDescription) SetTags(tags []string) {
	desc.Tags = desc.joinTags(tags)
}

//
// AppleRegistrationDescription
//

type AppleRegistrationDescription struct {
	RegistrationDescription

	DeviceToken string `xml:"DeviceToken"`
}

func (desc *AppleRegistrationDescription) GetPlatform() Platform {
	return PlatformApple
}

//
// GcmRegistrationDescription
//

type GcmRegistrationDescription struct {
	RegistrationDescription

	GcmRegistrationId string `xml:"GcmRegistrationId"`
}

func (desc *GcmRegistrationDescription) GetPlatform() Platform {
	return PlatformGcm
}
