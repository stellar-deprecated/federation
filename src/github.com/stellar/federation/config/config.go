package config

type Config struct {
	Port     int
	Domain   string
	Database ConfigDatabase
	Queries  ConfigQueries
	Tls      struct {
		CertificateFile string `mapstructure:"certificate-file"`
		PrivateKeyFile  string `mapstructure:"private-key-file"`
	}
}

type ConfigDatabase struct {
	Type string
	Url  string
	// Casandra only:
	Cluster      []string
	Keyspace     string
	ProtoVersion *int `mapstructure:"proto-version"`
}

type ConfigQueries struct {
	Federation        string
	ReverseFederation string `mapstructure:"reverse-federation"`
}
