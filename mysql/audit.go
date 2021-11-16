package mysql

import (
	"github.com/rock-go/rock-chameleon-go/mysql/auth"
	"github.com/rock-go/rock-chameleon-go/mysql/sql"
	"github.com/rock-go/rock/audit"
	"time"
)

type Audit struct{
	CodeVM 	func() string
}

func (a *Audit) Authentication(user, addr string, err error) {
	ev := audit.NewEvent("chameleon" ,
		audit.Subject("honey mysql auth") ,
		audit.From(a.CodeVM()),
		audit.Remote(addr) ,
		audit.User(user))

	if err == nil {
		ev.Set(audit.Msg("honey mysql auth success"))
	} else {
		ev.Set(audit.Msg("honey mysql auth error"), audit.E(err))
	}

	ev.Put()
}

func (a *Audit) Authorization(ctx *sql.Context, p auth.Permission, err error) {
	//	fmt.Printf("authO: %s %s %v\n" ,  ctx.Session , ctx.Client().Address, p)

}
func (a *Audit) Query(ctx *sql.Context, d time.Duration, err error) {
	//  "user":          ctx.Client().User,
	//	"query":         ctx.Query(),
	//	"address":       ctx.Client().Address,
	//	"connection_id": ctx.Session.ID(),
	//	"pid":           ctx.Pid(),
	//	"success":       true,

	//fmt.Printf("Query: %s %s %s %v %s\n" , d , ctx.Session , ctx.Client().Address, ctx.Query())
}
