resource "ansible_group" "group_1" {
  inventory_group_name = "group_1"
  children             = ["group_2"]

  vars = {
    foo = "bar"
  }
}

resource "ansible_group" "group_2" {
  inventory_group_name = "group_2"
}

resource "ansible_group" "other_groups" {
  count                = 2
  inventory_group_name = "some_group_${count.index}"
}

resource "ansible_host" "host_1" {
  inventory_hostname = "host_1"
  groups             = ["group_1"]

  vars = {
    ansible_user = "ubuntu"
    ansible_host = "1.2.3.4"
    test         = "host_1"
  }
}

resource "ansible_host" "host_2" {
  inventory_hostname = "host_2"
  groups             = ["group_1"]

  vars = {
    ansible_user = "ubuntu"
    ansible_host = "1.2.3.5"
    test         = "host_2"
  }
}

resource "ansible_host" "host_3" {
  inventory_hostname = "host_3"
  groups             = ["group_3"]

  vars = {
    ansible_user = "ubuntu"
    ansible_host = "1.2.3.6"
  }
}

resource "ansible_host" "host_4" {
  inventory_hostname = "host_4"

  vars = {
    ansible_user = "ubuntu"
    ansible_host = "1.2.3.7"
  }
}

resource "ansible_host" "other_hosts" {
  count              = 2
  inventory_hostname = "some_host_${count.index}"
  groups             = ["some_group_${count.index}"]

  vars = {
    ansible_user = "ubuntu"
    ansible_host = "1.2.4.${count.index}"
  }
}

module "more_hosts" {
  source = "./module"
}
