package cloud

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (d Domain) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Host, validation.Required, is.DNSName),
	)
}

// PrimaryDomain knows how to extract the second and top level domain i.e. tuuturu.org in dev.tuuturu.org
func (d Domain) PrimaryDomain() string {
	parts := strings.Split(d.Host, ".")

	topLevelDomain := parts[len(parts)-1]
	secondLevelDomain := parts[len(parts)-2]

	return fmt.Sprintf("%s.%s", secondLevelDomain, topLevelDomain)
}

// Subdomain knows how to extract the lowest level of domain i.e. dev in dev.tuuturu.org
func (d Domain) Subdomain() string {
	cleaned := strings.TrimSuffix(d.Host, ".")

	parts := strings.Split(cleaned, ".")

	if len(parts) >= requiredPartsForSubdomain {
		return parts[0]
	}

	return ""
}

// FQDN ensures the Host provided satisfies the fully qualified domain name format
func (d Domain) FQDN() string {
	if strings.HasSuffix(d.Host, ".") {
		return d.Host
	}

	return fmt.Sprintf("%s.", d.Host)
}

// String returns the fully qualified domain name excluding the dot
func (d Domain) String() string {
	return strings.TrimRight(d.FQDN(), ".")
}

const requiredPartsForSubdomain = 3
