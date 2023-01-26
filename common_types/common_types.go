package common_types

type PayloadFromMessageTwitch struct {
	UserId        string
	IsBroadcaster bool
	IsVip         bool
	IsMod         bool
	IsSubscriber  bool
	Message       string
}
