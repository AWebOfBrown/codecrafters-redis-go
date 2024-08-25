package stream

import "fmt"

type Stream struct {
	root *StreamNode
}

func NewStream() *Stream {
	root := NewStreamNode("", true, make(map[string]interface{}))
	return &Stream{
		root: root,
	}
}

// Todo: Handle inserting same prefix twice
func (s *Stream) Insert(id string, entries map[string]interface{}) error {
	if len(s.root.children) == 0 {
		s.root.isCompleted = false
		s.root.children[id] = NewStreamNode(id, true, entries)
		return nil
	}

	stack := make([]*StreamNode, 0)

	for _, v := range s.root.children {
		stack = append(stack, v)
	}
	ptr := 0
	parent := s.root

	// Loop through a list of nodes with maybe matching prefix
	for i := 0; i < len(stack); i++ {
		node := stack[i]

		// For each node prefix letter, compare with id[ptr]
		for index, r := range node.prefix {
			lookingFor := rune(id[ptr])

			// Abandon node if first letter is no match, children cannot match
			if r != lookingFor && index == 0 {
				break
			}

			// If there's a match and it's the end of the ID we're looking for
			if r == lookingFor && ptr == len(id)-1 {
				return fmt.Errorf("tried to insert existing ID, %s", id)
			}

			// If match and more characters to search
			if r == lookingFor && index < len(node.prefix)-1 && ptr < len(id)-1 {
				ptr++
				continue
			}

			// if match and end of prefix but node children exist
			if r == lookingFor && index == (len(node.prefix)-1) && len(node.children) >= 1 {
				for _, child := range node.children {
					stack = append(stack, child)
				}
				ptr++
				break
			}

			if index == (len(node.prefix) - 1) {
				restOfId := id[ptr+1:]
				newNode := NewStreamNode(restOfId, true, entries)
				node.children[restOfId] = newNode
			} else {
				// From here either split current ID and create two children, or add one
				newNode := NewStreamNode(id[ptr:], true, entries)
				node.isCompleted = false

				restOfNonMatchingPath := node.prefix[index:]
				nodeForExistingNonMatchingPath := NewStreamNode(restOfNonMatchingPath, true, node.entries)
				node.children[restOfNonMatchingPath] = nodeForExistingNonMatchingPath

				oldChildPrefix := node.prefix
				node.prefix = node.prefix[:index]
				parent.children[node.prefix] = node
				delete(parent.children, oldChildPrefix)

				newId := id[ptr:]
				node.children[newId] = newNode
			}
			return nil
		}
	}
	return nil
}

// func (s *Stream) SplitNode(node, parent, newChild *StreamNode) {

// }

func (s *Stream) Lookup(prefix string) *StreamNode {
	if s.root == nil || len(s.root.children) == 0 {
		return nil
	}

	stack := make([]*StreamNode, 0)

	for _, child := range s.root.children {
		stack = append(stack, child)
	}

	found := 0
	for i := 0; i < len(stack); i++ {
		curr := stack[i]

		for prefixIndex, char := range curr.prefix {
			lookingFor := rune(prefix[found])
			doneIfMatch := found == len(prefix)-1

			if lookingFor == char {
				if doneIfMatch && prefixIndex == len(curr.prefix)-1 {
					return curr
				}

				// If match and at end of prefix but not done, iterate children
				if prefixIndex == len(curr.prefix)-1 {
					found++
					for _, child := range curr.children {
						stack = append(stack, child)
					}
					continue
				}
				// else continue
				found++
			} else {
				fmt.Println("hi")
				break
			}
		}
	}
	return nil
}
