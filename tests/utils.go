package tests

import (
	"errors"
	"io"
	"reflect"
	"regexp"
	"time"
)

// No one should be using func Main anymore.
// See the doc comment on func Main and use MainStart instead.
var errMain = errors.New("testing: unexpected use of func Main")

type corpusEntry = struct {
	Parent     string
	Path       string
	Data       []byte
	Values     []any
	Generation int
	IsSeed     bool
}

type MatchStringOnly struct{}

func (f MatchStringOnly) StartCPUProfile(w io.Writer) error           { return errMain }
func (f MatchStringOnly) StopCPUProfile()                             {}
func (f MatchStringOnly) WriteProfileTo(string, io.Writer, int) error { return errMain }
func (f MatchStringOnly) ImportPath() string                          { return "" }
func (f MatchStringOnly) StartTestLog(io.Writer)                      {}
func (f MatchStringOnly) StopTestLog() error                          { return errMain }
func (f MatchStringOnly) SetPanicOnExit0(bool)                        {}
func (f MatchStringOnly) CheckCorpus([]any, []reflect.Type) error     { return nil }
func (f MatchStringOnly) ResetCoverage()                              {}
func (f MatchStringOnly) SnapshotCoverage()                           {}
func (f MatchStringOnly) CoordinateFuzzing(time.Duration, int64, time.Duration, int64, int, []corpusEntry, []reflect.Type, string, string) error {
	return errMain
}
func (f MatchStringOnly) RunFuzzWorker(func(corpusEntry) error) error { return errMain }
func (f MatchStringOnly) ReadCorpus(string, []reflect.Type) ([]corpusEntry, error) {
	return nil, errMain
}

var matchPat string
var matchRe *regexp.Regexp

func (f MatchStringOnly) MatchString(pat, str string) (result bool, err error) {
	if matchRe == nil || matchPat != pat {
		matchPat = pat
		matchRe, err = regexp.Compile(matchPat)
		if err != nil {
			return
		}
	}
	return matchRe.MatchString(str), nil
}
