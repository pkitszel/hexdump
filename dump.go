package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func getInFile(args []string) (io.Reader, error) {
	var r io.Reader
	var err error
	if len(args) == 0 {
		r = os.Stdin
	} else {
		r, err = os.Open(args[0])
	}
	if err != nil {
		r = bufio.NewReader(r)
	}
	return r, err
}

func repr(c byte) byte {
	if strconv.IsPrint(rune(c)) {
		return c
	}
	return '.'
}
func repr8(bs []byte) string {
	out := ""
	for _, b := range bs {
		out += fmt.Sprintf("%02x ", int(b))
	}
	return out
}
func reprTxt(bs []byte) string {
	out := "|"
	for _, b := range bs {
		out += string(repr(byte(b)))
	}
	return out + "|"
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func printLine(i int, buf []byte, out io.Writer) error {
	l := len(buf)
	m := min(8, l)
	txt := fmt.Sprintf("%08x  ", i)
	txt += repr8(buf[0:m]) + " " + repr8(buf[m:l])
	txt += strings.Repeat(" ", 60-len(txt))
	_, err := fmt.Fprintf(out, "%s%s\n", txt, reprTxt(buf))
	return err
}

func run(args []string, out io.Writer) error {
	in, err := getInFile(args)
	if err != nil {
		return err
	}

	var buf, prv []byte
	prevCmp := 1
	for i := 0; ; i += 16 {
		if len(buf) < 16 {
			buf = make([]byte, 16)
		}
		n, err := io.ReadFull(in, buf)
		if err == io.EOF {
			return nil
		}
		if n < 16 {
			buf = buf[:n]
		}
		cmp := bytes.Compare(buf, prv)
		if cmp != 0 {
			err = printLine(i, buf, out)
			if err != nil {
				return err
			}
		} else if prevCmp != 0 {
			_, err = out.Write([]byte("*\n"))
			if err != nil {
				return err
			}
		}
		prevCmp = cmp
		tmp := prv
		prv = buf
		buf = tmp
		if n < 16 {
			_, err = fmt.Fprintf(out, "%08x\n", i+n)
			return err
		}
	}
}

func main() {
	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	err := run(os.Args[1:], f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
