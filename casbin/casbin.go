package casbin

import (
	"strings"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/vsatcloud/mars/models"
)

func Casbin(modelPath string, db models.Database) *casbin.Enforcer {
	a, _ := gormadapter.NewAdapter("postgresql", db.User+":"+db.Password+"@("+db.Host+")/"+db.Dbname, true)
	e := casbin.NewEnforcer(modelPath, a)
	e.AddFunction("ParamsMatch", ParamsMatchFunc)
	_ = e.LoadPolicy()
	return e
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
