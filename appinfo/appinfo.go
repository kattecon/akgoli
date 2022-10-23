package appinfo

import (
	"runtime"
	"strings"
)

type AppInfo interface {
	AppIdName() string
	AppVersion() string
	GoVersion() string
}

type appInfoImpl struct{}

var (
	version string
	idName  string

	goVersion string = strings.TrimPrefix(runtime.Version(), "go")

	impl appInfoImpl = appInfoImpl{}
)

func Get() AppInfo {
	return impl
}

func (i appInfoImpl) AppIdName() string {
	if idName == "" {
		return "unknown"
	}
	return idName
}

func (i appInfoImpl) AppVersion() string {
	if version == "" {
		return "unknown"
	}
	return version
}

func (i appInfoImpl) GoVersion() string {
	if goVersion == "" {
		return "unknown"
	}
	return goVersion
}
