# rock-chameleon-go
磐石内置蜜罐系统

## mysql 
轻量内置mysql服务器, 支持数据伪造和驱动链接 支持audit event如下
- honey_mysql_auth
- honey_mysql_conn
- honey_mysql_query
- honey_mysql_raw
```lua
local mysql = chameleon.mysql
local auth = mysql.auth("root" , "123456") --设置认证


local users  = mysql.new_table("users" ,
    { name = "name", type = "text", null = false },
    { name = "mail", type = "text", null = false },
    { name = "age",  type = "int",  null = false },
    { name = "pass", type = "text", null = false })

users.insert("edunx" , "123@qq.com" , 18 , "123654")
users.insert("edunx" , "123@qq.com" , 18 , "123654")
users.insert("edunx" , "123@qq.com" , 18 , "123654")
users.insert("suncle" , "1@qq.com" , 20 , "123654")

local infos  = mysql.new_table("infos" , --新建表
    { name = "name", type = "text", null = false },
    { name = "mail", type = "text", null = false },
    { name = "age",  type = "int",  null = false },
    { name = "pass", type = "text", null = false })

--插入数据
infos.insert("edunx" , "123@qq.com" , 18 , "123654")
infos.insert("edunx" , "123@qq.com" , 18 , "123654")
infos.insert("edunx" , "123@qq.com" , 18 , "123654")
infos.insert("suncle" , "456@qq.com" , 20 , "123654")

local rock_db = mysql.new_db("rock_db") --新建database
rock_db.add_table( users , infos) --添加表结构

proc.start(mysql.new{
    name = "mysql", 
    bind = "0.0.0.0:3308",
    database = rock_db,
    auth = auth,
}
```

## ssh
轻量内置ssh服务器
```lua
    
local ssh_s = chameleon.ssh{
    name = "ssh_honey_pot",
    bind = "0.0.0.0:2222",
    prompt = "~$",
}
ssh_s.auth_root = "123456" --设置root账户密码
ssh_s.auth_app = "123456" --新增app账户密码
proc.start(ssh_s)

```

## proxy
代理型蜜罐
```lua
local proxy = chameleon.proxy
local proxy_mysql = proxy{
    name = "proxy_mysql",
    bind = "0.0.0.0:3310",     --对外端口
    remote = "127.0.0.1:3308", --后端地址
    protocol = "tcp",
}

proc.start(proxy_mysql)

local proxy_ssh = proxy{
    name = "proxy_ssh",
    bind = "0.0.0.0:1222",      --对外端口
    remote = "172.31.61.91:22", --后端地址
    protocol = "tcp",
}

proc.start(proxy_ssh)


```

## ftp
```lua
--todo
```

## smb
```lua
--todo
```

## rdp
```lua
--todo
```

## redis
```lua
--todo
```
