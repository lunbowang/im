package reply

import "time"

type ParamPublishFile struct {
	ID       int64     `json:"id,omitempty"`
	FileType string    `json:"file_type,omitempty"`
	FileSize int64     `json:"file_size,omitempty"`
	Url      string    `json:"url,omitempty"`
	CreateAt time.Time `json:"create_at"`
}
