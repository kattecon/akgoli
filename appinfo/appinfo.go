// Package appinfo exposes application identity and version metadata. The
// version and idName variables are populated at build time through Go linker
// flags. When they are not set, the corresponding methods return "unknown".
//
// Build example:
//
//	go build -ldflags="-X 'github.com/kattecon/akgoli/appinfo.version=1.0.0' -X 'github.com/kattecon/akgoli/appinfo.idName=myapp'"
package appinfo

import (
	"runtime"
	"strings"
)

// AppInfo provides read-only access to the application name, version, and
// Go runtime version. Use Get for production values and Mock for tests.
type AppInfo interface {
	// AppIdName returns the application identifier set at build time, or
	// "unknown" if the linker flag was not provided.
	AppIdName() string

	// AppVersion returns the application version set at build time, or
	// "unknown" if the linker flag was not provided.
	AppVersion() string

	// GoVersion returns the Go runtime version with the "go" prefix stripped.
	// For example, Go 1.25.3 produces "1.25.3". Returns "unknown" if the
	// runtime version string is empty.
	GoVersion() string
}

type appInfoImpl struct{}

var (
	// Populated by -ldflags at build time. Not set through normal Go
	// assignments.
	version string
	idName  string

	goVersion string = strings.TrimPrefix(runtime.Version(), "go")

	impl appInfoImpl = appInfoImpl{}
)

// Get returns the production AppInfo backed by linker-injected values.
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
