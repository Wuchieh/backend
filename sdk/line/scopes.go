package line

const (
	ScopeProfile = "profile"
	ScopeEmail   = "email"
	ScopeOpenid  = "openid"
)

type ScopeOption struct {
	ProfileInformation bool

	IDToken bool

	DisplayName bool

	ProfileImageURL bool

	EmailAddress bool
}

func GetScope(option ScopeOption) []string {
	scope := make(map[string]bool)

	if option.ProfileInformation || option.DisplayName || option.ProfileImageURL {
		scope[ScopeProfile] = true
	}

	if option.IDToken || option.EmailAddress {
		scope[ScopeEmail] = true
	}

	if option.DisplayName || option.ProfileImageURL {
		scope[ScopeOpenid] = true
	}

	// 如果只有 Email 範圍，添加 Profile 範圍
	if len(scope) == 1 && scope[ScopeEmail] {
		scope[ScopeProfile] = true
	}

	result := make([]string, 0, len(scope))
	for k := range scope {
		result = append(result, k)
	}

	return result
}
