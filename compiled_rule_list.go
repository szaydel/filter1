package main

// CompiledRuleList groups together multiple rules read in from configuration
// file.
type CompiledRuleList struct {
	list []CompiledRule
}

func (c *CompiledRuleList) Len() int           { return len(c.list) }
func (c *CompiledRuleList) Swap(i, j int)      { c.list[i], c.list[j] = c.list[j], c.list[i] }
func (c *CompiledRuleList) Less(i, j int) bool { return c.list[i].desc.Comment < c.list[j].desc.Comment }

// Append appends elem, a CompiledRule, to the CompiledRuleList.
func (c *CompiledRuleList) Append(elem CompiledRule) {
	c.list = append(c.list, elem)
}

// Slice returns the underlying data as a []CompiledRule type.
func (c *CompiledRuleList) Slice() []CompiledRule {
	return c.list
}
