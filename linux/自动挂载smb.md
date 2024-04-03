---
tags: 
date created: 2024-01-23 22:56
date modified: 2024-01-23 23:01
---
使用 nm-dispatcher 自动挂载 smb,在 `/etc/NetworkManager/dispatcher` 目录下创建文件 `30-samba.sh`
```bash
#!/bin/sh

# find connection uid with "nmcli c show"
WANTED_CON_UUID="a9fb169b-98a2-4f7f-bbfd-99455ba9a281"
USER="smbuser"
PASSWORD=123
TARGET="/mnt/share"
SMB_URL="//192.168.32.100/share"

#echo "current: $CONNECTION_UUID"
if [ "$CONNECTION_UUID" = "$WANTED_CON_UUID" ]; then
    # script paramter $1: network interface name
    # script paramter $2: dispatched event
    #echo "event $2"
    case "$2" in
        "up"|"vpn-up")
            #echo "enter up"
            mount --mkdir -t cifs $SMB_URL $TARGET -o username=$USER,password=$PASSWORD,iocharset=utf8,uid=1000,gid=1000
        ;;
       "pre-down"|"vpn-pre-down")
          #echo "enter down"
          umount $TARGET
        ;;
    esac
fi
```

在进行软链接
```bash
sudo ln -s /etc/NetworkManager/dispatcher.d/30-samba.sh /etc/NetworkManager/dispatcher.d/pre-down.d/30-samba.sh
```