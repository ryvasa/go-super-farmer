package casbin

import "github.com/casbin/casbin/v2"

type Casbin struct {
	enforcer *casbin.Enforcer
}

func Init(modelPath string, policyPath string) (*casbin.Enforcer, error) {
	e, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}
	return e, nil
}
