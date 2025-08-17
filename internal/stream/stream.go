package stream

import (
	"fmt"
)

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
	var parent = s.root
	// Loop through a list of nodes with maybe matching prefix
	for i := 0; i < len(stack); i++ {
		node := stack[i]

		// For each node prefix letter, compare with id[ptr]
		for index, r := range node.prefix {
			lookingFor := rune(id[ptr])
			lookingForStr := string(lookingFor)
			fmt.Printf(lookingForStr)

			// Abandon node if first letter is no match, children cannot match
			if r != lookingFor && index == 0 {
				break
			}

			if r == lookingFor {
				// If there's a match and it's the end of the ID we're looking for
				if ptr == len(id)-1 {
					return fmt.Errorf("tried to insert existing ID, %s", id)
				}

				// If match and more characters to search
				if index < len(node.prefix)-1 && ptr < len(id)-1 {
					ptr++
					continue
				}

				// if match and end of prefix but node children exist
				if index == (len(node.prefix)-1) && len(node.children) >= 1 {
					for _, child := range node.children {
						stack = append(stack, child)
					}
					parent = node
					ptr++
					break
				}
			} else {
				// If not a match but last character of current node, simply create new node.
				if index == (len(node.prefix) - 1) {
					restOfId := id[ptr+1:]
					newNode := NewStreamNode(restOfId, true, entries)
					node.children[restOfId] = newNode
					return nil
				} else {
					// If not a match and not final character of current node
					// Split the current node and move the rest to a new child node with current children / entries
					prefixOfRestOfCurrentNonMatchingNode := node.prefix[index:]
					nodeForAbovePrefix := NewStreamNode(prefixOfRestOfCurrentNonMatchingNode, node.isCompleted, node.entries)
					nodeForAbovePrefix.children = node.children
					node.children = make(map[string]*StreamNode)
					node.children[prefixOfRestOfCurrentNonMatchingNode] = nodeForAbovePrefix

					// Change the current node's prefix to be up-to the index
					parent.children[node.prefix[:index]] = node
					delete(parent.children, node.prefix)

					node.prefix = node.prefix[:index]
					node.entries = make(map[string]interface{})
					node.isCompleted = false

					// Then insert the rest of the new prefix as a new child of the current node.
					newNode := NewStreamNode(id[ptr:], true, entries)
					newNode.isCompleted = true
					node.children[id[ptr:]] = newNode
					return nil
				}
			}
			return nil
		}
	}
	s.root.children[id] = NewStreamNode(id, true, entries)
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
