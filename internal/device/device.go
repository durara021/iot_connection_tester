package device

type Device interface {
	Connect() error
	Close() error
}
