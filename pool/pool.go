package pool

//pool implements a generic resource allocator

type Pool struct {
	container []interface{}
}

func NewPool() *Pool {
	return &Pool{make([]interface{}, 0)}
}

func (p *Pool) AddToPool(resource interface{}) {
	p.container = append(p.container, resource)
}
func (p *Pool) IsPoolEmpty() bool {
	return len(p.container) == 0
}
func (p *Pool) GetResource() interface{} {
	topResource := p.container[len(p.container)-1]
	p.container = p.container[0 : len(p.container)-1]
	return topResource
}
func (p *Pool) AvaialbleResourceUnits() int {
	return len(p.container)
}
