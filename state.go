package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"
)

// The following structs are for Terraform State.
type State struct {
	Modules []Module `json:"modules"`
}

// GetGroups will return all ansible_group resources.
func (r State) GetGroups() ([]string, error) {
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
func (r State) GetHosts() ([]string, error) {
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

// GetChildrenForGroup will return the "children" members of an
// ansible_group resource.
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

	sort.Strings(children)
	return children, nil
}

// GetVarsForGroup will return the variables defined in an ansible_group
// resource.
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

// GetHostsForGroup will return the hosts in a specific ansible_group
// resource.
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

	sort.Strings(hosts)
	return hosts, nil
}

// GetHost will return a specific ansible_host.
func (r State) GetHost(host string) (*Resource, error) {
	for _, m := range r.Modules {
		for _, resource := range m.Resources {
			if resource.Type == "ansible_host" {
				if resource.Primary.ID == host {
					return &resource, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("Unable to find host %s", host)
}

// GetGroupsForHost will return the groups defined in an ansible_host resource.
func (r State) GetGroupsForHost(host string) ([]string, error) {
	groups := []string{}

	resource, err := r.GetHost(host)
	if err != nil {
		return nil, err
	}

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
func (r State) GetVarsForHost(host string) (map[string]interface{}, error) {
	vars := make(map[string]interface{})

	resource, err := r.GetHost(host)
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

func (r State) BuildInventory() (map[string]interface{}, error) {
	inv := make(map[string]interface{})
	meta := make(map[string]interface{})
	hostvars := make(map[string]interface{})
	allHosts := []string{}

	// Get all ansible_group resources.
	groups, err := r.GetGroups()
	if err != nil {
		return nil, err
	}

	// For each ansible_group defined...
	for _, group := range groups {
		g := make(map[string]interface{})
		hosts, err := r.GetHostsForGroup(group)
		if err != nil {
			return nil, err
		}

		// Get the children of the group.
		children, err := r.GetChildrenForGroup(group)
		if err != nil {
			return nil, err
		}

		// Get any variables for the group.
		vars, err := r.GetVarsForGroup(group)
		if err != nil {
			return nil, err
		}

		// Set the hosts.
		if len(hosts) > 0 {
			g["hosts"] = hosts
		}

		// Set the children.
		if len(children) > 0 {
			g["children"] = children
		}

		// Set the variables.
		g["vars"] = vars
		inv[group] = g
	}

	// Now that we've accounted for all explicitly defined
	// groups, let's find any groups which were implicitly
	// defined. These are ansible_host resources with group
	// memberships of groups that have no explicit
	// ansible_group resource.
	//
	// In addition, create an "ungrouped" group which will
	// contain hosts that have no group membership.
	var ungrouped []string

	// Get all ansible_host resources defined.
	hosts, err := r.GetHosts()
	if err != nil {
		return nil, err
	}

	// For each host...
	for _, host := range hosts {
		// Add the host to the set of all hosts.
		allHosts = append(allHosts, host)

		// Get any variable defined and set it in the inventory.
		vars, err := r.GetVarsForHost(host)
		if err != nil {
			return nil, err
		}

		hostvars[host] = vars

		// Find all groups that the host is a part of.
		groups, err := r.GetGroupsForHost(host)
		if err != nil {
			return nil, err
		}

		// If no groups were defined, add the host to the "ungrouped" group.
		if len(groups) == 0 {
			ungrouped = append(ungrouped, host)
		}

		// For each group defined in the host...
		for _, group := range groups {
			// Check and see if this group has already been accounted for.
			// If it has, check for the host membership.
			if v, ok := inv[group]; ok {
				groupInventory := v.(map[string]interface{})
				if hostInventory, ok := groupInventory["hosts"].([]string); ok {
					var found bool
					for _, h := range hostInventory {
						if h == host {
							found = true
						}
					}

					if !found {
						hostInventory = append(hostInventory, host)
					}
				}
			} else {
				// if the group wasn't already accounted for, do it now.
				inv[group] = map[string]interface{}{
					"hosts": []string{host},
					"vars":  map[string]interface{}{},
				}
			}
		}
	}

	// If there are any "ungrouped" hosts, create an inventory entry
	// for "ungrouped".
	if len(ungrouped) > 0 {
		inv["ungrouped"] = map[string]interface{}{
			"hosts": ungrouped,
			"vars":  map[string]interface{}{},
		}
	}

	// Create an "all" group if one was not defined.
	if _, ok := inv["all"]; !ok {
		sort.Strings(allHosts)
		all := map[string]interface{}{
			"hosts": allHosts,
			"vars":  map[string]interface{}{},
		}

		inv["all"] = all
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

	// If there was no output, return nil and no error
	if string(b) == "" {
		return nil, nil
	}

	if string(b[0]) == "o" && string(b[1]) == ":" {
		b = append(b[:0], b[2:]...)
	}

	var s State
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling state: %s\n", err)
	}

	return &s, nil
}
