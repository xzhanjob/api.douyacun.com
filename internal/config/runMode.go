package config

const (
	// DebugMode indicates gin mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates gin mode is release.
	ReleaseMode = "release"
	// TestMode indicates gin mode is test.
	TestMode = "test"
)
const (
	debugCode = iota
	releaseCode
	testCode
)

var runMode = debugCode

func SetRunMode(mode string) {
	switch mode {
	case "":
		runMode = debugCode
	case ReleaseMode:
		runMode = releaseCode
	case TestMode:
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
