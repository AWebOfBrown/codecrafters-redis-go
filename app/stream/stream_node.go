package stream

type StreamNode struct {
	prefix      string
	isCompleted bool
	children    map[string]*StreamNode
	entries     map[string]interface{}
}

func NewStreamNode(prefix string, isCompleted bool, entries map[string]interface{}) *StreamNode {
	sn := StreamNode{
		prefix:      prefix,
		isCompleted: isCompleted,
		entries:     entries,
		children:    make(map[string]*StreamNode),
	}
	return &sn
}

func (sn *StreamNode) GetChild(character byte) *StreamNode {
	for k, v := range sn.children {
		if k[0] == character {
			return v
		}
	}
	return nil
}
