package ers

type EP struct {
	ERSEndPoint struct {
		ID                      string `json:"id"`
		Name                    string `json:"name"`
		Description             string `json:"description"`
		Mac                     string `json:"mac"`
		ProfileID               string `json:"profileId"`
		StaticProfileAssignment bool   `json:"staticProfileAssignment"`
		GroupID                 string `json:"groupId"`
		StaticGroupAssignment   bool   `json:"staticGroupAssignment"`
		PortalUser              string `json:"portalUser"`
		IdentityStore           string `json:"identityStore"`
		IdentityStoreID         string `json:"identityStoreId"`
		Link                    struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
			Type string `json:"type"`
		} `json:"link"`
	} `json:"ERSEndPoint"`
}

type EPGroup struct {
	EndPointGroup struct {
		Description   string `json:"description"`
		ID            string `json:"id"`
		Name          string `json:"name"`
		SystemDefined bool   `json:"systemDefined"`
	} `json:"EndPointGroup"`
}
