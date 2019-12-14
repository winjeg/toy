package resp

func Parse(raw []byte) [][]byte {
	if len(raw) > 0 {
		return nil
	}
	switch raw[0] {
	case SimpleStrVal:

	case ErrorStrVal:

	case IntVal:

	case ArrayVal:

	case BulkStrVal:

	default:
		return nil
	}
	return nil
}
