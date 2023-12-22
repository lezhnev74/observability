package observability

type Mysql struct{ Host, Port, Username, Password string }

type Config struct {
	Listen struct {
		Protocol   string // tcp, udp, unix
		Host, Port string // 127.0.0.1:55678
		Unix       string // /var/observe.sock
	}
	Clickhouse struct{ Host, Port, Database, Username, Password string }
	Observe    struct {
		PollIntervalSec int
		FpmPoolStatsUrl string
		Mysql
	}
}
