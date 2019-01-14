package sorts

//Int64KeyStruct the object to be sorted
type Int64KeyStruct struct {
	Key   int64
	Value interface{}
}

type Int64KeyInterface []Int64KeyStruct

func (a Int64KeyInterface) Len() int {
	return len(a)
}
func (a Int64KeyInterface) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a Int64KeyInterface) Less(i, j int) bool {
	return a[i].Key < a[j].Key
}
