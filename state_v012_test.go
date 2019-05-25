package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedStateV012 = StateV012{
	Resources: []ResourceV012{
		{
			Name: "group_1",
			Type: "ansible_group",
			Instances: []InstanceV012{
				{
					Attributes: map[string]interface{}{
						"id":                   "group_1",
						"inventory_group_name": "group_1",
						"children":             []interface{}{"group_2"},
						"vars": map[string]interface{}{
							"foo": "bar",
						},
					},
				},
			},
		},
		{
			Name: "group_2",
			Type: "ansible_group",
			Instances: []InstanceV012{
				{
					Attributes: map[string]interface{}{
						"id":                   "group_2",
						"inventory_group_name": "group_2",
						"children":             nil,
						"vars":                 nil,
					},
				},
			},
		},
		{
			Name: "other_groups",
			Type: "ansible_group",
			Instances: []InstanceV012{
				{
					Attributes: map[string]interface{}{
						"id":                   "some_group_0",
						"inventory_group_name": "some_group_0",
						"children":             nil,
						"vars":                 nil,
					},
				},
				{
					Attributes: map[string]interface{}{
						"id":                   "some_group_1",
						"inventory_group_name": "some_group_1",
						"children":             nil,
						"vars":                 nil,
					},
				},
			},
		},
		{
			Name: "host_1",
			Type: "ansible_host",
			Instances: []InstanceV012{
				{
					Attributes: map[string]interface{}{
						"id":                 "host_1",
						"inventory_hostname": "host_1",
						"groups":             []interface{}{"group_1"},
						"vars": map[string]interface{}{
							"ansible_host": "1.2.3.4",
							"ansible_user": "ubuntu",
							"test":         "host_1",
						},
					},
				},
			},
		},
		{
			Name: "host_2",
			Type: "ansible_host",
			Instances: []InstanceV012{
				{
					Attributes: map[string]interface{}{
						"id":                 "host_2",
						"inventory_hostname": "host_2",
						"groups":             []interface{}{"group_1"},
						"vars": map[string]interface{}{
							"ansible_host": "1.2.3.5",
							"ansible_user": "ubuntu",
							"test":         "host_2",
						},
					},
				},
			},
		},
		{
			Name: "host_3",
			Type: "ansible_host",
			Instances: []InstanceV012{
				{
					Attributes: map[string]interface{}{
						"id":                 "host_3",
						"inventory_hostname": "host_3",
						"groups":             []interface{}{"group_3"},
						"vars": map[string]interface{}{
							"ansible_host": "1.2.3.6",
							"ansible_user": "ubuntu",
						},
					},
				},
			},
		},
		{
			Name: "host_4",
			Type: "ansible_host",
			Instances: []InstanceV012{
				{
					Attributes: map[string]interface{}{
						"id":                 "host_4",
						"inventory_hostname": "host_4",
						"groups":             nil,
						"vars": map[string]interface{}{
							"ansible_host": "1.2.3.7",
							"ansible_user": "ubuntu",
						},
					},
				},
			},
		},
		{
			Name: "host_5",
			Type: "ansible_host",
			Instances: []InstanceV012{
				{
					Attributes: map[string]interface{}{
						"id":                 "host_5",
						"inventory_hostname": "host_5",
						"groups":             nil,
						"vars": map[string]interface{}{
							"ansible_host": "1.2.3.8",
							"ansible_user": "ubuntu",
						},
					},
				},
			},
		},
		{
			Name: "other_hosts",
			Type: "ansible_host",
			Instances: []InstanceV012{
				{
					Attributes: map[string]interface{}{
						"id":                 "some_host_0",
						"inventory_hostname": "some_host_0",
						"groups":             []interface{}{"some_group_0"},
						"vars": map[string]interface{}{
							"ansible_host": "1.2.4.0",
							"ansible_user": "ubuntu",
						},
					},
				},
				{
					Attributes: map[string]interface{}{
						"id":                 "some_host_1",
						"inventory_hostname": "some_host_1",
						"groups":             []interface{}{"some_group_1"},
						"vars": map[string]interface{}{
							"ansible_host": "1.2.4.1",
							"ansible_user": "ubuntu",
						},
					},
				},
			},
		},
	},
}

var expectedInventoryV012 = map[string]interface{}{
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

func TestStateV012_basic(t *testing.T) {
	actual, err := getState("fixtures/v012")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedStateV012, actual)

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

	assert.Equal(t, expectedInventoryV012, actualInventory)
}
