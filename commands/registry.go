package commands

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func All() []Runner {
	return []Runner{
		BuildListeners(),
		BuildVersion(),
		BuildExpiredCertificatesCommand(),
		BuildSetLogLevelCommand(),
	}
}