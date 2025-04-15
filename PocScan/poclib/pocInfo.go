package poclib

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type Poc struct {
	Name   string  `yaml:"name"`
	Set    StrMap  `yaml:"set"`
	Sets   ListMap `yaml:"sets"`
	Rules  []Rules `yaml:"rules"`
	Groups RuleMap `yaml:"groups"`
	//Detail Detail  `yaml:"detail"`
}

type Rules struct {
	Method          string            `yaml:"method"`
	Path            string            `yaml:"path"`
	Headers         map[string]string `yaml:"headers"`
	Body            string            `yaml:"body"`
	Search          string            `yaml:"search"`
	FollowRedirects bool              `yaml:"follow_redirects"`
	Expression      string            `yaml:"expression"`
	Continue        bool              `yaml:"continue"`
}

//type Detail struct {
//	Author      string   `yaml:"author"`
//	Links       []string `yaml:"links"`
//	Description string   `yaml:"description"`
//	Version     string   `yaml:"version"`
//}

type MapSlice = yaml.MapSlice

type StrMap []StrItem
type ListMap []ListItem
type RuleMap []RuleItem

type StrItem struct {
	Key, Value string
}

type ListItem struct {
	Key   string
	Value []string
}

type RuleItem struct {
	Key   string
	Value []Rules
}

func (r *RuleMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp1 yaml.MapSlice
	if err := unmarshal(&tmp1); err != nil {
		return err
	}
	var tmp = make(map[string][]Rules)
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	for _, one := range tmp1 {
		key := one.Key.(string)
		value := tmp[key]
		*r = append(*r, RuleItem{key, value})
	}
	return nil
}

func (r *StrMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp yaml.MapSlice
	if err := unmarshal(&tmp); err != nil {
		return err
	}
	for _, one := range tmp {
		key, value := one.Key.(string), one.Value.(string)
		*r = append(*r, StrItem{key, value})
	}
	return nil
}

func (r *ListMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp yaml.MapSlice
	if err := unmarshal(&tmp); err != nil {
		return err
	}
	for _, one := range tmp {
		key := one.Key.(string)
		var value []string
		for _, val := range one.Value.([]interface{}) {
			v := fmt.Sprintf("%v", val)
			value = append(value, v)
		}
		*r = append(*r, ListItem{key, value})
	}
	return nil
}
