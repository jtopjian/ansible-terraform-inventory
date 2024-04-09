ansible-terraform-inventory
===========================

A dynamic inventory script for Ansible and Terraform.

Quickstart
----------

To use this inventory script, you must first create Terraform resources
using the [terraform-provider-ansible](https://github.com/nbering/terraform-provider-ansible)
plugin:

```hcl
resource "ansible_host" "example" {
  inventory_hostname = "example.com"
  groups = ["web"]
  vars {
    ansible_user = "admin"
  }
}

resource "ansible_group" "web" {
  inventory_group_name = "web"
  children = ["foo", "bar", "baz"]
  vars {
    foo = "bar"
    bar = 2
  }
}
```

Next, use this script as your Ansible dynamic inventory script.

If your Ansible playbooks are in a different directory than your Terraform
resources, then set the `TF_STATE` environment variable to the location
of the Terraform directory.

If you want to use [terragrunt](https://terragrunt.gruntwork.io/) instead of
terraform, set the `TF_TERRAGRUNT` environment variable to any non-empty
value and set `TF_STATE` to the directory where the `terragrunt.hcl` file is
located.

Installation
------------

Download the latest [release](https://github.com/jtopjian/ansible-terraform-inventory/releases).

Building From Source
--------------------

```shell
$ go get github.com/jtopjian/ansible-terraform-inventory
$ go build -o $GOPATH/bin/terraform-inventory
$ ln -s $GOPATH/bin/terraform-inventory /path/to/ansible/hosts
```

Public GPG Key
--------------

The following is the public key which can be used to validate releases in this repository:

```-----BEGIN PGP PUBLIC KEY BLOCK-----

mQENBGMfm0gBCAC4HubI5t9Dt21PibwDlmW1yqD6R4ix1c/mX9k3DbD2rBD4y5zY
DcVxqEbdhGkWg00diPkoTJAvjYs6Byy0Y32PDYtubK4qck33qaWcM4dwh5X+Y2n0
1/YV/UBaKdrVgS8f8Haza6wLXTB/vkCb6/3Vs3NUhNonjmkF6W2RVW0WxfJ97oO3
PP6mXocT5ES6Xnbphfih2U4C92owep0+JRUiZs9OCfG9nroGVEAZuvELJUfZVnwh
TyucBBcDCXclpG+EcKnETKZTGFVd55iVKA7OEGQJxi/1Wi5vUS1nOpNcbK6j/n+6
Z8y3aokGf6iNoiAaFOcjmmxF+yIQELr5SuRdABEBAAG0HUpvZSBUb3BqaWFuIDxq
b2VAdG9wamlhbi5uZXQ+iQFOBBMBCgA4FiEEfRTM4xWef6ivXegDh+r7T9BsimcF
AmMfm0gCGwMFCwkIBwIGFQoJCAsCBBYCAwECHgECF4AACgkQh+r7T9BsimfS4QgA
osyLeHHSEn1FZcgGizi3cPwdcAfcbFP2bSdw7fQkE77LQYPLjHZ4zXSzyqgUU7Gi
qan5Q/tjJ72Z2vZqCfC99VIpE+XTH9AakKI+hsdkcKoUtkJ0LNja0IQ+qTkt9csH
rB3tb4LOp+InHRjnQQObeXAYVrvcZauRdmvesHSJwe78+OQrtFfmB4cgCU4dHl9n
HwOJbx2/LLqWH/juv5CYURj8/qfhEJtruHqcx+pkfMwfUnXuDTJhEr2wsJPWntOf
J0mKhTZLZ/S7fYuPejF/tszB2rEHXoNMJbOzLip35OnG+qdbkKaATE2o9buDtw3q
nPlbHSXqGKIow3+DsKITe7kBDQRjH5tIAQgA16dQGx+EhA8tmLLZ9Shur+DFjHKP
tbkkGKf+OWE89qoB7VcUgr8WP7XUXVjpUufWT/IYF6mq+M92WrTPrY38C5eObRkZ
C6Ua7W2aneWxSIl/yD/TsOPh75uBfQuM3f+PPFN5Rxo/1QAdePppZvTPJZLf8f52
h9N/yOi27SsgyEyBT7fTYfngVlevPKXaTTZdEgjH92QJ6oVm6k9jorROJ+eYNKjS
qxVX8ZWE5ZcsEa/tUz5sUv5tr0wAIkz+xbJxgS2P5Nl44bUtbjHX7L4reDimus/0
2+3z0cF/SD18V9obieJNDBiT1AKB6JeN0+W0tOMP5TsGfOd/vlgTY00BkwARAQAB
iQE2BBgBCgAgFiEEfRTM4xWef6ivXegDh+r7T9BsimcFAmMfm0gCGwwACgkQh+r7
T9BsimeDsAf+I8Wk10I5WKQnIXEdw/XXbZKjXxy8oobqCDG0Ayys2QIujH94HuD1
HvQjeWLQHNta+ySNNir3q81iGMPyiGuOdjHMRn/5Iykx+S5z3vZv8r7MYSAl5clB
2nQ/1ch9Zp37hxJ1fdz9aCF00qvQ4jb935lWxs8u7XpheXViz5Vcn1xwZgcJ4/fG
8AF5qpSUXvWG6x8xrG6FvAGK/Hpz4mwveQ4h/4QWfsrOyo5ssxnNOQDQ5YKy5hLE
5IrXRNCKTbROKeYSN2eaO3TcXHz9J+caWL3Jz/XyyIicOlf36TXhKucBipI00uYi
D/RQUP57qh/KrRtYQ+q62lcG0oJwUPWp4A==
=Tpdk
-----END PGP PUBLIC KEY BLOCK-----``
