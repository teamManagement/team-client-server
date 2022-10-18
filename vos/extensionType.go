package vos

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// ReleaseDigest 发布摘要说明
type ReleaseDigest []string

func (s ReleaseDigest) Value() (driver.Value, error) {
	marshal, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return string(marshal), nil
}

func (s *ReleaseDigest) Scan(src any) error {
	jsonStr, ok := src.(string)
	if !ok {
		return errors.New("轮播图原始数据类型不正确")
	}

	return json.Unmarshal([]byte(jsonStr), &s)
}
