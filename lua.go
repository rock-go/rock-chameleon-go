package chameleon

import (
	"github.com/rock-go/rock-chameleon-go/mysql"
	"github.com/rock-go/rock-chameleon-go/proxy"
	"github.com/rock-go/rock-chameleon-go/ssh"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
)

func LuaInjectApi(env xcall.Env) {
	uv := lua.NewUserKV()
	proxy.Inject(uv)
	mysql.Inject(uv)
	ssh.Inject(uv)
	env.SetGlobal("chameleon", uv)
}
