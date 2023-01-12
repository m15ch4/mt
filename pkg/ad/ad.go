package ad

type AD struct {
	Server string
	Port   int32
	BaseDN string
}

func (a *AD) GetUserGroups(username string) []string {
	groups := make([]string, 3)
	groups[0] = "T01"
	groups[1] = "N01"
	groups[2] = "U01"
	return groups
}

func NewAD(adserver string, basedn string, port int32) *AD {
	return &AD{
		Server: adserver,
		Port:   port,
		BaseDN: basedn,
	}
}
