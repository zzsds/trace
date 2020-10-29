package queue

type Options {

}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}