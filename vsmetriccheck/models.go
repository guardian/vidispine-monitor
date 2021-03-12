package vsmetriccheck

import (
	"errors"
	"fmt"
	"strconv"
)

type MetricGauge struct {
	Value interface{} `json:"value"`
}

/**
try to convert the gauge value into a decimal
*/
func (g MetricGauge) FloatValue() (float64, error) {
	if floatVal, isFloat := g.Value.(float64); isFloat {
		return floatVal, nil
	}
	if stringVal, isString := g.Value.(string); isString {
		return strconv.ParseFloat(stringVal, 64)
	}
	return 0, errors.New("gauge did not contain a string or float")
}

/**
returns a float64 of the value or panics
*/
func (g MetricGauge) MustFloat() float64 {
	v, err := g.FloatValue()
	if err != nil {
		errString := fmt.Sprintf("%s is not a float value", g.Value)
		panic(errString)
	}
	return v
}

type MetricCounter struct {
	Count int64 `json:"count"`
}

type MetricMeter struct {
	Count    int64   `json:"count"`
	M15Rate  float64 `json:"m15_rate"`
	M1Rate   float64 `json:"m1_rate"`
	M5Rate   float64 `json:"m5_rate"`
	MeanRate float64 `json:"mean_rate"`
	Units    string  `json:"units"`
}

type MetricsResponse struct {
	Version  string                   `json:"version"`
	Gauges   map[string]MetricGauge   `json:"gauges"`
	Counters map[string]MetricCounter `json:"counters"`
	Meters   map[string]MetricMeter   `json:"meters"`
}
