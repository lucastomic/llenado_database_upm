package models

import "time"

type Insercion interface {
	// Columnas en el orden que seran insertadas
	GetProps() []any
}

type Municipio struct {
	Provincia string
	Nombre    string
}

func (m Municipio) GetProps() []any {
	return []any{m.Nombre, m.Provincia}
}

type Calle struct {
	Nombre    string
	Municipio string
}

func (c Calle) GetProps() []any {
	return []any{c.Nombre, c.Municipio}
}

type Cartero struct {
	Nombre   string
	Apellido string
	DNI      string
}

func (c Cartero) GetProps() []any {
	return []any{c.Nombre, c.Apellido, c.DNI}
}

type Direccion struct {
	Numero    int
	Letra     string
	Piso      int
	Portal    int
	Calle     string
	Municipio string
}

func (d Direccion) GetProps() []any {
	return []any{d.Numero, d.Letra, d.Piso, d.Portal, d.Calle, d.Municipio}
}

type RutaPredefinida struct {
	ID int
}

func (r RutaPredefinida) GetProps() []any {
	return []any{r.ID}
}

type Segmento struct {
	Inicio        int
	Final         int
	Calle         string
	Municipio     string
	Ruta          int
	IDAreaDeEnvio string
}

func (s Segmento) GetProps() []any {
	return []any{s.Inicio, s.Final, s.Calle, s.Municipio, s.IDAreaDeEnvio}
}

type Usuario struct {
	ID        int
	Nombre    string
	Apellidos string
	DNI       string
	Correo    string
	Direccion Direccion
}

func (u Usuario) GetProps() []any {
	// if u.Direccion.Calle != "" {
	return []any{u.ID, u.Nombre, u.Apellidos, u.DNI, u.Correo, u.Direccion.Numero, u.Direccion.Letra, u.Direccion.Piso, u.Direccion.Portal, u.Direccion.Calle, u.Direccion.Municipio}
	// } else {
	// 	return []any{u.ID, u.Nombre, u.Apellidos, u.DNI, u.Correo, nil, nil, nil, nil, nil, nil}
	// }
}

type AutorizaA struct {
	IDReceptor int
	IDEmisor   int
}

func (a AutorizaA) GetProps() []any {
	return []any{a.IDEmisor, a.IDReceptor}
}

type AreaDeEnvio struct {
	ID      string
	Oficina string
	IDPadre string
}

func (a AreaDeEnvio) GetProps() []any {
	return []any{a.ID, a.Oficina}
}

type AsosianA struct {
	IDAreaDeEnvio string
	DNICartero    string
}

func (a AsosianA) GetProps() []any {
	return []any{a.IDAreaDeEnvio, a.DNICartero}
}

type Recogida struct {
	Identificador int
	Fecha         time.Time
	DNICartero    string
	Direccion     Direccion
}

func (r Recogida) GetProps() []any {
	return []any{r.Identificador, r.Fecha, r.DNICartero, r.Direccion.Numero, r.Direccion.Letra, r.Direccion.Piso, r.Direccion.Portal, r.Direccion.Calle, r.Direccion.Municipio}
}

type CentroDeClasificacion struct {
	Codigo          int
	NombreMunicipio string
	Nombre          string
	NEnviosMax      int
}

func (c CentroDeClasificacion) GetProps() []any {
	return []any{c.Codigo, c.NombreMunicipio, c.Nombre, c.NEnviosMax}
}

type Oficina struct {
	Codigo    string
	Centro    CentroDeClasificacion
	Direccion Direccion
}

func (c Oficina) GetProps() []any {
	return []any{c.Codigo, c.Centro.Codigo, c.Direccion.Numero, c.Direccion.Letra, c.Direccion.Piso, c.Direccion.Portal, c.Direccion.Calle, c.Direccion.Municipio}
}

type Coche struct {
	Matricula     string
	CodigoOficina string
	Capacidad     float64
}

func (c Coche) GetProps() []any {
	return []any{c.Matricula, c.CodigoOficina, c.Capacidad}
}

type TrabajaEn struct {
	DNICartero    string
	CodigoOficina string
	Horario       string
	FechaInicio   time.Time
	Turno         string
}

func (c TrabajaEn) GetProps() []any {
	return []any{c.DNICartero, c.CodigoOficina, c.Horario, c.FechaInicio, c.Turno}
}

type Reparto struct {
	IDRutaPredefinida int
	ID                int
	MatriculaCoche    string
	CodigoOficina     string
	DNICartero        string
}

func (r Reparto) GetProps() []any {
	return []any{r.IDRutaPredefinida, r.ID, r.MatriculaCoche, r.CodigoOficina, r.DNICartero}
}

type CartaCertificada struct {
	Identificador string
	Fecha         time.Time
	IDReparto     int
	Urgencia      int
	IDEmisor      int
	IDReceptor    int
}

func (c CartaCertificada) GetProps() []any {
	return []any{c.Identificador, c.Fecha, c.IDReparto, c.Urgencia, c.IDEmisor, c.IDReceptor}
}

type Carta struct {
	Identificador string
	Fecha         time.Time
	IDReparto     int
	Formato       string
	IDEmisor      int
	IDReceptor    int
}

func (c Carta) GetProps() []any {
	return []any{c.Identificador, c.Fecha, c.IDReparto, c.Formato, c.IDEmisor, c.IDReceptor}
}

type Paquete struct {
	Identificador string
	Fecha         time.Time
	IDReparto     int
	Comentario    string
	Peso          float64
	Dimensiones   string
	IDEmisor      int
	IDReceptor    int
}

func (c Paquete) GetProps() []any {
	return []any{c.Identificador, c.Fecha, c.IDReparto, c.Comentario, c.Peso, c.Dimensiones, c.IDEmisor, c.IDReceptor}
}

type CuentaCon struct {
	Ruta     int
	Segmento Segmento
	Orden    int
}

func (c CuentaCon) GetProps() []any {
	return []any{c.Ruta, c.Segmento.Inicio, c.Segmento.Final, c.Segmento.Calle, c.Segmento.Municipio, c.Orden}
}
