package semtools

import (
	"testing"
	"github.com/sirupsen/logrus"
)


func init() {

	logrus.SetLevel(logrus.DebugLevel)

}

func TestNewNamespace(t *testing.T) {

	n := NewEmptyNamespace()
	if n == nil {
		t.Errorf("NewEmptyNamespace() fails to create new namespace")
	}

	n = NewNamespace()
	if n == nil {
		t.Errorf("NewNamespace() fails to create new namespace")
	}
	if len(n.ListKeys()) <= 0 {
		t.Errorf("NewNamespace() doesn't contain any defaults")
	}

}

func TestListKeysValues(t *testing.T) {

	n := NewEmptyNamespace()
	n.Set("test", "fullvalue")

	if len(n.ListKeys()) != 1 || n.ListKeys()[0] != "test" {
		t.Errorf("ListKeys() returns unexpected data")
	}

	if len(n.ListValues()) != 1 || n.ListValues()[0] != "fullvalue" {
		t.Errorf("ListValues() returns unexpected data")
	}

}

func TestContainsGetSet(t *testing.T) {

	n := NewEmptyNamespace()

	n.Set("test", "fullvalue")
	if n.MustGet("test") != "fullvalue" {
		t.Errorf("MustGet() returns unexpected result")
	}
	if v, ok := n.Get("test"); !ok || v != "fullvalue" {
		t.Errorf("Get() returns unexpected result")
	}
	if _, ok := n.Get("test2"); ok {
		t.Errorf("Get() returns unexpected result")
	}

	if !n.Contains("test") {
		t.Errorf("Contains() returns unexpected result")
	}
	if n.Contains("test2") {
		t.Errorf("Contains() returns unexpected result")
	}

}

func TestInclude(t *testing.T) {

	n := NewEmptyNamespace()
	n2 := NewNamespace()

	n3 := n.Include(n2)

	if len(n3.ListKeys()) <= 0 {
		t.Errorf("Include() doesn't copy over any namespaces")
	}

}

func TestCopy(t *testing.T) {

	n1 := NewNamespace()

	n2 := n1.Copy()

	if len(n2.ListKeys()) <= 0 {
		t.Errorf("Copy() doesn't create a valid copy")
	}

}
