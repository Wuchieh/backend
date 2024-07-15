package line

import "testing"

func TestGetScope(t *testing.T) {
	t.Log(GetScope(ScopeOption{
		ProfileInformation: true,
		IDToken:            false,
		DisplayName:        false,
		ProfileImageURL:    false,
		EmailAddress:       false,
	}))
}
