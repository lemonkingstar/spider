package iserver

var (
	AppName     = "unknown"
	Version     = "unknown"
	BuildBranch = ""
	BuildCommit = ""
	BuildTime   = ""
)

func GetAppName() string     { return AppName }
func GetVersion() string     { return Version }
func GetBuildBranch() string { return BuildBranch }
func GetBuildCommit() string { return BuildCommit }
func GetBuildTime() string   { return BuildTime }
