package battery

import (
	"fmt"
	"sync"
)

type Charger interface {
	SetCharging(float64) error
}

type Controller struct {
	Chargers []Charger
}

func NewController(chargers []Charger) *Controller {
	return &Controller{Chargers: chargers}
}

func (c *Controller) SetValue(val float64) {
	for _, oneCharger := range c.Chargers {
		// TODO: add retry
		// TODO: add waitGroup. In current task it would have been overengineering.
		err := oneCharger.SetCharging(val)
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}

func (c *Controller) SetValueAsync(val float64) {
	var wg *sync.WaitGroup
	for _, oneCharger := range c.Chargers {
		wg.Add(1)
		go func(charger Charger, group *sync.WaitGroup) {
			defer group.Done()
			err := charger.SetCharging(val)
			if err != nil {
				fmt.Println(err)
			}
		}(oneCharger, wg)
	}
	wg.Wait()
	return
}
