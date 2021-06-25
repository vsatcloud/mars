package casbin

import (
	"errors"
	"strings"
	"sync"

	"gorm.io/gorm"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
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
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == "admin"
`

var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func Casbin(db *gorm.DB) *casbin.SyncedEnforcer {
	once.Do(func() {
		//dsn := fmt.Sprintf("hoast=%s user=%s password=%s port=%s dbname=%s sslmode=disable", db.Host, db.User, db.Password, db.Port, db.Dbname)
		a, _ := gormadapter.NewAdapterByDB(db)
		m, _ := model.NewModelFromString(modelText)
		syncedEnforcer, _ = casbin.NewSyncedEnforcer(m, a)
		syncedEnforcer.AddFunction("ParamsMatch", ParamsMatchFunc)
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

func UpdateCasbin(db *gorm.DB, authorityID, path, method string) error {
	ClearCasbin(db, 0, authorityID)
	rules := [][]string{}
	rules = append(rules, []string{authorityID, path, method})
	e := Casbin(db)
	success, _ := e.AddPolicies(rules)
	if success == false {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

func ClearCasbin(db *gorm.DB, v int, p ...string) bool {
	e := Casbin(db)
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success

}

type CasbinInfo struct {
	Path   string
	Method string
}

func GetPolicyPathByAuthorityId(db *gorm.DB, authorityID string) (pathMaps []CasbinInfo) {
	e := Casbin(db)
	list := e.GetFilteredPolicy(0, authorityID)
	for _, v := range list {
		pathMaps = append(pathMaps, CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}
