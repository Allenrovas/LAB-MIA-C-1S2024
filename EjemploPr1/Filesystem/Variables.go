package Filesystem

type Mount struct {
	Id          string
	LetterValor string
	Name        string
	Part_type   [1]byte
	Start       int32
	Size        int32
}

var particionesMontadas []Mount

func VerificarParticionMontada(id string) int {
	for i := 0; i < len(particionesMontadas); i++ {
		if particionesMontadas[i].Id == id {
			return i
		}
	}
	return -1
}
