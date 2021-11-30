### 配置多账号git

1. `git config --global user.name` 和 `git config --global user.email` 不能存在值，若存在

`git config --global --unset user.name` 删除

2. 配置 `sshKey` 

```
ssh-keygen -t rsa -C 'email';
// Enter file in which to save the key (/c/Users/Excailbur/.ssh/id_rsa):这里的名字根据需要自定义；如：
/c/Users/Excailbur/.ssh/git_id_rsa
// 后面的提示可以直接 enter 过
```

如果有多个账号重复 2 步骤即可

3. 在 `~/.ssh/` 目录下创建一个 `config` 文件，配置如下：

```
Host github.com
HostName github.com
User shea  # 这个名字随便命
IdentityFile ~/.ssh/github_id_rsa

Host gitlab.com
HostName gitlab.com
User shea
IdentityFile ~/.ssh/gitlab_id_rsa
```

4. 配置完将 `xx_id_rsa.pub` 公钥复制到对应的网站即可，随后测试 `ssh -T git@github.com`

#### 让git 的 ssh 协议走socks5代理

> 首先下载 connect ，可以使用 scoop 下载
>
> Host github.com
> 	HostName github.com
> 	User shea  # 这个名字随便命
> 	IdentityFile ~/.ssh/github_id_rsa
>
> ProxyCommand connect -S 127.0.0.1:51837 %h %p