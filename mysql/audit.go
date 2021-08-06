package mysql

import (
	"github.com/rock-go/rock-chameleon-go/mysql/auth"
	"github.com/rock-go/rock-chameleon-go/mysql/sql"
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/audit/event"
	"time"
)

type Audit struct{}

func (a *Audit) Authentication(user, addr string, err error) {
	ev := event.New("honey_mysql_auth",
		event.Subject("honey mysql auth"),
		event.Addr(addr),
		event.User(user),
	)

	if err == nil {
		ev.Set(event.Infof("auth succeed"))
	} else {
		ev.Set(event.Infof("%s", err))
	}
	audit.Put(ev)
}

func (a *Audit) Authorization(ctx *sql.Context, p auth.Permission, err error) {
	//	fmt.Printf("authO: %s %s %v\n" ,  ctx.Session , ctx.Client().Address, p)

}
func (a *Audit) Query(ctx *sql.Context, d time.Duration, err error) {
	//"user":          ctx.Client().User,
	//	"query":         ctx.Query(),
	//	"address":       ctx.Client().Address,
	//	"connection_id": ctx.Session.ID(),
	//	"pid":           ctx.Pid(),
	//	"success":       true,

	//fmt.Printf("Query: %s %s %s %v %s\n" , d , ctx.Session , ctx.Client().Address, ctx.Query())
}
