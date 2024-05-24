package battery

import (
	"fmt"
	"time"
)

type Standard struct {
	Name string
}

func NewStandard(name string) *Standard {
	return &Standard{Name: name}
}

func (s *Standard) SetCharging(v float64) error {
	fmt.Println(fmt.Sprintf("%s:Battery %s charging value %v", time.Now().UTC().String(), s.Name, v))
	return nil
}
