package main

import "fmt"

type Err struct {
	Pkg  string
	Info string
	Prev error
}

func (e *Err) Error() string {
	if e.Prev == nil {
		return fmt.Sprintf("%s: %s", e.Pkg, e.Info)
	}
	return fmt.Sprintf("%s: %s\n%v", e.Pkg, e.Info, e.Prev)
}

func me(err error, format string, args ...interface{}) *Err {
	if len(args) > 0 {
		return &Err{
			Pkg:  `v3-golang`,
			Info: fmt.Sprintf(format, args...),
			Prev: err,
		}
	}
	return &Err{
		Pkg:  `v3-golang`,
		Info: format,
		Prev: err,
	}
}

func ce(err error, format string, args ...interface{}) {
	if err != nil {
		panic(me(err, format, args...))
	}
}

func ct(err *error) {
	if p := recover(); p != nil {
		if e, ok := p.(error); ok {
			*err = e
		} else {
			panic(p)
		}
	}
}
