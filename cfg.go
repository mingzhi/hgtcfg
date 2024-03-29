package hgtcfg

import (
	"bytes"
	"fmt"
)

type ParamSet struct {
	Sizes           []int
	Lengths         []int
	MutationRates   []float64
	TransferRates   []float64
	TransferFrags   []int
	TransferDists   []int
	TransferEffs    []float64
	TransferEffLens []int
	SampleSizes     []int
	SampleTimes     []int
	SampleRepls     []int
	SampleBias      int
	ClusterSize     int
	ClusterNumber   int
	LinkageSize     int
	CovMaxls        []int
	FitnessRates    []float64
	FitnessScales   []float64
	FitnessShapes   []float64
	FitnessCoupled  int
	Model           int
	AlphabetSize    int
	Save            int
}

type Cfg struct {
	Population Population
	Mutation   Mutation
	Transfer   Transfer
	Sample     Sample
	Fitness    Fitness
	Output     Output
	Linkage    Linkage
	Cov        Cov
	Genome     Genome
}

func (c Cfg) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, c.Population)
	fmt.Fprintln(&b, c.Genome)
	fmt.Fprintln(&b, c.Mutation)
	fmt.Fprintln(&b, c.Transfer)
	fmt.Fprintln(&b, c.Sample)
	fmt.Fprintln(&b, c.Fitness)
	fmt.Fprintln(&b, c.Output)
	fmt.Fprintln(&b, c.Linkage)
	fmt.Fprintln(&b, c.Cov)
	return b.String()
}

type Genome struct {
	AlphabetSize int
}

func (g Genome) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "[genome]\n")
	fmt.Fprintf(&b, "alphabet_size = %d\n", g.AlphabetSize)
	return b.String()
}

type Cov struct {
	Maxl int
}

func (c Cov) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "[cov]\n")
	fmt.Fprintf(&b, "maxl = %d\n", c.Maxl)
	return b.String()
}

type Mutation struct {
	Rate float64
}

func (m Mutation) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "[mutation]")
	fmt.Fprintf(&b, "rate = %g\n", m.Rate)
	return b.String()
}

type Transfer struct {
	Rate         float64
	Frag         int
	Distribution int
	Efficiency   float64
	EffLen       int
}

func (t Transfer) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "[transfer]")
	fmt.Fprintf(&b, "rate = %g\n", t.Rate)
	fmt.Fprintf(&b, "fragment = %d\n", t.Frag)
	fmt.Fprintf(&b, "distribution = %d\n", t.Distribution)
	fmt.Fprintf(&b, "efficiency = %g\n", t.Efficiency)
	fmt.Fprintf(&b, "efficiency_length = %d\n", t.EffLen)
	return b.String()

}

type Population struct {
	Size       int
	Length     int
	Model      int
	Generation int
}

func (p Population) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "[population]")
	fmt.Fprintf(&b, "size = %d\n", p.Size)
	fmt.Fprintf(&b, "length = %d\n", p.Length)
	fmt.Fprintf(&b, "model = %d\n", p.Model)
	fmt.Fprintf(&b, "generations = %d\n", p.Generation)
	return b.String()
}

type Sample struct {
	Size          int
	Time          int
	Replicates    int
	ClusterNumber int
	ClusterSize   int
	Bias          int
}

func (s Sample) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "[sample]")
	fmt.Fprintf(&b, "size = %d\n", s.Size)
	fmt.Fprintf(&b, "time = %d\n", s.Time)
	fmt.Fprintf(&b, "replicates = %d\n", s.Replicates)
	fmt.Fprintf(&b, "bias = %d\n", s.Bias)
	fmt.Fprintf(&b, "cluster_size = %d\n", s.ClusterSize)
	fmt.Fprintf(&b, "cluster_number = %d\n", s.ClusterNumber)
	return b.String()
}

type Fitness struct {
	Rate    float64
	Scale   float64
	Shape   float64
	Coupled int
}

type Linkage struct {
	Size int
}

func (l Linkage) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "[linkage]\n")
	fmt.Fprintf(&b, "size = %d\n", l.Size)
	return b.String()
}

type Output struct {
	Prefix string
	Save   int
}

func (o Output) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "[output]")
	fmt.Fprintf(&b, "prefix = %s\n", o.Prefix)
	fmt.Fprintf(&b, "save = %d\n", o.Save)
	return b.String()
}

func (f Fitness) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "[fitness]")
	fmt.Fprintf(&b, "rate = %g\n", f.Rate)
	fmt.Fprintf(&b, "scale = %g\n", f.Scale)
	fmt.Fprintf(&b, "shape = %g\n", f.Shape)
	fmt.Fprintf(&b, "coupled = %d\n", f.Coupled)
	return b.String()
}
