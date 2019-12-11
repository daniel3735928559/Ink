package parser

import (
	"fmt"
	"strings"
	"strconv"
	"libink/universe"
	"github.com/docopt/docopt-go"
	"github.com/google/shlex"
)

func ParseEffect(effect_str string) (*universe.Effect, error) {
	effect_def := `effect

Usage: 
  effect <item> change <propname> <oldpropval> to <newpropval>
  effect <item> do <verb> <directobj>
  effect <item> add <propname> <num>
  effect <item> get <gotten>
  effect <item> drop <dropped>
  effect describe <message>
`
	es, err := shlex.Split(effect_str)
	if err != nil {
		return nil, err
	}
	args, err := docopt.ParseArgs(effect_def, es, "")
	if err != nil {
		return nil, err
	}
        if args["change"].(bool) {
		return &universe.Effect{Type: universe.EFFECT_CHANGE, ActorName: args["<item>"].(string), PropName: args["<propname>"].(string), OldVal: args["<oldpropval>"].(string), NewVal: args["<newpropval>"].(string)}, nil
	} else if args["do"].(bool) {
		return &universe.Effect{Type: universe.EFFECT_DO, ActorName: args["<item>"].(string), Verb: args["<verb>"].(string), DirObj: args["<directobj>"].(string)}, nil
	} else if args["add"].(bool) {
		to_add, err := strconv.Atoi(args["<num>"].(string))
		if err != nil {
			return nil, err
		}
		return &universe.Effect{Type: universe.EFFECT_ADD, ActorName: args["<item>"].(string), PropName: args["<propname>"].(string), ToAdd: to_add}, nil
	} else if args["get"].(bool) {
		return &universe.Effect{Type: universe.EFFECT_GET, ActorName: args["<item>"].(string), Target: args["<gotten>"].(string)}, nil
	} else if args["drop"].(bool) {
		return &universe.Effect{Type: universe.EFFECT_DROP, ActorName: args["<item>"].(string), Target: args["<gotten>"].(string)}, nil
	} else if args["describe"].(bool) {
		return &universe.Effect{Type: universe.EFFECT_DESCRIBE, Message: args["<message>"].(string)}, nil
	}
	return nil, nil
}

func ParseCondition(condition_str string) (*universe.Condition, error) {
	condition_def := `condition

Usage: 
  condition <item> property <propname> is <propval>
  condition <item> property <propname> greater than <propval>
  condition <item> property <propname> less than <propval>
  condition <item> has <child>
`
	cs, err := shlex.Split(condition_str)
	if err != nil {
		return nil, err
	}
	args, err := docopt.ParseArgs(condition_def, cs, "")
	if err != nil {
		return nil, err
	}
	
        if args["is"].(bool) {
		return &universe.Condition{Type: universe.CONDITION_PROPIS, PropName: args["<propname>"].(string), PropVal: args["<propval>"].(string)}, nil
	} else if args["greater"].(bool) {
		return &universe.Condition{Type: universe.CONDITION_PROPGT, PropName: args["<propname>"].(string), PropVal: args["<propval>"].(string)}, nil
	} else if args["less"].(bool) {
		return &universe.Condition{Type: universe.CONDITION_PROPLT, PropName: args["<propname>"].(string), PropVal: args["<propval>"].(string)}, nil
	} else if args["has"].(bool) {
		return &universe.Condition{Type: universe.CONDITION_HAS, ChildName: args["<child>"].(string)}, nil
	}
	return nil, nil
}

