---
created: 2022-08-31 21:09
updated: 2024-08-13
---

## ping socket operation not permitted

```shell
sudo setcap 'cap_net_admin,cap_net_raw=eip' $(which ping)
```