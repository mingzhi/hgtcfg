package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type ParamSet struct {
	Sizes          []int
	Lengths        []int
	MutationRates  []float64
	TransferRates  []float64
	TransferFrags  []int
	TransferDists  []int
	SampleSizes    []int
	SampleTimes    []int
	SampleRepls    []int
	CovMaxls       []int
	FitnessRates   []float64
	FitnessScales  []float64
	FitnessShapes  []float64
	FitnessCoupled int
	Model          int
	AlphabetSize   int
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
}

func (t Transfer) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "[transfer]")
	fmt.Fprintf(&b, "rate = %g\n", t.Rate)
	fmt.Fprintf(&b, "fragment = %d\n", t.Frag)
	fmt.Fprintf(&b, "distribution = %d\n", t.Distribution)
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
	Size       int
	Time       int
	Replicates int
}

func (s Sample) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "[sample]")
	fmt.Fprintf(&b, "size = %d\n", s.Size)
	fmt.Fprintf(&b, "time = %d\n", s.Time)
	fmt.Fprintf(&b, "replicates = %d\n", s.Replicates)
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
}

func (o Output) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "[output]")
	fmt.Fprintf(&b, "prefix = %s\n", o.Prefix)
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

var cfgFileName string
var prefix string
var nodes int
var ppn int
var walltime int
var factor int
var message string
var replicates int

func init() {
	flag.StringVar(&prefix, "prefix", "test", "prefix")
	flag.StringVar(&message, "message", "a", "message")
	flag.IntVar(&nodes, "nodes", 1, "nodes")
	flag.IntVar(&ppn, "ppn", 20, "ppn")
	flag.IntVar(&walltime, "walltime", 48, "walltime in hours")
	flag.IntVar(&factor, "factor", 1, "factor multiple to generations")
	flag.IntVar(&replicates, "r", 1, "replicates")
	flag.Parse()
	if flag.NArg() <= 0 {
		fmt.Println("need config file!")
		os.Exit(1)
	}
	cfgFileName = flag.Arg(0)
}

func main() {
	ps := parse(cfgFileName)
	cs := create(ps, prefix)
	writeCfgs(cs)
	writeQSub(cs)
	for _, c := range cs {
		writeCfgJson(c)
		writeCfgIni(c)
		writePbs(c)
	}
}

func parse(filename string) (c ParamSet) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d := json.NewDecoder(f)
	if err := d.Decode(&c); err != nil {
		panic(err)
	}
	return
}

func create(ps ParamSet, prefix string) (cs []Cfg) {
	count := 0
	for _, size := range ps.Sizes {
		for _, length := range ps.Lengths {
			for _, mutationRate := range ps.MutationRates {
				for _, transferRate := range ps.TransferRates {
					fragments := ps.TransferFrags
					if transferRate == 0 {
						fragments = []int{0}
					}
					for _, transferFrag := range fragments {
						for _, transferDist := range ps.TransferDists {
							for _, sampleSize := range ps.SampleSizes {
								for _, sampleTime := range ps.SampleTimes {
									for _, sampleRepl := range ps.SampleRepls {
										for _, sampleMaxl := range ps.CovMaxls {
											for _, fitnessRate := range ps.FitnessRates {
												fitnessScales := ps.FitnessScales
												if fitnessRate == 0 {
													fitnessScales = []float64{0}
												}
												for _, fitnessScale := range fitnessScales {
													for _, fitnessShape := range ps.FitnessShapes {
														// create population
														pop := Population{
															Size:   size,
															Length: length,
															Model:  ps.Model,
														}

														if ps.Model > 0 {
															pop.Generation = pop.Size * factor
														} else {
															pop.Generation = pop.Size * pop.Size * factor
														}

														// create mutation
														mut := Mutation{
															Rate: mutationRate,
														}

														// create transfer
														tra := Transfer{
															Rate:         transferRate,
															Frag:         transferFrag,
															Distribution: transferDist,
														}

														// create sample
														smp := Sample{
															Size:       sampleSize,
															Time:       sampleTime,
															Replicates: sampleRepl,
														}

														cov := Cov{
															Maxl: sampleMaxl,
														}

														// create fitness
														fit := Fitness{
															Rate:    fitnessRate,
															Scale:   fitnessScale,
															Shape:   fitnessShape,
															Coupled: ps.FitnessCoupled,
														}

														out := Output{
															Prefix: fmt.Sprintf("%s_individual_%d", prefix, count),
														}

														lin := Linkage{
															Size: transferFrag,
														}

														genome := Genome{}
														genome.AlphabetSize = ps.AlphabetSize

														cfg := Cfg{
															Population: pop,
															Mutation:   mut,
															Transfer:   tra,
															Sample:     smp,
															Fitness:    fit,
															Output:     out,
															Linkage:    lin,
															Cov:        cov,
															Genome:     genome,
														}
														for i := 0; i < replicates; i++ {
															cs = append(cs, cfg)
															count++
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func writeCfgJson(cfg Cfg) {
	filename := cfg.Output.Prefix + ".cfg.json"
	writeJson(cfg, filename)
}

func writeJson(s interface{}, filename string) {
	w, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	e := json.NewEncoder(w)
	if err := e.Encode(s); err != nil {
		panic(err)
	}
}

func writeCfgIni(cfg Cfg) {
	filename := cfg.Output.Prefix + ".cfg.ini"
	w, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	w.WriteString(fmt.Sprintf("%s", cfg))
}

func writeCfgs(cfgs []Cfg) {
	filename := prefix + "_configs.json"
	writeJson(cfgs, filename)
}

func writePbs(c Cfg) {
	filename := c.Output.Prefix + ".pbs"
	w, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	w.WriteString("#!/bin/bash\n")
	w.WriteString(fmt.Sprintf("#PBS -N %s\n", c.Output.Prefix))
	w.WriteString(fmt.Sprintf("#PBS -l nodes=%d:ppn=%d\n", nodes, ppn))
	w.WriteString(fmt.Sprintf("#PBS -l walltime=%d:00:00\n", walltime))
	w.WriteString(fmt.Sprintf("#PBS -M ml3365@nyu.edu\n"))
	w.WriteString(fmt.Sprintf("#PBS -m %s\n", message))
	w.WriteString(fmt.Sprintf("module load openmpi/intel/1.6.5\n"))
	w.WriteString(fmt.Sprintf("cd %s\n", wd))
	w.WriteString(fmt.Sprintf("mpirun -n %d hgt_mpi_moran_const -C %s\n", ppn*nodes, c.Output.Prefix+".cfg.ini"))
}

func writeQSub(cfgs []Cfg) {
	filename := prefix + "_qsub.sh"
	w, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	w.WriteString("#!/bin/bash\n")
	for _, c := range cfgs {
		w.WriteString(fmt.Sprintf("qsub %s.pbs\n", c.Output.Prefix))
	}
}
