package go_gen_server

const (
	_call uint8 = 1
	_cast uint8 = 2
	_info uint8 = 3
)

type pid struct {
	opt   uint8
	name  string
	pname string
}
