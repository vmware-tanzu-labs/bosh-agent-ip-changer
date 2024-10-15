package bosh

type VM struct {
	InstanceName  string
	Deployment    string
	InstanceGroup string
	JobState      string
	IPs           []string
}