func ParseGame(game string) *universe.Universe {
	obj_def := `object

Usage:
  obj <item> is <description>
  obj <item> has <child>
  obj <item> property <propname> can be <propval> <description>
  obj <item> property <propname> num <description>
  obj <item> property <propname> max <propval>
  obj <item> property <propname> min <propval>
  obj <item> property <propname> is <propval>
  obj <item> can <verb> <directobj> then <effects>...
  obj <item> can <verb> <directobj> if <condition> then <effects>...
  obj win <description> if <condition>...
  obj lose <description> if <condition>...
`
	u := universe.MakeUniverse()
	
	for lineno, l := range strings.Split(game, "\n") {
		fmt.Println("Line:",l)
		line, _ := shlex.Split(l)
		args, err := docopt.ParseArgs(obj_def, line, "")
		if err != nil {
			fmt.Println("Line",lineno,":",l,">>> Error:",err)
		} else {
			if args["win"].(bool) {
				conds := []*universe.Condition{}
				for _, cond_str := range args["<condition>"].([]string) {
					cond, err := ParseCondition(cond_str)
					if err != nil {
						fmt.Println("Line",lineno,":",l,">>> Condition: ",args["<condition>"].(string),">>> Error:",err)
						continue
					}
					conds = append(conds, cond)
				}
				u.AddWinCondition(args["<description>"].(string), conds)
			} else if args["lose"].(bool) {
				conds := []*universe.Condition{}
				for _, cond_str := range args["<condition>"].([]string) {
					cond, err := ParseCondition(cond_str)
					if err != nil {
						fmt.Println("Line",lineno,":",l,">>> Condition: ",args["<condition>"].(string),">>> Error:",err)
						continue
					}
					conds = append(conds, cond)
				}
				u.AddLoseCondition(args["<description>"].(string), conds)
			} else {
				item := u.GetItem(args["<item>"].(string))
				if args["is"].(bool) && !args["property"].(bool) {
					item.SetDescription(args["<description>"].(string))
				} else if args["has"].(bool) {
					item.AddChild(u.GetItem(args["<child>"].(string)))
				} else if args["can"].(bool) && args["be"].(bool) {
					item.AddStateProperty(args["<propname>"].(string), args["<propval>"].(string), args["<description>"].(string))
				} else if args["num"].(bool) {
					item.AddNumProperty(args["<propname>"].(string), args["<description>"].(string))
				} else if args["max"].(bool) {
					val, _ := strconv.Atoi(args["<propval>"].(string))
					item.SetNumPropertyMax(args["<propname>"].(string), val)
				} else if args["min"].(bool) {
					val, _ := strconv.Atoi(args["<propval>"].(string))
					u.GetItem(args["<item>"].(string)).SetNumPropertyMin(args["<propname>"].(string), val)
				} else if args["property"].(bool) && args["is"].(bool) {
					if item.HasStateProperty(args["<propname>"].(string)) {
						item.SetStatePropertyValue(args["<propname>"].(string), args["<propval>"].(string))
					} else if item.HasNumProperty(args["<propname>"].(string)) {
						val, _ := strconv.Atoi(args["<propval>"].(string))
						item.SetNumPropertyValue(args["<propname>"].(string), val)
					}
				} else if args["can"].(bool) && !args["if"].(bool) && args["then"].(bool) {
					effects := []*universe.Effect{}
					for _, es := range args["<effects>"].([]string) {
						eff, err := ParseEffect(es)
						if err != nil {
							fmt.Println("Line",lineno,":",l,">>> effect: ",es,">>> Error:",err)
							continue
						}
						effects = append(effects, eff)
					}
					item.AddAction(args["<verb>"].(string), u.GetItem(args["<directobj>"].(string)), effects, nil)
				} else if args["can"].(bool) && !args["if"].(bool) && args["then"].(bool) {
					effects := []*universe.Effect{}
					for _, es := range args["<effects>"].([]string) {
						eff, err := ParseEffect(es)
						if err != nil {
							fmt.Println("Line",lineno,":",l,">>> effect: ",es,">>> Error:",err)
							continue
						}
						effects = append(effects, eff)
					}
					cond, err := ParseCondition(args["<condition>"].(string))
					if err != nil {
						fmt.Println("Line",lineno,":",l,">>> Condition: ",args["<condition>"].(string),">>> Error:",err)
						continue
					}
					item.AddAction(args["<verb>"].(string), u.GetItem(args["<directobj>"].(string)), effects, cond)
				} else {
					fmt.Println("Line",lineno,":",l,">>> Something unexpected happened. Args:",args)
				}
			} 
		}
	}
	return u
}
