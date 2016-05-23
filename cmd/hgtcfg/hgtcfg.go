package main

import (
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/mingzhi/hgtcfg"
	"os"
)

var cfgFileName string
var prefix string
var nodes int
var ppn int
var walltime int
var factor int
var message string
var replicates int
var mpirun bool
var exec string

func init() {
	flag.StringVar(&prefix, "prefix", "test", "prefix")
	flag.StringVar(&message, "message", "a", "message")
	flag.IntVar(&nodes, "nodes", 1, "nodes")
	flag.IntVar(&ppn, "ppn", 1, "ppn")
	flag.IntVar(&walltime, "walltime", 48, "walltime in hours")
	flag.IntVar(&factor, "factor", 1, "factor multiple to generations")
	flag.IntVar(&replicates, "r", 1, "replicates")
	flag.BoolVar(&mpirun, "mpi", false, "mpi run?")
	flag.StringVar(&exec, "exec", "hgt_simu", "exec name")
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
							for _, transferEffLen := range ps.TransferEffLens {
								for _, transferEff := range ps.TransferEffs {
									for _, sampleSize := range ps.SampleSizes {
										for _, sampleTime := range ps.SampleTimes {
											for _, sampleRepl := range ps.SampleRepls {
												for _, sampleMaxl := range ps.CovMaxls {
													for _, fitnessRate := range ps.FitnessRates {
														fitnessScales := ps.FitnessScales
														for _, fitnessScale := range fitnessScales {
															for _, fitnessShape := range ps.FitnessShapes {
																// create population
																if fitnessRate == 0 {
																	if fitnessScale != 0 {
																		continue
																	}
																}
																if fitnessScale == 0 {
																	if fitnessRate != 0 {
																		continue
																	}
																}
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
																	Efficiency:   transferEff,
																	EffLen:       transferEffLen,
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

																lin := Linkage{
																	Size: 0,
																}

																genome := Genome{}
																genome.AlphabetSize = ps.AlphabetSize

																cfg := Cfg{
																	Population: pop,
																	Mutation:   mut,
																	Transfer:   tra,
																	Sample:     smp,
																	Fitness:    fit,
																	Linkage:    lin,
																	Cov:        cov,
																	Genome:     genome,
																}
																for i := 0; i < replicates; i++ {
																	cfg.Output.Prefix = fmt.Sprintf("%s_individual_%d", prefix, count)
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
	w.WriteString(fmt.Sprintf("cd %s\n", wd))
	if mpirun {
		w.WriteString(fmt.Sprintf("module load openmpi/intel/1.6.5\n"))
		w.WriteString(fmt.Sprintf("mpirun -n %d %s -C %s\n", ppn*nodes, exec, c.Output.Prefix+".cfg.ini"))
	} else {
		w.WriteString(fmt.Sprintf("%s -C %s\n", exec, c.Output.Prefix+".cfg.ini"))
	}

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
