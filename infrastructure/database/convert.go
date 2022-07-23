package database

import (
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/wzzfarewell/go-mod/infrastructure/utils/copyutil"
	"gorm.io/gorm"
)

type ModelConverter[MODEL, DOMAIN any] struct {
}

func (c *ModelConverter[MODEL, DOMAIN]) ToDomain(m *MODEL) (*DOMAIN, error) {
	return copyutil.Copy[MODEL, DOMAIN](m)
}

func (c *ModelConverter[MODEL, DOMAIN]) ToDomainWithError(m *MODEL, err error) (*DOMAIN, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return c.ToDomain(m)
}

func (c *ModelConverter[MODEL, DOMAIN]) ToDomains(ms []*MODEL) ([]*DOMAIN, error) {
	return copyutil.CopySlice[MODEL, DOMAIN](ms)
}

func (c *ModelConverter[MODEL, DOMAIN]) ToModel(d *DOMAIN) (*MODEL, error) {
	return copyutil.Copy[DOMAIN, MODEL](d)
}

func (c *ModelConverter[MODEL, DOMAIN]) ToModels(ds []*DOMAIN) ([]*MODEL, error) {
	return copyutil.CopySlice[DOMAIN, MODEL](ds)
}

func (c *ModelConverter[MODEL, DOMAIN]) ToModelFrom(d *DOMAIN, m *MODEL) error {
	return copyutil.CopyTo(d, m)
}

func (c *ModelConverter[MODEL, DOMAIN]) ToDomainFrom(m *MODEL, d *DOMAIN) error {
	return copyutil.CopyTo(m, d)
}

func (c *ModelConverter[MODEL, DOMAIN]) ModelToMap(m *MODEL) (map[string]any, error) {
	if !structs.IsStruct(m) {
		return nil, errors.New("model is not a struct")
	}
	s := structs.New(m)
	s.TagName = "json"
	return s.Map(), nil
}

func (c *ModelConverter[MODEL, DOMAIN]) DomainToMap(d *DOMAIN) (map[string]any, error) {
	if !structs.IsStruct(d) {
		return nil, errors.New("domain is not a struct")
	}
	s := structs.New(d)
	s.TagName = "json"
	return s.Map(), nil
}
