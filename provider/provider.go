package provider

//Provider represents a cloud provider.
type Provider interface {
	GetInstanceMap(group *string) (map[string]string, error)
}

//Config configures a Provider.
type Config struct {
	Region      string
	UsePublicIP bool
}
