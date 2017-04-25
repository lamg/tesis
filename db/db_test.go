package db

import (
	"github.com/lamg/tesis"
	a "github.com/stretchr/testify/assert"
	"testing"
)

/*
func TestConn(t *testing.T) {
	var r *sql.Rows
	var e error
	r, e = AllStudents()
	if a.NoError(t, e) {
		r.Close()
	}
}*/

func TestPageSlice(t *testing.T) {
	var s, r []interface{}
	s = make([]interface{}, 20)
	for i, _ := range s {
		s[i] = i
	}
	var n, l, m, ps int
	n = 3 // page size
	r, l, m, ps = pageSlice(s, n, 6)
	a.True(t, len(r) == 2 && l == 18 && m == 20 && ps == 7)
}

func TestPending(t *testing.T) {
	var um *UPRManager
	var e error
	var dm tesis.UserDB
	dm = &tesis.DummyManager{}
	um, e = NewUPRManager("dtFile.json", dm)
	var pd *tesis.PageD
	if a.NoError(t, e) {
		pd, e = um.Pending(0)
	}
	a.True(t, e == nil && pd != nil)
}

func TestPropose(t *testing.T) {
	var pr []tesis.Diff
	pr = []tesis.Diff{
		tesis.Diff{
			LDAPRec:  tesis.DBRecord{Id: "0", IN: "8901191122", Name: "LUIS"},
			DBRec:    tesis.DBRecord{Id: "0", IN: "8901191122", Name: "Luis"},
			Exists:   true,
			Mismatch: true,
			Src:      "SIGENU",
		},
	}
	//each proposed exists in pending
	var um *UPRManager
	var e error
	var dm tesis.UserDB
	dm = &tesis.DummyManager{}
	//FIXME no esta leyendo el archivo
	um, e = NewUPRManager("dtFile.json", dm)
	if a.NoError(t, e) {
		e = um.Propose("lamg", pr)
	}
	a.NoError(t, e)
	//each proposed doesn't exists in pending
}

func TestSymDiff0(t *testing.T) {
	var x, y []tesis.Diff
	x, y =
		[]tesis.Diff{
			tesis.Diff{
				LDAPRec:  tesis.DBRecord{Id: "0", IN: "8901191122", Name: "LUIS"},
				DBRec:    tesis.DBRecord{Id: "0", IN: "8901191122", Name: "Luis"},
				Exists:   true,
				Mismatch: true,
				Src:      "SIGENU",
			},
		},
		[]tesis.Diff{
			tesis.Diff{
				LDAPRec:  tesis.DBRecord{Id: "0", IN: "8901191122", Name: "LUIS"},
				DBRec:    tesis.DBRecord{Id: "0", IN: "8901191122", Name: "Luis"},
				Exists:   true,
				Mismatch: true,
				Src:      "SIGENU",
			},
		}
	a.True(t, x[0].DBRec.Equals(y[0].DBRec))
	a.True(t, x[0].Equals(y[0]))
	var u, v []tesis.Eq
	u, v = make([]tesis.Eq, len(x)), make([]tesis.Eq, len(y))
	u[0], v[0] = x[0], y[0]
	var w, z []tesis.Eq
	w, z = tesis.DiffInt(u, v)
	a.True(t, len(w) == 0 && len(z) == 1 &&
		z[0].Equals(x[0]))
}
