package dcmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/text/dregex"
	"github.com/osgochina/donkeygo/text/dstr"
	"os"
	"strings"
)

type Parser struct {
	strict           bool
	parsedArgs       []string
	parsedOptions    map[string]string
	passedOptions    map[string]bool
	supportedOptions map[string]bool
	commandFuncMap   map[string]func()
}

// Parse creates and returns a new Parser with os.Args and supported options.
//
// Note that the parameter <supportedOptions> is as [option name: need argument], which means
// the value item of <supportedOptions> indicates whether corresponding option name needs argument or not.
//
// The optional parameter <strict> specifies whether stops parsing and returns error if invalid option passed.
func Parse(supportedOptions map[string]bool, strict ...bool) (*Parser, error) {
	return ParseWithArgs(os.Args, supportedOptions, strict...)
}

// ParseWithArgs creates and returns a new Parser with given arguments and supported options.
//
// Note that the parameter <supportedOptions> is as [option name: need argument], which means
// the value item of <supportedOptions> indicates whether corresponding option name needs argument or not.
//
// The optional parameter <strict> specifies whether stops parsing and returns error if invalid option passed.
func ParseWithArgs(args []string, supportedOptions map[string]bool, strict ...bool) (*Parser, error) {
	strictParsing := false
	if len(strict) > 0 {
		strictParsing = strict[0]
	}
	parser := &Parser{
		strict:           strictParsing,
		parsedArgs:       make([]string, 0),
		parsedOptions:    make(map[string]string),
		passedOptions:    supportedOptions,
		supportedOptions: make(map[string]bool),
		commandFuncMap:   make(map[string]func()),
	}
	for name, needArgument := range supportedOptions {
		for _, v := range strings.Split(name, ",") {
			parser.supportedOptions[strings.TrimSpace(v)] = needArgument
		}
	}
	for i := 0; i < len(args); {
		if option := parser.parseOption(args[i]); option != "" {
			array, _ := dregex.MatchString(`^(.+?)=(.+)$`, option)
			if len(array) == 3 {
				if parser.isOptionValid(array[1]) {
					parser.setOptionValue(array[1], array[2])
				}
			} else {
				if parser.isOptionValid(option) {
					if parser.isOptionNeedArgument(option) {
						if i < len(args)-1 {
							parser.setOptionValue(option, args[i+1])
							i += 2
							continue
						}
					} else {
						parser.setOptionValue(option, "")
						i++
						continue
					}
				} else {
					// Multiple options?
					if arr := parser.parseMultiOption(option); len(array) > 0 {
						for _, v := range arr {
							parser.setOptionValue(v, "")
						}
						i++
						continue
					} else if parser.strict {
						return nil, errors.New(fmt.Sprintf(`invalid option '%s'`, args[i]))
					}
				}
			}
		} else {
			parser.parsedArgs = append(parser.parsedArgs, args[i])
		}
		i++
	}

	return parser, nil
}

// parseMultiOption parses option to multiple valid options like: --dav.
// It returns nil if given option is not multi-option.
func (that *Parser) parseMultiOption(option string) []string {
	for i := 1; i <= len(option); i++ {
		s := option[:i]
		if that.isOptionValid(s) && !that.isOptionNeedArgument(s) {
			if i == len(option) {
				return []string{s}
			}
			array := that.parseMultiOption(option[i:])
			if len(array) == 0 {
				return nil
			}
			return append(array, s)
		}
	}
	return nil
}

//解析option
func (that *Parser) parseOption(argument string) string {
	array, _ := dregex.MatchString(`^\-{1,2}(.+)$`, argument)
	if len(array) == 2 {
		return array[1]
	}
	return ""
}

//判断option是否合法
func (that *Parser) isOptionValid(name string) bool {
	_, ok := that.supportedOptions[name]
	return ok
}

//判断是否需要这个option
func (that *Parser) isOptionNeedArgument(name string) bool {
	return that.supportedOptions[name]
}

// setOptionValue sets the option value for name and according alias.
func (that *Parser) setOptionValue(name, value string) {
	for optionName, _ := range that.passedOptions {
		array := dstr.SplitAndTrim(optionName, ",")
		for _, v := range array {
			if strings.EqualFold(v, name) {
				for _, v1 := range array {
					that.parsedOptions[v1] = value
				}
				return
			}
		}
	}
}

// GetOpt returns the option value named <name>.
func (that *Parser) GetOpt(name string, def ...string) string {
	if v, ok := that.parsedOptions[name]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetOptVar returns the option value named <name> as gvar.Var.
func (that *Parser) GetOptVar(name string, def ...interface{}) *dvar.Var {
	if that.ContainsOpt(name) {
		return dvar.New(that.GetOpt(name))
	}
	if len(def) > 0 {
		return dvar.New(def[0])
	}
	return dvar.New(nil)
}

// GetOptAll returns all parsed options.
func (that *Parser) GetOptAll() map[string]string {
	return that.parsedOptions
}

// ContainsOpt checks whether option named <name> exist in the arguments.
func (that *Parser) ContainsOpt(name string) bool {
	_, ok := that.parsedOptions[name]
	return ok
}

// GetArg returns the argument at <index>.
func (that *Parser) GetArg(index int, def ...string) string {
	if index < len(that.parsedArgs) {
		return that.parsedArgs[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetArgVar returns the argument at <index> as gvar.Var.
func (that *Parser) GetArgVar(index int, def ...string) *dvar.Var {
	return dvar.New(that.GetArg(index, def...))
}

// GetArgAll returns all parsed arguments.
func (that *Parser) GetArgAll() []string {
	return that.parsedArgs
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (that *Parser) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"parsedArgs":       that.parsedArgs,
		"parsedOptions":    that.parsedOptions,
		"passedOptions":    that.passedOptions,
		"supportedOptions": that.supportedOptions,
	})
}
