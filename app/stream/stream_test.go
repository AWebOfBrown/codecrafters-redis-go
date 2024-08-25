package stream

import (
	"fmt"
	"testing"
)

func Test_s(t *testing.T) {
	t.Run("Test inserting to a stream", func(t *testing.T) {
		s := NewStream()

		m := make(map[string]interface{})
		m["stuff"] = "stuff"

		s.Insert("Alexander", m)

		if s.root.children["Alexander"] == nil {
			t.Error("Failed to insert node")
		}

		s.Insert("Alexandra", m)

		res := s.Insert("Alexandrara", m)

		fmt.Println(res)
		if s.root.children["Alexand"] == nil {
			t.Error("Failed to split node after matching prefix")
		}
		if s.root.children["Alexand"].children["er"] == nil {
			t.Error("Failed to insert child after matching prefix")
		}
		if s.root.children["Alexand"].children["ra"] == nil {
			t.Error("Failed to insert child after matching prefix")
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
