package vm

import "slices"

type Filter struct {
	Deployment     string
	InstanceGroups []string
}

func (f Filter) ShouldProcessDeployment(deployment string) bool {
	return len(f.Deployment) == 0 || f.Deployment == deployment
}

func (f Filter) ShouldProcessInstanceGroup(instanceGroup string) bool {
	return len(f.InstanceGroups) == 0 || slices.Contains(f.InstanceGroups, instanceGroup)
}
