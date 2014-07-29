package main

type Soul struct {
	*Item
	*stats
}

type SoulPouch []*Soul

func (sp *SoulPouch) getList() (list string) {
	var slist []string
	for _, s := range *sp {

		slist = append(slist, s.name())
	}

	for k, v := range slist {
		if k != len(slist)-1 {
			v += "\n"
		}
		list += v
	}

	return
}

func (s *Soul) name() string {
	if s.atk <= 1 && s.def <= 1 && s.maxhp <= 1 {
		return "weak soul"
	}

	if s.atk > s.def && s.atk > s.maxhp {
		return "red soul"
	}

	if s.maxhp > s.def+5 && s.maxhp > s.atk+5 {
		return "green soul"
	}

	return "blue soul"
}
