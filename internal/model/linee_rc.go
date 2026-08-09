package model

type LineeRC struct {
	Linea string `db:"linea"`
	Rc string `db:"rc"`
	Disp_linea *string `db:"disp_linea"`
	Disp_dest *string `db:"disp_dest"`
	Desc *string `db:"desc"`
}