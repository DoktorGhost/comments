package types

type Comment struct {
	ID              int    `json:"id"`
	NewsID          int    `json:"news_id"`
	Text            string `json:"text"`
	ParentCommentID int    `json:"parent_id"`
}
