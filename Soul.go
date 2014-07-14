package main

type Soul struct {
	*Item
	*stats
}

type SoulPouch []*Soul

func (sp *SoulPouch) getList() (list string) {
	var slist []string
	for _, s := range *sp {

		if s.atk > s.def {
			slist = append(slist, "red soul")
		} else {
			slist = append(slist, "blue soul")
		}
	}

	for k, v := range slist {
		if k != len(slist)-1 {
			v += "\n"
		}
		list += v
	}

	return
}
