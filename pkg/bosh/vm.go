package bosh

type VM struct {
	InstanceName string
	Deployment   string
	JobName      string
	JobState     string
	IPs          []string
}
