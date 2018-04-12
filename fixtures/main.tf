resource "ansible_group" "group_1" {
  inventory_group_name = "group_1"
  children             = ["group_2"]

  vars {
    foo = "bar"
  }
}

resource "ansible_group" "group_2" {
  inventory_group_name = "group_2"
}

resource "ansible_host" "host_1" {
  inventory_hostname = "host_1"
  groups             = ["group_1"]

  vars {
    ansible_user = "ubuntu"
    ansible_host = "1.2.3.4"
    test         = "host_1"
  }
}

resource "ansible_host" "host_2" {
  inventory_hostname = "host_2"
  groups             = ["group_1"]

  vars {
    ansible_user = "ubuntu"
    ansible_host = "1.2.3.5"
    test         = "host_2"
  }
}
