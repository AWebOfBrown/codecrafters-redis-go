package stream

import (
	"testing"
)

func Test_s(t *testing.T) {
	t.Run("Test inserting to a stream", func(t *testing.T) {
		s := NewStream()

		m := make(map[string]interface{})
		m["stuff"] = "stuff"

		s.Insert("romanus", m)
		s.Insert("Alex", m)
		s.Insert("romane", m)
		s.Insert("romulus", m)
		s.Insert("rubicon", m)
		s.Insert("rubicundus", m)
		s.Insert("rubens", m)
		s.Insert("Allen", m)
		s.Insert("Hamish", m)
		s.Insert("ruber", m)

		roman := s.root.children["r"].children["om"].children["an"]
		romanus := roman.children["us"]
		romane := roman.children["e"]
		if romanus == nil || romane == nil {
			t.Errorf("Node compression not working as expected")
		}

		romulus := s.root.children["r"].children["om"].children["ulus"]
		if romulus == nil {
			t.Errorf("Node compression not working as expected")
		}

		rubic := s.root.children["r"].children["ub"].children["ic"]
		rubicundus := rubic.children["undus"]
		rubicon := rubic.children["on"]
		if rubicundus == nil || rubicon == nil {
			t.Errorf("Node compression not working as expected")
		}

		rube := s.root.children["r"].children["ub"].children["e"]
		rubens := rube.children["ns"]
		ruber := rube.children["r"]

		if rubens == nil || ruber == nil {
			t.Errorf("Node compression not working as expected")
		}
	})

	t.Run("Test something", func(t *testing.T) {
		s := NewStream()

		m := make(map[string]interface{})
		m["stuff"] = "stuff"

		s.Insert("Alexander", m)

		if s.root.children["Alexander"] == nil {
			t.Error("Failed to insert node")
		}

		s.Insert("Alessandro", m)

		if s.root.children["Alexand"] == nil {
			t.Error("Failed to split node after matching prefix")
		}
	})

	t.Run("Test looking up nodes of a stream", func(t *testing.T) {
		s := NewStream()
		m := make(map[string]interface{})

		s.Insert("Alexander", m)
		s.Insert("Alexandra", m)

		alexand := s.Lookup("Alexand")

		alexander := s.Lookup("Alexander")
		alexandra := s.Lookup("Alexandra")

		if alexand == nil {
			t.Error("Failed to lookup node")
		}

		if alexander == nil {
			t.Error("Failed to lookup node")
		}

		if alexandra == nil {
			t.Error("Failed to lookup node")
		}
	})
}
