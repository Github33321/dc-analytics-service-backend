package models

type TaskStat struct {
	CreatedAtStr string `json:"created_at_str"`
	Count        uint64 `json:"count"`
}
