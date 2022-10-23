package appinfo

type mockImpl struct{}

var mock mockImpl = mockImpl{}

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
