package cli

type Param struct {
	Binding     ParamBinding
	Description string
	Name        string
}

type ParamBinding interface {
	Assign(string) error
	Param()
	Reset()
	String() string
	Type() string
}
