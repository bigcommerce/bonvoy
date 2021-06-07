package commands

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func All() []Runner {
	cmds := []Runner{
		BuildListeners(),
		BuildVersion(),
		BuildExpiredCertificatesCommand(),
	}
	return cmds
}