package flags

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type MultiString []string

func (m *MultiString) String() string {
	return strings.Join(*m, "\n")
}

func (m *MultiString) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func (m *MultiString) Contains(value string) bool {

	for _, test := range *m {

		if test == value {
			return true
		}
	}

	return false
}

type MultiDSNString []map[string]string

func (m *MultiDSNString) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *MultiDSNString) Set(value string) error {

	value = strings.Trim(value, " ")
	pairs := strings.Split(value, " ")

	dict := make(map[string]string)

	for _, str_pair := range pairs {

		str_pair = strings.Trim(str_pair, " ")
		pair := strings.Split(str_pair, "=")

		if len(pair) != 2 {
			return errors.New("Invalid pair")
		}

		k := pair[0]
		v := pair[1]

		dict[k] = v
	}

	*m = append(*m, dict)
	return nil
}

type MultiInt64 []int64

func (m *MultiInt64) String() string {

	str_values := make([]string, len(*m))

	for i, v := range *m {
		str_values[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(str_values, "\n")
}

func (m *MultiInt64) Set(str_value string) error {

	value, err := strconv.ParseInt(str_value, 10, 64)

	if err != nil {
		return err
	}

	*m = append(*m, value)
	return nil
}

func (m *MultiInt64) Contains(value int64) bool {

	for _, test := range *m {

		if test == value {
			return true
		}
	}

	return false
}
