package federation

type Config struct {
	Port     int
	Domain   string
	Database struct {
		Type string
		Url  string
	}
	Queries ConfigQueries
	Tls     struct {
		CertificateFile string `mapstructure:"certificate-file"`
		PrivateKeyFile  string `mapstructure:"private-key-file"`
	}
}

type ConfigQueries struct {
	Federation        string
	ReverseFederation string `mapstructure:"reverse-federation"`
}
