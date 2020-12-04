package watchdog

import (
	"testing"

	"github.com/raulk/clock"
	"github.com/stretchr/testify/require"
)

var (
	watermarks = []float64{0.50, 0.75, 0.80}
	thresholds = func() []uint64 {
		var ret []uint64
		for _, w := range watermarks {
			ret = append(ret, uint64(float64(limit)*w))
		}
		return ret
	}()
)

func TestProgressiveWatermarks(t *testing.T) {
	clk := clock.NewMock()
	Clock = clk

	p, err := NewWatermarkPolicy(watermarks...)(limit)
	require.NoError(t, err)

	// at zero
	next, immediate := p.Evaluate(UtilizationSystem, uint64(0))
	require.False(t, immediate)
	require.EqualValues(t, thresholds[0], next)

	// before the watermark.
	next, immediate = p.Evaluate(UtilizationSystem, uint64(float64(limit)*watermarks[0])-1)
	require.False(t, immediate)
	require.EqualValues(t, thresholds[0], next)

	// exactly at the watermark; gives us the next watermark, as the watchdodg would've
	// taken care of triggering the first watermark.
	next, immediate = p.Evaluate(UtilizationSystem, uint64(float64(limit)*watermarks[0]))
	require.False(t, immediate)
	require.EqualValues(t, thresholds[1], next)

	// after the watermark gives us the next watermark.
	next, immediate = p.Evaluate(UtilizationSystem, uint64(float64(limit)*watermarks[0])+1)
	require.False(t, immediate)
	require.EqualValues(t, thresholds[1], next)

	// last watermark; always triggers.
	next, immediate = p.Evaluate(UtilizationSystem, uint64(float64(limit)*watermarks[2]))
	require.True(t, immediate)
	require.EqualValues(t, uint64(float64(limit)*watermarks[2]), next)

	next, immediate = p.Evaluate(UtilizationSystem, uint64(float64(limit)*watermarks[2]+1))
	require.True(t, immediate)
	require.EqualValues(t, uint64(float64(limit)*watermarks[2])+1, next)

	next, immediate = p.Evaluate(UtilizationSystem, limit)
	require.True(t, immediate)
	require.EqualValues(t, limit, next)
}
