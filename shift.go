package cli

func shiftArg(arg *string) string {
	str := *arg
	*arg = ""
	return str
}

func shiftArgs(args *[]string) string {
	str := (*args)[0]
	*args = (*args)[1:]
	return str
}
