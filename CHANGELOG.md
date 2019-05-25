## 0.3.0

IMPROVEMENTS

* Initial support for Terraform v0.12 [GH-6](https://github.com/jtopjian/ansible-terraform-inventory/pull/6)
* Support for "implicit" groups has been added. When an ansible_host declares that it is part of a group, but no ansible_group resource exists, the group will automatically be created. [GH-5](https://github.com/jtopjian/ansible-terraform-inventory/pull/5)
* Support for an "all" group has been added [GH-5](https://github.com/jtopjian/ansible-terraform-inventory/pull/5)
* Support for an "ungrouped" group has been added [GH-5](https://github.com/jtopjian/ansible-terraform-inventory/pull/5)
