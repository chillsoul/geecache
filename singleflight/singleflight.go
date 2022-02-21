package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//如果在map中找到正在进行的call
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() //阻塞，等查完再返回
		return c.val, c.err
	}
	//如果没有正在进行的call
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn() //执行查询数据
	c.wg.Done()         //查完释放锁
	g.mu.Lock()         //要把已经查完的删掉，以免死锁
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}
