package qa

import "os"

type Args []string

func (_ Def) Args() Args {
	return Args(os.Args[1:])
}
