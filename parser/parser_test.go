package parser

import "testing"

func TestBasicLit(t *testing.T) {
	consumes(t, BasicLit, map[string]bool{
		``:    false,
		`1`:   true,
		`1.1`: true,
		`1i`:  true,
		`'a'`: true,
		`"a"`: true,
		`a`:   false,
		`_`:   false,
	})
}

func TestIdentifierList(t *testing.T) {
	consumes(t, IdentifierList, map[string]bool{
		`a`:   true,
		`a,a`: true,
		`1`:   false,
		`_`:   true,
	})
}

func TestTypeName(t *testing.T) {
	consumes(t, TypeName, map[string]bool{
		`a`:   true,
		`a.a`: true,
		`1`:   false,
		`_`:   true,
	})
}

func TestQualifiedIdent(t *testing.T) {
	consumes(t, QualifiedIdent, map[string]bool{
		`1`:   false,
		`a`:   false,
		`a.a`: true,
		`_.a`: false,
		`a._`: true,
		`_._`: false,
	})
}

func TestPackageName(t *testing.T) {
	consumes(t, PackageName, map[string]bool{
		`a`: true,
		`1`: false,
		`_`: false,
	})
}
