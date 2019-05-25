package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedStateV011 = StateV011{
	Modules: []ModuleV011{
		ModuleV011{
			Resources: map[string]ResourceV011{
				"ansible_host.host_1": ResourceV011{
					Type: "ansible_host",
					Primary: PrimaryV011{
						ID: "host_1",
						Attributes: map[string]string{
							"id":                 "host_1",
							"inventory_hostname": "host_1",
							"groups.#":           "1",
							"groups.0":           "group_1",
							"vars.%":             "3",
							"vars.ansible_host":  "1.2.3.4",
							"vars.ansible_user":  "ubuntu",
							"vars.test":          "host_1",
						},
					},
				},
				"ansible_host.host_2": ResourceV011{
					Type: "ansible_host",
					Primary: PrimaryV011{
						ID: "host_2",
						Attributes: map[string]string{
							"id":                 "host_2",
							"inventory_hostname": "host_2",
							"groups.#":           "1",
							"groups.0":           "group_1",
							"vars.%":             "3",
							"vars.ansible_host":  "1.2.3.5",
							"vars.ansible_user":  "ubuntu",
							"vars.test":          "host_2",
						},
					},
				},
				"ansible_host.host_3": ResourceV011{
					Type: "ansible_host",
					Primary: PrimaryV011{
						ID: "host_3",
						Attributes: map[string]string{
							"id":                 "host_3",
							"inventory_hostname": "host_3",
							"groups.#":           "1",
							"groups.0":           "group_3",
							"vars.%":             "2",
							"vars.ansible_host":  "1.2.3.6",
							"vars.ansible_user":  "ubuntu",
						},
					},
				},
				"ansible_host.host_4": ResourceV011{
					Type: "ansible_host",
					Primary: PrimaryV011{
						ID: "host_4",
						Attributes: map[string]string{
							"id":                 "host_4",
							"inventory_hostname": "host_4",
							"vars.%":             "2",
							"vars.ansible_host":  "1.2.3.7",
							"vars.ansible_user":  "ubuntu",
						},
					},
				},
				"ansible_host.other_hosts.0": ResourceV011{
					Type: "ansible_host",
					Primary: PrimaryV011{
						ID: "some_host_0",
						Attributes: map[string]string{
							"id":                 "some_host_0",
							"inventory_hostname": "some_host_0",
							"groups.#":           "1",
							"groups.0":           "some_group_0",
							"vars.%":             "2",
							"vars.ansible_host":  "1.2.4.0",
							"vars.ansible_user":  "ubuntu",
						},
					},
				},
				"ansible_host.other_hosts.1": ResourceV011{
					Type: "ansible_host",
					Primary: PrimaryV011{
						ID: "some_host_1",
						Attributes: map[string]string{
							"id":                 "some_host_1",
							"inventory_hostname": "some_host_1",
							"groups.#":           "1",
							"groups.0":           "some_group_1",
							"vars.%":             "2",
							"vars.ansible_host":  "1.2.4.1",
							"vars.ansible_user":  "ubuntu",
						},
					},
				},
				"ansible_group.group_1": ResourceV011{
					Type: "ansible_group",
					Primary: PrimaryV011{
						ID: "group_1",
						Attributes: map[string]string{
							"id":                   "group_1",
							"inventory_group_name": "group_1",
							"children.#":           "1",
							"children.0":           "group_2",
							"vars.%":               "1",
							"vars.foo":             "bar",
						},
					},
				},
				"ansible_group.group_2": ResourceV011{
					Type: "ansible_group",
					Primary: PrimaryV011{
						ID: "group_2",
						Attributes: map[string]string{
							"id":                   "group_2",
							"inventory_group_name": "group_2",
						},
					},
				},
				"ansible_group.other_groups.0": ResourceV011{
					Type: "ansible_group",
					Primary: PrimaryV011{
						ID: "some_group_0",
						Attributes: map[string]string{
							"id":                   "some_group_0",
							"inventory_group_name": "some_group_0",
						},
					},
				},
				"ansible_group.other_groups.1": ResourceV011{
					Type: "ansible_group",
					Primary: PrimaryV011{
						ID: "some_group_1",
						Attributes: map[string]string{
							"id":                   "some_group_1",
							"inventory_group_name": "some_group_1",
						},
					},
				},
			},
		},
		{
			Resources: map[string]ResourceV011{
				"ansible_host.host_5": ResourceV011{
					Type: "ansible_host",
					Primary: PrimaryV011{
						ID: "host_5",
						Attributes: map[string]string{
							"id":                 "host_5",
							"inventory_hostname": "host_5",
							"vars.%":             "2",
							"vars.ansible_host":  "1.2.3.8",
							"vars.ansible_user":  "ubuntu",
						},
					},
				},
			},
		},
	},
}

