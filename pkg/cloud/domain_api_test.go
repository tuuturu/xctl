package cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubdomain_PrimaryDomain(t *testing.T) {
	testCases := []struct {
		name               string
		withURL            string
		expectParentDomain string
	}{
		{
			name:               "Should work with one subdomain",
			withURL:            "test.tuuturu.org",
			expectParentDomain: "tuuturu.org",
		},
		{
			name:               "Should work with multiple subdomains",
			withURL:            "a.b.test.tuuturu.org",
			expectParentDomain: "tuuturu.org",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			d := Domain{Host: tc.withURL}

			assert.Equal(t, tc.expectParentDomain, d.PrimaryDomain())
		})
	}
}

func TestSubdomain_Subdomain(t *testing.T) {
	testCases := []struct {
		name            string
		withFullDomain  string
		expectSubdomain string
	}{
		{
			name:            "Should work",
			withFullDomain:  "test.tuuturu.org",
			expectSubdomain: "test",
		},
		{
			name:            "Should return empty string with no subdomain",
			withFullDomain:  "klokkinn.no",
			expectSubdomain: "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			d := Domain{Host: tc.withFullDomain}

			assert.Equal(t, tc.expectSubdomain, d.Subdomain())
		})
	}
}

func TestDomain_FQDN(t *testing.T) {
	testCases := []struct {
		name       string
		withDomain Domain
		expectFQDN string
	}{
		{
			name:       "Should add a punctuation to a primary domain",
			withDomain: Domain{Host: "tuuturu.org"},
			expectFQDN: "tuuturu.org.",
		},
		{
			name:       "Should add a punctuation to a domain with one subdomain",
			withDomain: Domain{Host: "cluster.tuuturu.org"},
			expectFQDN: "cluster.tuuturu.org.",
		},
		{
			name:       "Should return the original Host when Host contains a punctuation",
			withDomain: Domain{Host: "cluster.tuuturu.org."},
			expectFQDN: "cluster.tuuturu.org.",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.withDomain.FQDN(), tc.expectFQDN)
		})
	}
}

func TestDomain_String(t *testing.T) {
	testCases := []struct {
		name         string
		withDomain   Domain
		expectString string
	}{
		{
			name:         "Should return original Host when Host is a primary domain without punctuation",
			withDomain:   Domain{Host: "tuuturu.org"},
			expectString: "tuuturu.org",
		},
		{
			name:         "Should return original Host when Host has a subdomain without punctuation",
			withDomain:   Domain{Host: "cluster.tuuturu.org"},
			expectString: "cluster.tuuturu.org",
		},
		{
			name:         "Should return Host without punctuation when Host is a primary domain FQDN",
			withDomain:   Domain{Host: "tuuturu.org."},
			expectString: "tuuturu.org",
		},
		{
			name:         "Should return Host without punctuation when Host is a FQDN with a subdomain",
			withDomain:   Domain{Host: "cluster.tuuturu.org."},
			expectString: "cluster.tuuturu.org",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.withDomain.String(), tc.expectString)
		})
	}
}

func TestDomain_Validate(t *testing.T) {
	testCases := []struct {
		name       string
		withDomain Domain
		expectErr  string
	}{
		{
			name:       "Should accept a primary domain",
			withDomain: Domain{Host: "tuuturu.org"},
		},
		{
			name:       "Should accept a domain with subdomain",
			withDomain: Domain{Host: "cluster.tuuturu.org"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.withDomain.Validate()

			if tc.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.expectErr, err.Error())
			}
		})
	}
}
