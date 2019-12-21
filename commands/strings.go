package commands

type RedisStringCommands interface {
	Set(k, v []byte) error
	SetEx(k, v []byte, exp int) error
}
