package main

import (
	"encoding/json"
	"sort"
)

// Interface State represents the methods a state struct has to implement to
// parse a Terraform state.
type State interface {
	GetGroups() ([]string, error)
	GetGroup(group string) (interface{}, error)
	GetGroupsForHost(host string) ([]string, error)

	GetChildrenForGroup(group string) ([]string, error)

	GetVarsForGroup(group string) (map[string]interface{}, error)
	GetVarsForHost(host string) (map[string]interface{}, error)

	GetHosts() ([]string, error)
	GetHost(host string) (interface{}, error)
	GetHostsForGroup(group string) ([]string, error)
}

func BuildInventory(state State) (map[string]interface{}, error) {
	inv := make(map[string]interface{})
	meta := make(map[string]interface{})
	hostvars := make(map[string]interface{})
	allHosts := []string{}

	// Get all ansible_group resources.
	groups, err := state.GetGroups()
	if err != nil {
		return nil, err
	}

	// For each ansible_group defined...
	for _, group := range groups {
		g := make(map[string]interface{})
		hosts, err := state.GetHostsForGroup(group)
		if err != nil {
			return nil, err
		}

		// Get the children of the group.
		children, err := state.GetChildrenForGroup(group)
		if err != nil {
			return nil, err
		}

		// Get any variables for the group.
		vars, err := state.GetVarsForGroup(group)
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
	hosts, err := state.GetHosts()
	if err != nil {
		return nil, err
	}

	// For each host...
	for _, host := range hosts {
		// Add the host to the set of all hosts.
		allHosts = append(allHosts, host)

		// Get any variable defined and set it in the inventory.
		vars, err := state.GetVarsForHost(host)
		if err != nil {
			return nil, err
		}

		hostvars[host] = vars

		// Find all groups that the host is a part of.
		groups, err := state.GetGroupsForHost(host)
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

func ToJSON(state State) (string, error) {
	var s string

	inv, err := BuildInventory(state)
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
