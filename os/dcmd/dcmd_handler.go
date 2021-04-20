package dcmd

import "errors"

// BindHandle 为命令绑定执行方法
func BindHandle(cmd string, f func()) error {
	if _, ok := defaultCommandFuncMap[cmd]; ok {
		return errors.New("duplicated handle for command:" + cmd)
	} else {
		defaultCommandFuncMap[cmd] = f
	}

	return nil
}

// BindHandleMap 通过map传入批量绑定方法
func BindHandleMap(m map[string]func()) error {
	var err error
	for k, v := range m {
		if err = BindHandle(k, v); err != nil {
			return err
		}
	}
	return err
}

// RunHandle 运行某个方法
func RunHandle(cmd string) error {
	if handle, ok := defaultCommandFuncMap[cmd]; ok {
		handle()
	} else {
		return errors.New("no handle found for command:" + cmd)
	}
	return nil
}

// AutoRun 根据命令行传入自动执行绑定的方法
func AutoRun() error {
	if cmd := GetArg(1); cmd != "" {
		if handle, ok := defaultCommandFuncMap[cmd]; ok {
			handle()
		} else {
			return errors.New("no handle found for command:" + cmd)
		}
	} else {
		return errors.New("no command found")
	}
	return nil
}
