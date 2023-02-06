package types

type PayloadFromMessageTwitch struct {
	UserId        string
	UserName      string
	IsBroadcaster bool
	IsVip         bool
	IsMod         bool
	IsSubscriber  bool
	Emotes        string
	Message       string
}

type PayloadToClient struct {
	IsBroadcaster bool    `json:"isBroadcaster"`
	UserId        string  `json:"userId"`
	Username      string  `json:"username"`
	DisplayName   string  `json:"displayName"`
	UrlUserImage  string  `json:"urlUserImage"`
	Emotes        string  `json:"emotes"`
	IsMod         bool    `json:"isMod"`
	Message       string  `json:"message"`
	SvgGlobe      string  `json:"svgGlobe"`
	CPGlobe       string  `json:"cPGlobe"`
	CSGlobe       string  `json:"cSGlobe"`
	OpacityGlobe  float64 `json:"opacityGlobe"`
}
