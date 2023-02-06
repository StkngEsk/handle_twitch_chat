package types

type GetUserTwitchResponse struct {
	Data []DataTwitch `json:"data"`
}

type DataTwitch struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	ViewCount       int64  `json:"view_count"`
	CreatedAt       string `json:"created_at"`
}
