package resp

const (
	SimpleStrPrefix = "+"
	ErrorStrPrefix  = "-"
	IntegerPrefix   = ":"
	BulkStrPrefix   = "$"
	ArrayPrefix     = "*"

	EndStr      = "\r\n"
	NullBulkStr = "$-1\r\n"
	EmptyArray  = "*0\r\n"

	SimpleStrVal = '+'
	ErrorStrVal  = '-'
	IntVal       = ':'
	BulkStrVal   = '$'
	ArrayVal     = '*'
	CtrlVal      = '\r'
	LfVal        = '\n'
)
