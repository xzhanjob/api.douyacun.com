package article

import (
	"context"
	"dyc/internal/db"
	"dyc/internal/logger"
	"encoding/json"
)

func NewInfo(id string) (*Article, error) {
	resp, err := db.ES.Get().Index(TopicCost).Id(id).Do(context.Background())
	if err != nil {
		return nil, nil
	}
	var data Article
	if err = json.Unmarshal(resp.Source, &data); err != nil {
		return nil, err
	}
	if resp.Source == nil {
		return nil, nil
	}
	logger.Debugf("%s", resp.Id)
	return &data, nil
}
