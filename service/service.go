package service

import (
	"fmt"

	"github.com/Drofff/revsynth/aco"
	"github.com/Drofff/revsynth/circuit"
	"github.com/Drofff/revsynth/logging"
)

type AlgorithmConfig struct {
	NumOfAnts              int     `json:"num_of_ants"`
	NumOfIterations        int     `json:"num_of_iterations"`
	Alpha                  float32 `json:"alpha"`
	Beta                   float32 `json:"beta"`
	EvaporationRate        float32 `json:"evaporation_rate"`
	LocalLoops             int     `json:"local_loops"`
	SearchDepth            int     `json:"search_depth"`
	DisableNegativeControl bool    `json:"disable_negative_control"`
}

type TruthTable struct {
	Inputs  [][]int `json:"inputs"`
	Outputs [][]int `json:"outputs"`
}

type SynthesiseInput struct {
	Config AlgorithmConfig `json:"config"`
	Target TruthTable      `json:"target"`
}

type Gate struct {
	TypeName    string `json:"type"`
	ControlBits []int  `json:"control_bits"`
	TargetBits  []int  `json:"target_bits"`
}

type SynthesiseOutput struct {
	ErrorsCount int    `json:"errors_count"`
	Cost        int    `json:"cost"`
	Gates       []Gate `json:"gates"`
}

const (
	acoDepositStrength = 100
	gateTypeToffoli    = "toffoli"
)

func toAcoConfig(ac AlgorithmConfig) aco.Config {
	allowedCV := circuit.ControlBitValues
	if ac.DisableNegativeControl {
		allowedCV = []int{circuit.ControlBitIgnore, circuit.ControlBitPositive}
	}

	return aco.Config{
		NumOfAnts:               ac.NumOfAnts,
		NumOfIterations:         ac.NumOfIterations,
		Alpha:                   float64(ac.Alpha),
		Beta:                    float64(ac.Beta),
		EvaporationRate:         float64(ac.EvaporationRate),
		DepositStrength:         acoDepositStrength,
		LocalLoops:              ac.LocalLoops,
		SearchDepth:             ac.SearchDepth,
		AllowedControlBitValues: allowedCV,
	}
}

func toTruthVector(target TruthTable) circuit.TruthVector {
	tt := circuit.TruthTable{}
	for ri := 0; ri < len(target.Inputs); ri++ {
		row := circuit.TruthTableRow{Input: target.Inputs[ri], Output: target.Outputs[ri]}
		tt.Rows = append(tt.Rows, row)
	}
	return tt.ToVector()
}

func toSynthesisOutput(res aco.SynthesisResult) *SynthesiseOutput {
	gates := make([]Gate, 0)
	for i := len(res.Gates) - 1; i >= 0; i-- {
		resGate := res.Gates[i]
		gates = append(gates, Gate{
			TypeName:    gateTypeToffoli,
			TargetBits:  resGate.TargetBits(),
			ControlBits: resGate.ControlBits(),
		})
	}

	return &SynthesiseOutput{ErrorsCount: res.Complexity, Cost: len(res.Gates), Gates: gates}
}

func Synthesise(in *SynthesiseInput) (*SynthesiseOutput, error) {
	conf := toAcoConfig(in.Config)
	synth := aco.NewSynthesizer(conf, circuit.NewToffoliGateFactory(), logging.NewLogger(logging.LevelInfo))
	res := synth.Synthesise(toTruthVector(in.Target))

	if len(res.Gates) == 0 {
		return nil, fmt.Errorf("failed to synthesize")
	}

	return toSynthesisOutput(res), nil
}
