package workerpool2

type Option func(*Pool)

func WithBlock(block bool) Option {
	return func(pool *Pool) {
		pool.block = block
	}
}
