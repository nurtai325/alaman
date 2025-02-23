package timezone

import "os"

const (
	defaultZone = "Asia/Yekaterinburg"
)

func init() {
	os.Setenv("TZ", defaultZone)
}

func Get() string {
	return defaultZone
}
