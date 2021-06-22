package casbin

import (
	"strings"
	"sync"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/vsatcloud/mars/models"
)

const modelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func Casbin(db models.Database) *casbin.SyncedEnforcer {
	once.Do(func() {
		a, _ := gormadapter.NewAdapter("postgresql", db.User+":"+db.Password+"@("+db.Host+")/"+db.Dbname, true)
		m := casbin.NewModel(modelText)
		e := casbin.NewSyncedEnforcer(m, a)
		e.AddFunction("ParamsMatch", ParamsMatchFunc)
	})
	_ = syncedEnforcer.LoadPolicy()
	return syncedEnforcer
}

func ParamsMatch(fullNameKey1 string, key2 string) bool {
	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用casbin的keyMatch2
	return util.KeyMatch2(key1, key2)
}

func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return ParamsMatch(name1, name2), nil
}
