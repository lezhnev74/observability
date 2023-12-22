package observability

type Config struct {
	Listen struct {
		Host, Port string // 127.0.0.1:55678 (udp)
	}
	Clickhouse struct{ Host, Port, Database, Table, Username, Password string }
	Observe    []struct {
		PollIntervalSec int `yaml:"poll_interval_sec"`

		FpmPoolStatsAddr string `yaml:"fpm_addr"`
		FpmPoolStatsPath string `yaml:"fpm_path"`

		Mysql struct{ Host, Port, Username, Password string }
	}
}
