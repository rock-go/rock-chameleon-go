package chameleon

import (
	"github.com/rock-go/rock-chameleon-go/mysql"
	"github.com/rock-go/rock-chameleon-go/proxy"
	"github.com/rock-go/rock-chameleon-go/ssh"
	"github.com/rock-go/rock-chameleon-go/stream"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xbase"
)

func LuaInjectApi(env *xbase.EnvT) {
	uv := lua.NewUserKV()
	proxy.Inject(env, uv)
	stream.Inject(env, uv)
	mysql.Inject(env, uv)
	ssh.Inject(env, uv)
	env.Global("chameleon", uv)
}
