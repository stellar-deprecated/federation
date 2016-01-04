package federation

type Config struct {
	Port                   int
	Domain                 string
	DatabaseType           string
	DatabaseUrl            string
	FederationQuery        string
	ReverseFederationQuery string
}
