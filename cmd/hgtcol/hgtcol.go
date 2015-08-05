package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/mingzhi/hgtcfg"
	"os"
	"strconv"
)

var prefix string

func init() {
	flag.StringVar(&prefix, "prefix", "test", "prefix")
	flag.Parse()
}

func main() {
	originalCfgs := parse(prefix)
	m := make(map[Cfg][]string)
	for _, cfg := range originalCfgs {
		p := cfg.Output.Prefix
		// remove unique prefix.
		cfg.Output.Prefix = ""
		m[cfg] = append(m[cfg], p)
	}

	processedCfgs := []Cfg{}
	index := 0
	for cfg, prefixes := range m {
		// add unique prefix.
		cfg.Output.Prefix = fmt.Sprintf("%s_merge_%d", prefix, index)
		processedCfgs = append(processedCfgs, cfg)
		// write merged cfg.
		writeCfg(cfg, cfg.Output.Prefix+".cfg.json")

		// merge Ks.
		allKs := []Ks{}
		for _, p := range prefixes {
			ksArr := readKs(p)
			allKs = append(allKs, ksArr...)
		}
		// write Ks.
		writeKs(allKs, cfg.Output.Prefix+".ks.json")

		// merge t2.
		allT2 := []T{}
		for _, p := range prefixes {
			allT2 = append(allT2, readT2(p)...)
		}
		// write T2.
		writeJSON(allT2, cfg.Output.Prefix+".t2.json")

		index++
	}

	writeJSON(processedCfgs, prefix+"_merge_configs.json")
}

func parse(prefix string) (cfgs []Cfg) {
	filename := prefix + "_configs.json"
	// open file.
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// create json reader.
	d := json.NewDecoder(f)
	if err := d.Decode(&cfgs); err != nil {
		panic(err)
	}

	return
}

type Ks struct {
	Ks, Vd float64
	G      int
}

func readKs(prefix string) (ksArr []Ks) {
	filename := prefix + ".ks.txt"
	records := readTable(filename)
	for i := 0; i < len(records); i++ {
		cols := records[i]
		ks := stringToFloat64(cols[0])
		vd := stringToFloat64(cols[1])
		g := stringToInt(cols[5])
		ksArr = append(ksArr, Ks{Ks: ks, Vd: vd, G: g})
	}
	return
}

func writeKs(ksArr []Ks, filename string) {
	writeJSON(ksArr, filename)
}

type T struct {
	T float64
	G int
}

func readT2(prefix string) (ts []T) {
	filename := prefix + ".t2.txt"
	records := readTable(filename)
	for i := 0; i < len(records); i++ {
		cols := records[i]
		t := stringToFloat64(cols[0])
		g := stringToInt(cols[1])
		ts = append(ts, T{T: t, G: g})
	}

	return ts
}

func readTable(filename string) (records [][]string) {
	// open file.
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = '\t'
	reader.Comment = '#'

	records, err = reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return
}

func stringToFloat64(s string) (v float64) {
	var err error
	v, err = strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return
}

func stringToInt(s string) (v int) {
	var err error
	v, err = strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return
}

func writeCfg(cfg Cfg, filename string) {
	writeJSON(cfg, filename)
}

func writeJSON(o interface{}, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e := json.NewEncoder(f)
	if err := e.Encode(o); err != nil {
		panic(err)
	}
}
