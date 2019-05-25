package main

import (
	"fmt"
	"sort"
)

// The following structs are for Terraform State
// for version v0.12.
type StateV012 struct {
	Resources []ResourceV012 `json:"resources"`
}

// GetGroups will return all ansible_group resources.
func (r StateV012) GetGroups() ([]string, error) {
	var groups []string

	for _, resource := range r.Resources {
		if resource.Type == "ansible_group" {
			for _, instance := range resource.Instances {
				if v, ok := instance.Attributes["inventory_group_name"].(string); ok {
					groups = append(groups, v)
				}
			}
		}
	}

	sort.Strings(groups)

	return groups, nil
}

// GetHosts will return all ansible_hosts resources.
func (r StateV012) GetHosts() ([]string, error) {
	var hosts []string

	for _, resource := range r.Resources {
		if resource.Type == "ansible_host" {
			for _, instance := range resource.Instances {
				if v, ok := instance.Attributes["inventory_hostname"].(string); ok {
					hosts = append(hosts, v)
				}
			}
		}
	}

	sort.Strings(hosts)

	return hosts, nil
}

// GetGroup will find and return a specific ansible_group resource.
func (r StateV012) GetGroup(group string) (interface{}, error) {
	for _, resource := range r.Resources {
		if resource.Type == "ansible_group" {
			for _, instance := range resource.Instances {
				if v, ok := instance.Attributes["inventory_group_name"].(string); ok {
					if v == group {
						return instance, nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("Unable to find group %s", group)
}

// GetChildrenForGroup will return the "children" members of an
// ansible_group resource.
func (r StateV012) GetChildrenForGroup(group string) ([]string, error) {
	var children []string
	var instance InstanceV012

	v, err := r.GetGroup(group)
	if err != nil {
		return nil, err
	}

	instance = v.(InstanceV012)

	if v, ok := instance.Attributes["children"].([]interface{}); ok {
		for _, c := range v {
			children = append(children, c.(string))
		}
	}

	sort.Strings(children)
	return children, nil
}

// GetVarsForGroup will return the variables defined in an ansible_group
// resource.
func (r StateV012) GetVarsForGroup(group string) (map[string]interface{}, error) {
	var instance InstanceV012
	vars := make(map[string]interface{})

	v, err := r.GetGroup(group)
	if err != nil {
		return nil, err
	}

	instance = v.(InstanceV012)

	if v, ok := instance.Attributes["vars"].(map[string]interface{}); ok {
		vars = v
	}

	return vars, nil
}

// GetHostsForGroup will return the hosts that belong to a defined group.
func (r StateV012) GetHostsForGroup(group string) ([]string, error) {
	var hosts []string

	for _, resource := range r.Resources {
		if resource.Type == "ansible_host" {
			for _, instance := range resource.Instances {
				hostname, ok := instance.Attributes["inventory_hostname"].(string)
				if !ok {
					continue
				}

				groups, ok := instance.Attributes["groups"].([]interface{})
				if !ok {
					continue
				}

				for _, g := range groups {
					if group == g.(string) {
						hosts = append(hosts, hostname)
					}
				}
			}
		}
	}

	sort.Strings(hosts)
	return hosts, nil
}

// GetHost will return a specific ansible_host.
func (r StateV012) GetHost(host string) (interface{}, error) {
	for _, resource := range r.Resources {
		if resource.Type == "ansible_host" {
			for _, instance := range resource.Instances {
				if v, ok := instance.Attributes["inventory_hostname"].(string); ok {
					if v == host {
						return instance, nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("Unable to find host %s", host)
}

// GetGroupsForHost will return the groups defined in an ansible_host resource.
func (r StateV012) GetGroupsForHost(host string) ([]string, error) {
	var instance InstanceV012
	groups := []string{}

	v, err := r.GetHost(host)
	if err != nil {
		return nil, err
	}

	instance = v.(InstanceV012)

	if v, ok := instance.Attributes["groups"].([]interface{}); ok {
		for _, group := range v {
			groups = append(groups, group.(string))
		}
	}

	return groups, nil
}

// GetVarsForHost will return the variables defined in an ansible_host resource.
func (r StateV012) GetVarsForHost(host string) (map[string]interface{}, error) {
	var instance InstanceV012
	vars := make(map[string]interface{})

	v, err := r.GetHost(host)
	if err != nil {
		return nil, err
	}

	instance = v.(InstanceV012)

	if v, ok := instance.Attributes["vars"].(map[string]interface{}); ok {
		vars = v
	}

	return vars, nil
}

type ResourceV012 struct {
	Type      string         `json:"type"`
	Name      string         `json:"name"`
	Instances []InstanceV012 `json:"instances"`
}

type InstanceV012 struct {
	Attributes map[string]interface{} `json:"attributes"`
}
