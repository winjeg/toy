package commands

type RedisStringCommands interface {
	Set(k string, v []byte) error
	SetEx(k string, v []byte, exp int) error
	Get(k string) []byte
}
