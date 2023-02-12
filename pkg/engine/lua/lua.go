package lua

import (
	"fmt"

	"github.com/sirupsen/logrus"
	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

const (
	fnWrapper = `
function getInput()
	json = require("json")
	payloadObject, err = json.decode(inputPayload)
	if err then
		return nil
	end
	return payloadObject
end

function evaluate()
	%s
end
`
	fnBody = `
payload = getInput()
%s
return false, "default error"
`
)

func Execute(payload, script string) (bool, error) {

	fnEvaluateBody := fmt.Sprintf(fnBody, script)
	fn := fmt.Sprintf(fnWrapper, fnEvaluateBody)

	logrus.Debugf("payload:\n %s", payload)
	logrus.Debugf("executing function evaluate() in script:\n %s", fn)

	return executeFn(fn, payload)
}

func executeFn(fn, payload string) (bool, error) {

	state := lua.NewState()
	defer state.Close()

	// load libraries
	// TODO: not all libraries should be loaded
	libs.Preload(state)

	// set globals and load function
	state.SetGlobal("inputPayload", lua.LString(payload))
	if err := state.DoString(fn); err != nil {
		return false, err
	}

	evaluate := state.GetGlobal("evaluate")
	if err := state.CallByParam(lua.P{
		Fn:      evaluate,
		NRet:    2,
		Protect: true,
	}); err != nil {
		return false, err
	}

	res := state.ToBool(1)
	msg := state.ToString(2)

	return res, fmt.Errorf(msg)
}
