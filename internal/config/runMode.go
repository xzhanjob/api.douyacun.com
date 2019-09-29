package config

const (
	// DebugMode indicates gin mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates gin mode is release.
	ReleaseMode = "release"
	// TestMode indicates gin mode is test.
	InfoMode = "info"
)
const (
	debugCode = iota
	releaseCode
	testCode
)

var runMode = debugCode

func SetRunMode(mode string) {
	switch mode {
	case DebugMode:
		runMode = debugCode
	case ReleaseMode:
		runMode = releaseCode
	case InfoMode:
		runMode = testCode
	}
}

func IsRelease() bool {
	return runMode == releaseCode
}

func IsDebug() bool {
	return runMode == debugCode
}

func IsTest() bool {
	return runMode == testCode
}
