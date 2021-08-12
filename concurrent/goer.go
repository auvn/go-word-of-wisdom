package concurrent

type Goer interface {
	Go(fn func())
}
