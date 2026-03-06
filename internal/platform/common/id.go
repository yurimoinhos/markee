package common

import "github.com/svicknesh/tsid"

// ID é um value object que representa um identificador único baseado em TSID.
type ID struct {
	Value string
}

func GenID() ID {
	return ID{Value: tsid.NewDefault().Generate().String()}
}
