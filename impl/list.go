package impl

type StrImpl struct {
}

func (s *StrImpl) Set(k, v []byte) error {
	return nil
}

func (s *StrImpl) SetEx(k, v []byte, exp int) error {
	return nil
}
