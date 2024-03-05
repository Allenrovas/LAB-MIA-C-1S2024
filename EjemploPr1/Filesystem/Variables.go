package Filesystem

type Mount struct {
	Id          string
	LetterValor string
	Name        string
	Part_type   [1]byte
	Start       int32
	Size        int32
}

type usuarioActual struct {
	Uid int32
	Gid int32
	Grp string
	Usr string
	Pwd string
	Pid string
}

func NuevoUsuarioActual() usuarioActual {
	return usuarioActual{-1, -1, "", "", "", ""}
}

var particionesMontadas []Mount

var Usr_sesion usuarioActual = NuevoUsuarioActual()

func VerificarParticionMontada(id string) int {
	for i := 0; i < len(particionesMontadas); i++ {
		if particionesMontadas[i].Id == id {
			return i
		}
	}
	return -1
}
