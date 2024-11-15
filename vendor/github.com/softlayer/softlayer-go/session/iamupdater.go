package session

//counterfeiter:generate . IAMUpdater
type IAMUpdater interface {
	Update(token string, refresh string)
}

type LogIamUpdater struct {
	debug bool
}

func NewLogIamUpdater(debug bool) *LogIamUpdater {
	return &LogIamUpdater{
		debug: debug,
	}
}

func (iamupdater *LogIamUpdater) Update(token string, refresh string) {
	if iamupdater.debug {
		Logger.Printf("[DEBUG] New Token: %s\n", token)
		Logger.Printf("[DEBUG] New Refresh Token: %s\n", refresh)
	}
}
