package anh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InstallationRegistration_GetInstallationId(t *testing.T) {
	registration := installationRegistration{
		Tags: "$InstallationId:{123}",
	}
	id, found := registration.GetInstallationId()
	assert.True(t, found)
	assert.Equal(t, "123", id)

	registration = installationRegistration{
		Tags: "$InstallationId:{123",
	}
	id, found = registration.GetInstallationId()
	assert.False(t, found)
	assert.Equal(t, "", id)

	registration = installationRegistration{
		Tags: "InstallationId:{123}",
	}
	id, found = registration.GetInstallationId()
	assert.False(t, found)
	assert.Equal(t, "", id)
}

func Test_InstallationRegistration_IsInstallation(t *testing.T) {
	registration := installationRegistration{
		Tags: "$InstallationId:{123}",
	}
	assert.True(t, registration.IsInstallation())

	registration = installationRegistration{
		Tags: "$InstallationId:{123",
	}
	assert.False(t, registration.IsInstallation())

	registration = installationRegistration{
		Tags: "InstallationId:{123}",
	}
	assert.False(t, registration.IsInstallation())
}
