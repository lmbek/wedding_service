package env

type Env interface {
	HttpPort() string
	HttpsPort() string
	MySQLHost() string
	MySQLPort() string
	MySQLUser() string
	MySQLPassword() string
	MySQLDatabase() string
	MySQLRootPassword() string
}

type env struct {
	httpPort  string
	httpsPort string
	mysqlHost string
	mysqlPort string
	mysqlUser string
	mysqlPass string
	mysqlDB   string
	mysqlRoot string
}

// NewEnv reads configuration exclusively from the process environment.
// Required variables: DEBUG, WEDDING_SERVICE_HTTP_PORT.
// Optional variables:
//   - WEDDING_SERVICE_HTTPS_PORT (defaults to "8443" if empty)
//   - WEDDING_SERVICE_HOSTNAMES (format: "host:alias1,alias2|host2:alias...")
//   - SELF_SIGNED_CERT_PATH (defaults to "/data/certs/localhost_wedding_service.crt")
//   - SELF_SIGNED_KEY_PATH (defaults to "/data/certs/localhost_wedding_service.key")
func NewEnv(path string) (Env, error) {
	e := &env{
		httpPort:  "80",
		httpsPort: "443",
		//certPath:   os.Getenv("SELF_SIGNED_CERT_PATH"),
		//keyPath:    os.Getenv("SELF_SIGNED_KEY_PATH"),
		mysqlHost: "localhost",
		mysqlPort: "3306",
		mysqlUser: "lmbek",
		mysqlPass: "kp-o34-Aa,e4.FW/EfeKLA2Rt,mfk",
		mysqlDB:   "wedding_db",
		mysqlRoot: "kp-o,e4.g434erw,.-FW/EfwdweK34-AaLA2Rt02912la,mfk",
	}
	return e, nil
}

func (e *env) HttpPort() string          { return e.httpPort }
func (e *env) HttpsPort() string         { return e.httpsPort }
func (e *env) MySQLHost() string         { return e.mysqlHost }
func (e *env) MySQLPort() string         { return e.mysqlPort }
func (e *env) MySQLUser() string         { return e.mysqlUser }
func (e *env) MySQLPassword() string     { return e.mysqlPass }
func (e *env) MySQLDatabase() string     { return e.mysqlDB }
func (e *env) MySQLRootPassword() string { return e.mysqlRoot }
