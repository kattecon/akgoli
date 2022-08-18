package utils

type ConstError string

func (e ConstError) Error() string { return string(e) }
