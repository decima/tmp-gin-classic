package configuration

type ServerConfiguration struct {
	HostAndPort string
	Security    ServerSecurityConfiguration
}

type ServerSecurityConfiguration struct {
	Enable                      bool
	RedirectOnUnauthorizedPath  string
	RedirectOnLogin             string
	UsePreviousIfDefinedOnLogin bool
}
