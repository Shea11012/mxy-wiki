---
title: "ping not permitted"
created: 2022-08-31
updated: 2026-07-22
status: needs-review
type: concept
tags:
  - linux
  - linux/network
confidence: INFERRED
---


## ping socket operation not permitted

```shell
sudo setcap 'cap_net_admin,cap_net_raw=eip' $(which ping)
```