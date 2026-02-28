module github.com/kattecon/akgoli

go 1.25.3

require (
	github.com/prometheus/common v0.67.5
	github.com/stretchr/testify v1.11.1
	golang.org/x/crypto v0.48.0
)

require (
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/tools/go/expect v0.1.1-deprecated // indirect
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/alecthomas/kingpin v2.2.6+incompatible // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/caarlos0/svu v1.12.0 // indirect
	github.com/cespare/reflex v0.3.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/creack/pty v1.1.11 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/ogier/pflag v0.0.1 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/ramya-rao-a/go-outline v0.0.0-20210608161538-9736a4bde949 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	golang.org/x/exp/typeparams v0.0.0-20250620022241-b7579e27df2b // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/telemetry v0.0.0-20260109210033-bd525da824e2 // indirect
	golang.org/x/tools v0.41.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	honnef.co/go/tools v0.6.1 // indirect
)

require (
	github.com/agnivade/levenshtein v1.2.1
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.23.2
	go.uber.org/zap v1.27.1
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// NOTE: Do NOT add golang.org/x/tools/gopls here! It depends (at least it did that at time of writing) on a dev version
// of honnef.co/go/tools which breaks Renovate with weird conflicting deps.
tool (
	github.com/caarlos0/svu
	github.com/cespare/reflex
	github.com/ramya-rao-a/go-outline
	golang.org/x/tools/cmd/goimports
	honnef.co/go/tools/cmd/staticcheck
)
