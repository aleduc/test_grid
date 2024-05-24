package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"

	"test_grid/internal/grid"
)

const (
	maxRetries = 10
	retryPause = time.Second
)

var DefaultActionMap = map[grid.Intensity]float64{
	grid.VeryHigh: -1,
	grid.High:     -1,
	grid.Low:      1,
	grid.VeryLow:  1,
}

type ChargerController interface {
	SetValue(val float64)
}

type NationalGrid interface {
	GetIntensity(ctx context.Context) (val grid.Intensity, err error)
}

type Processor struct {
	ChargerController ChargerController
	NationalGrid      NationalGrid
	ParentCtx         context.Context
	ActionMap         map[grid.Intensity]float64
	backOff           backoff.BackOff
}

func NewProcessor(chargerController ChargerController, nationalGrid NationalGrid, parentCtx context.Context, actionMap map[grid.Intensity]float64) *Processor {
	return &Processor{
		ChargerController: chargerController,
		NationalGrid:      nationalGrid,
		ParentCtx:         parentCtx,
		// Just default, out of the task's scope.
		backOff:   backoff.WithMaxRetries(backoff.WithContext(backoff.NewConstantBackOff(retryPause), parentCtx), maxRetries),
		ActionMap: actionMap,
	}
}

func (p *Processor) Run() {
	p.backOff.Reset()
	var intensity grid.Intensity
	operation := func() (err error) {
		intensity, err = p.NationalGrid.GetIntensity(p.ParentCtx)
		if err != nil {
			fmt.Println(fmt.Errorf("[retry in %s] GetIntensity error: %w", p.backOff.NextBackOff(), err))
		}
		return err

	}

	err := backoff.Retry(operation, p.backOff)
	if err != nil {
		return
	}

	if val, ok := p.ActionMap[intensity]; ok {
		p.ChargerController.SetValue(val)
	} else {
		fmt.Println(fmt.Sprintf("there is not action for intensity: %s", intensity))
	}
}
