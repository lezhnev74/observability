package observability

type Config struct {
	Listen struct {
		Protocol   string // tcp, udp, unix
		Host, Port string // 127.0.0.1:55678
		Unix       string // /var/observe.sock
	}
	Clickhouse struct{ Host, Port, Database, Table, Username, Password string }
	Observe    []struct {
		PollIntervalSec int `yaml:"poll_interval_sec"`

		FpmPoolStatsAddr string `yaml:"fpm_addr"`
		FpmPoolStatsPath string `yaml:"fpm_path"`

		Mysql struct{ Host, Port, Username, Password string }
	}
}
