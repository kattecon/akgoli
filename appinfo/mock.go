package appinfo

type mockImpl struct{}

var mock mockImpl = mockImpl{}

// Mock returns an AppInfo with fixed values for tests: AppIdName returns
// "mock", AppVersion returns "1.2.3", and GoVersion returns "100.500". These
// values do not change between calls or test runs.
func Mock() AppInfo {
	return mock
}

func (i mockImpl) AppIdName() string {
	return "mock"
}

func (i mockImpl) AppVersion() string {
	return "1.2.3"
}

func (i mockImpl) GoVersion() string {
	return "100.500"
}