var expectedInventoryV011 = map[string]interface{}{
	"all": map[string]interface{}{
		"hosts": []string{"host_1", "host_2", "host_3", "host_4", "host_5", "some_host_0", "some_host_1"},
		"vars":  map[string]interface{}{},
	},
	"ungrouped": map[string]interface{}{
		"hosts": []string{"host_4", "host_5"},
		"vars":  map[string]interface{}{},
	},
	"group_1": map[string]interface{}{
		"hosts":    []string{"host_1", "host_2"},
		"children": []string{"group_2"},
		"vars": map[string]interface{}{
			"foo": "bar",
		},
	},
	"group_2": map[string]interface{}{
		"vars": map[string]interface{}{},
	},
	"group_3": map[string]interface{}{
		"hosts": []string{"host_3"},
		"vars":  map[string]interface{}{},
	},
	"some_group_0": map[string]interface{}{
		"hosts": []string{"some_host_0"},
		"vars":  map[string]interface{}{},
	},
	"some_group_1": map[string]interface{}{
		"hosts": []string{"some_host_1"},
		"vars":  map[string]interface{}{},
	},
	"_meta": map[string]interface{}{
		"hostvars": map[string]interface{}{
			"host_1": map[string]interface{}{
				"ansible_host": "1.2.3.4",
				"ansible_user": "ubuntu",
				"test":         "host_1",
			},
			"host_2": map[string]interface{}{
				"ansible_host": "1.2.3.5",
				"ansible_user": "ubuntu",
				"test":         "host_2",
			},
			"host_3": map[string]interface{}{
				"ansible_host": "1.2.3.6",
				"ansible_user": "ubuntu",
			},
			"host_4": map[string]interface{}{
				"ansible_host": "1.2.3.7",
				"ansible_user": "ubuntu",
			},
			"host_5": map[string]interface{}{
				"ansible_host": "1.2.3.8",
				"ansible_user": "ubuntu",
			},
			"some_host_0": map[string]interface{}{
				"ansible_host": "1.2.4.0",
				"ansible_user": "ubuntu",
			},
			"some_host_1": map[string]interface{}{
				"ansible_host": "1.2.4.1",
				"ansible_user": "ubuntu",
			},
		},
	},
}

func TestStateV011_basic(t *testing.T) {
	actual, err := getState("fixtures/v011")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedStateV011, actual)

	expectedGroups := []string{"group_1", "group_2", "some_group_0", "some_group_1"}
	actualGroups, err := actual.GetGroups()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedGroups, actualGroups)

	expectedHosts := []string{"host_1", "host_2"}
	actualHosts, err := actual.GetHostsForGroup(expectedGroups[0])
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedHosts, actualHosts)

	expectedVars := map[string]interface{}{
		"ansible_host": "1.2.3.4",
		"ansible_user": "ubuntu",
		"test":         "host_1",
	}

	actualVars, err := actual.GetVarsForHost("host_1")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedVars, actualVars)

	actualInventory, err := BuildInventory(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedInventoryV011, actualInventory)
}
