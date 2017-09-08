package gedcom

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

type stringTestCase struct {
	tested   string
	expected string
	actual   string
}

type stringTestCases []stringTestCase

func (testCases stringTestCases) run(t *testing.T) {
	for _, tc := range testCases {
		if tc.expected != tc.actual {
			t.Errorf("%s was [%s], expected [%s]", tc.tested, tc.actual, tc.expected)
		}
	}
}

type intTestCase struct {
	format   string
	expected int
	actual   int
}

type intTestCases []intTestCase

func (testCases intTestCases) run(t *testing.T) {
	for _, tc := range testCases {
		if tc.expected != tc.actual {
			t.Fatalf(tc.format+", expected [%d]", tc.actual, tc.expected)
		}
	}
}

type boolTestCase struct {
	tested   string
	expected bool
	actual   bool
}

type boolTestCases []boolTestCase

func (testCases boolTestCases) run(t *testing.T) {
	for _, tc := range testCases {
		if tc.expected != tc.actual {
			t.Fatalf("%s was [%t], expected [%t]", tc.tested, tc.actual, tc.expected)
		}
	}
}

func TestMain(m *testing.M) {

	d := NewDecoder(bytes.NewReader(data))
	d.SetUnrecTagFunc(func(l int, t, v, x string) {
		if t[0:1] == "_" {
			return
		}
		fmt.Printf("Unrecognized: %d %s %s", l, t, v)
		if x != "" {
			fmt.Printf(" (%s)", x)
		}
		fmt.Println("")
	})

	var err error
	g, err = d.Decode()
	if err != nil {
		log.Fatal("Result of decoding gedcom gave error, expected no error")
	}

	retCode := m.Run()
	os.Exit(retCode)
}
