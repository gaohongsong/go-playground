package emptystruct

type Set map[string]struct{}

func (set Set) Add(s string) {
	set[s] = struct{}{}
}

func (set Set) Remove(s string) {
	delete(set, s)
}

func (set Set) Has(s string) bool {
	_, ok := set[s]
	return ok
}
