package mars

import (
	"os"
)

// EnvMarsMode indicates environment name for mars mode.
const EnvMarsMode = "MARS_MODE"

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
	// TestMode indicates mode is test.
	TestMode = "test"
)
const (
	debugCode = iota
	releaseCode
	testCode
)

var marsMode = debugCode
var modeName = DebugMode

func init() {
	mode := os.Getenv(EnvMarsMode)
	SetMode(mode)
}

// SetMode sets mode according to input string.
func SetMode(value string) {
	switch value {
	case DebugMode, "":
		marsMode = debugCode
	case ReleaseMode:
		marsMode = releaseCode
	case TestMode:
		marsMode = testCode
	default:
		panic("mode unknown: " + value)
	}
	if value == "" {
		value = DebugMode
	}
	modeName = value
}

// Mode returns currently mode.
func Mode() string {
	return modeName
}
