package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"unicode/utf8"
)

var commands = map[int]func(*Interp) error{
	int('h'): func(i *Interp) error { i.p = i.s; return nil },
	int('c'): func(i *Interp) error { i.r = i.memory[i.s]; return nil },
	int('.'): func(i *Interp) error { print(i.memory[i.s]); return nil },
	int('+'): func(i *Interp) error { i.memory[i.s]++; return nil },
	int('-'): func(i *Interp) error { i.memory[i.s]--; return nil },
	int('>'): func(i *Interp) error { i.s++; return nil },
	int('<'): func(i *Interp) error { i.s--; return nil },
	int('%'): func(i *Interp) error { return errSuccess },
}

var (
	errSuccess = errors.New("success")
)

type Interp struct {
	memory     []int
	f, s, p, r int
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		return
	}
	buf, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
	}
	interp := NewInterp(string(buf))
	interp.Run()
}

func NewInterp(code string) *Interp {
	p := new(Interp)
	p.memory = make([]int, 4096)
	var i int
	var comment bool
	for _, c := range code {
		if comment {
			if c == '\n' {
				comment = false
				continue
			}
		} else {
			if c == '#' {
				comment = true
				continue
			}
			if utf8.RuneLen(c) != 1 {
				continue
			}
			if c == '%' && p.s == 0 {
				p.s = i + 1
			}
			if commands[int(c)] != nil {
				p.memory[i] = int(c)
				i++
			}
		}
	}
	if p.s == 0 {
		p.memory[i] = int('%')
		p.s = i + 1
	}
	return p
}

func (i *Interp) Run() {
	for {
		err := commands[i.memory[i.f]](i)
		if err != nil {
			if err != errSuccess {
				fmt.Println(err)
			}
			return
		}
		i.f++
	}
}

func print(c int) {
	fmt.Printf("%c", c)
}
