package dcmd

import (
	"errors"
)

// BindHandle registers callback function <f> with <cmd>.
func (that *Parser) BindHandle(cmd string, f func()) error {
	if _, ok := that.commandFuncMap[cmd]; ok {
		return errors.New("duplicated handle for command:" + cmd)
	} else {
		that.commandFuncMap[cmd] = f
	}
	return nil
}

// BindHandleMap registers callback function with map <m>.
func (that *Parser) BindHandleMap(m map[string]func()) error {
	var err error
	for k, v := range m {
		if err = that.BindHandle(k, v); err != nil {
			return err
		}
	}
	return err
}

// RunHandle executes the callback function registered by <cmd>.
func (that *Parser) RunHandle(cmd string) error {
	if handle, ok := that.commandFuncMap[cmd]; ok {
		handle()
	} else {
		return errors.New("no handle found for command:" + cmd)
	}
	return nil
}

// AutoRun automatically recognizes and executes the callback function
// by value of index 0 (the first console parameter).
func (that *Parser) AutoRun() error {
	if cmd := that.GetArg(1); cmd != "" {
		if handle, ok := that.commandFuncMap[cmd]; ok {
			handle()
		} else {
			return errors.New("no handle found for command:" + cmd)
		}
	} else {
		return errors.New("no command found")
	}
	return nil
}
