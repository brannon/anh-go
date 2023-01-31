package anh

import (
	"fmt"
	"strings"
)

type WellKnownTagName string

const (
	WellKnownTagNameInstallationId WellKnownTagName = "InstallationId"
	WellKnownTagNameUserId         WellKnownTagName = "UserId"
)

type WellKnownTag struct {
	Name  string
	Value string
}

func (t *WellKnownTag) String() string {
	return fmt.Sprintf("$%s:{%s}", t.Name, t.Value)
}

func ParseWellKnownTag(tag string) (WellKnownTag, bool) {
	if strings.HasPrefix(tag, "$") {
		remain := tag[1:]

		name, remain, found := strings.Cut(remain, ":{")
		if found {
			value, _, found := strings.Cut(remain, "}")
			if found {
				return WellKnownTag{name, value}, true
			}
		}
	}
	return WellKnownTag{}, false
}
