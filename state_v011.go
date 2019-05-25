package main

import (
	"fmt"
	"sort"
	"strings"
)

// The following structs are for Terraform State
// from version v0.11 and prior.
type StateV011 struct {
	Modules []ModuleV011 `json:"modules"`
}

// GetGroups will return all ansible_group resources.
func (r StateV011) GetGroups() ([]string, error) {
	var groups []string

	for _, m := range r.Modules {
		for _, resource := range m.Resources {
			if resource.Type == "ansible_group" {
				groups = append(groups, resource.Primary.ID)
			}
		}
	}

	sort.Strings(groups)
	return groups, nil
}

// GetHosts will return all ansible_group resources.
func (r StateV011) GetHosts() ([]string, error) {
	var hosts []string

	for _, m := range r.Modules {
		for _, resource := range m.Resources {
			if resource.Type == "ansible_host" {
				hosts = append(hosts, resource.Primary.ID)
			}
		}
	}

	sort.Strings(hosts)
	return hosts, nil
}

// GetGroup will find and return a specific ansible_group resource.
func (r StateV011) GetGroup(group string) (interface{}, error) {
	for _, m := range r.Modules {
		for _, resource := range m.Resources {
			if resource.Type == "ansible_group" {
				if resource.Primary.ID == group {
					return resource, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("Unable to find group %s", group)
}

// GetChildrenForGroup will return the "children" members of an
// ansible_group resource.
func (r StateV011) GetChildrenForGroup(group string) ([]string, error) {
	var children []string
	var resource ResourceV011

	v, err := r.GetGroup(group)
	if err != nil {
		return nil, err
	}

	resource = v.(ResourceV011)

	for attrName, attr := range resource.Primary.Attributes {
		if strings.HasPrefix(attrName, "children.") {
			if attrName == "children.#" {
				continue
			}
			children = append(children, attr)
		}
	}

	sort.Strings(children)
	return children, nil
}

// GetVarsForGroup will return the variables defined in an ansible_group
// resource.
func (r StateV011) GetVarsForGroup(group string) (map[string]interface{}, error) {
	var resource ResourceV011
	vars := make(map[string]interface{})

	v, err := r.GetGroup(group)
	if err != nil {
		return nil, err
	}

	resource = v.(ResourceV011)

	for attrName, attr := range resource.Primary.Attributes {
		if strings.HasPrefix(attrName, "vars.") {
			if attrName == "vars.%" {
				continue
			}

			pieces := strings.SplitN(attrName, ".", 2)
			if len(pieces) == 2 {
				vars[pieces[1]] = attr
			}
		}
	}

	return vars, nil
}

// GetHostsForGroup will return the hosts that belong to a defined group.
func (r StateV011) GetHostsForGroup(group string) ([]string, error) {
	var hosts []string

	for _, m := range r.Modules {
		for _, resource := range m.Resources {
			if resource.Type == "ansible_host" {
				for attrName, attr := range resource.Primary.Attributes {
					if strings.HasPrefix(attrName, "groups.") {
						if group == attr {
							hosts = append(hosts, resource.Primary.ID)
						}
					}
				}
			}
		}
	}

	sort.Strings(hosts)
	return hosts, nil
}

// GetHost will return a specific ansible_host.
func (r StateV011) GetHost(host string) (interface{}, error) {
	for _, m := range r.Modules {
		for _, resource := range m.Resources {
			if resource.Type == "ansible_host" {
				if resource.Primary.ID == host {
					return resource, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("Unable to find host %s", host)
}

// GetGroupsForHost will return the groups defined in an ansible_host resource.
func (r StateV011) GetGroupsForHost(host string) ([]string, error) {
	var resource ResourceV011
	groups := []string{}

	v, err := r.GetHost(host)
	if err != nil {
		return nil, err
	}

	resource = v.(ResourceV011)

	for attrName, attr := range resource.Primary.Attributes {
		if strings.HasPrefix(attrName, "groups.") {
			if attrName == "groups.#" {
				continue
			}

			pieces := strings.SplitN(attrName, ".", 2)
			if len(pieces) == 2 {
				groups = append(groups, attr)
			}
		}
	}

	return groups, nil
}

// GetVarsForHost will return the variables defined in an ansible_host resource.
func (r StateV011) GetVarsForHost(host string) (map[string]interface{}, error) {
	var resource ResourceV011
	vars := make(map[string]interface{})

	v, err := r.GetHost(host)
	if err != nil {
		return nil, err
	}

	resource = v.(ResourceV011)

	for attrName, attr := range resource.Primary.Attributes {
		if strings.HasPrefix(attrName, "vars.") {
			if attrName == "vars.%" {
				continue
			}

			pieces := strings.SplitN(attrName, ".", 2)
			if len(pieces) == 2 {
				vars[pieces[1]] = attr
			}
		}
	}

	return vars, nil
}

type ModuleV011 struct {
	Resources map[string]ResourceV011 `json:"resources"`
}

type ResourceV011 struct {
	Type    string      `json:"type"`
	Primary PrimaryV011 `json:"primary"`
}

type PrimaryV011 struct {
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}
