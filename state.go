package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// The following structs are for Terraform State.
type State struct {
	Modules []Module `json:"modules"`
}

func (r State) GetGroups() ([]string, error) {
	var groups []string

	for _, m := range r.Modules {
		for _, resource := range m.Resources {
			if resource.Type == "ansible_group" {
				groups = append(groups, resource.Primary.ID)
			}
		}
	}

	return groups, nil
}

func (r State) GetGroup(group string) (*Resource, error) {
	for _, m := range r.Modules {
		for _, resource := range m.Resources {
			if resource.Type == "ansible_group" {
				if resource.Primary.ID == group {
					return &resource, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("Unable to find group %s", group)
}
func (r State) GetChildrenForGroup(group string) ([]string, error) {
	var children []string

	resource, err := r.GetGroup(group)
	if err != nil {
		return nil, err
	}

	for attrName, attr := range resource.Primary.Attributes {
		if strings.HasPrefix(attrName, "children.") {
			if attrName == "children.#" {
				continue
			}
			children = append(children, attr)
		}
	}

	return children, nil
}

func (r State) GetVarsForGroup(group string) (map[string]interface{}, error) {
	vars := make(map[string]interface{})

	resource, err := r.GetGroup(group)
	if err != nil {
		return nil, err
	}

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

func (r State) GetHostsForGroup(group string) ([]string, error) {
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

	return hosts, nil
}

func (r State) GetVarsForHost(host string) (map[string]interface{}, error) {
	vars := make(map[string]interface{})

	for _, m := range r.Modules {
		for _, resource := range m.Resources {
			if resource.Type == "ansible_host" {
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
			}
		}
	}

	return vars, nil
}

func (r State) BuildInventory() (map[string]interface{}, error) {
	inv := make(map[string]interface{})
	meta := make(map[string]interface{})
	hostvars := make(map[string]interface{})

	groups, err := r.GetGroups()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		g := make(map[string]interface{})
		hosts, err := r.GetHostsForGroup(group)
		if err != nil {
			return nil, err
		}

		children, err := r.GetChildrenForGroup(group)
		if err != nil {
			return nil, err
		}

		vars, err := r.GetVarsForGroup(group)
		if err != nil {
			return nil, err
		}

		if len(hosts) > 0 {
			g["hosts"] = hosts
		}

		if len(children) > 0 {
			g["children"] = children
		}

		g["vars"] = vars
		inv[group] = g

		for _, host := range hosts {
			vars, err := r.GetVarsForHost(host)
			if err != nil {
				return nil, err
			}

			hostvars[host] = vars
		}

	}

	meta["hostvars"] = hostvars
	inv["_meta"] = meta
	return inv, nil
}

func (r State) ToJSON() (string, error) {
	var s string

	inv, err := r.BuildInventory()
	if err != nil {
		return s, err
	}

	b, err := json.Marshal(inv)
	if err != nil {
		return s, err
	}

	s = string(b)

	return s, nil
}

type Module struct {
	Resources map[string]Resource `json:"resources"`
}

type Resource struct {
	Type    string  `json:"type"`
	Primary Primary `json:"primary"`
}

type Primary struct {
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}

func getState(path string) (*State, error) {
	cmd := exec.Command("terraform", "state", "pull")
	cmd.Dir = path
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Error running `terraform state pull` in directory %s, %s\n", path, err)
	}

	b, err := ioutil.ReadAll(&out)
	if err != nil {
		return nil, fmt.Errorf("Error reading output of `terraform state pull`: %s\n", err)
	}

	var s State
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling state: %s\n", err)
	}

	return &s, nil
}
