package util

import "fmt"

type SliceFlag []string

func (f *SliceFlag) String() string {
	return fmt.Sprintf("%v", []string(*f))
}

func (f *SliceFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}
