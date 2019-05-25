resource "ansible_host" "host_5" {
  inventory_hostname = "host_5"

  vars = {
    ansible_user = "ubuntu"
    ansible_host = "1.2.3.8"
  }
}
