package validate

type ChannelCreateValidator struct {
	Title   string   `validate:"" json:"title"`
	Members []string `validate:"required" json:"members"`
	Type    string   `validate:"required,oneof=private public" json:"type"`
}
