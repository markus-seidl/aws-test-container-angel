package angel

type Configuration struct {
	AngelDirectory string
	TempDirectory  string
}

func DefaultConfiguration() Configuration {
	return Configuration{
		AngelDirectory: "/angel",
		TempDirectory:  "/tmp",
	}
}
