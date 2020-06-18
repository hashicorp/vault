package cos

import (
	"encoding/json"
)

type PicOperations struct {
	IsPicInfo int                  `json:"is_pic_info,omitempty"`
	Rules     []PicOperationsRules `json:"rules,omitemtpy"`
}

type PicOperationsRules struct {
	Bucket string `json:"bucket,omitempty"`
	FileId string `json:"fileid"`
	Rule   string `json:"rule"`
}

func EncodePicOperations(pic *PicOperations) string {
	bs, err := json.Marshal(pic)
	if err != nil {
		return ""
	}
	return string(bs)
}
