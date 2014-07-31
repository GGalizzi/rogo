package main

import "fmt"

type Soul struct {
	*Item
	*stats

	name string
}

type SoulPouch []*Soul

func (sp *SoulPouch) getList() (list string) {
	var slist []string
	if len(*sp) == 0 {
		log("There are no souls in your pouch.")
		return
	}
	for _, s := range *sp {

		slist = append(slist, s.name)
	}

	for k, v := range slist {
		if k != len(slist)-1 {
			v += "\n"
		}
		list += v
	}

	return
}

func (s *Soul) genName(entName string) string {
	if s.atk <= 1 && s.def <= 1 && s.maxhp <= 1 {
		return fmt.Sprintf("weak %s soul", entName)
	}

	if s.atk > s.def && s.atk > s.maxhp {
		return fmt.Sprintf("red %s soul", entName)
	}

	if s.maxhp > s.def+5 && s.maxhp > s.atk+5 {
		return fmt.Sprintf("green %s soul", entName)
	}

	return fmt.Sprintf("blue %s soul", entName)
}
