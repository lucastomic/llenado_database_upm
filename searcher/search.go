package search

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/lucastomic/llenado_database_upm/models"
)

func GetMunicipios() []models.Municipio {
	c := colly.NewCollector()

	var municipios []models.Municipio
	c.OnHTML("#mw-content-text > div.mw-parser-output > center > table > tbody > tr", func(e *colly.HTMLElement) {
		municipio := models.Municipio{
			Nombre:    e.ChildText("td:nth-child(2) > a"),
			Provincia: e.ChildText("td:nth-child(4) > a"),
		}
		if municipio.Nombre != "" && municipio.Provincia != "" {
			municipios = append(municipios, municipio)
		}
	})

	c.Visit("https://es.wikipedia.org/wiki/Anexo:Municipios_de_España_por_población")

	return municipios
}

func GetCalles(municipios []models.Municipio) []models.Calle {
	c := colly.NewCollector()

	var calles []models.Calle
	c.OnHTML("body > fieldset > table > tbody > tr > td> a", func(e *colly.HTMLElement) {
		nombre := strings.Split(e.Text, "(")[0]
		var municipio string
		if len(calles) >= len(municipios) {
			municipio = municipios[len(calles)%(len(municipios))].Nombre

		} else {
			municipio = municipios[len(municipios)%(len(calles)+1)].Nombre

		}
		calle := models.Calle{
			Nombre:    nombre,
			Municipio: municipio,
		}
		// fmt.Println(nombre, " :", municipio)

		calles = append(calles, calle)
	})

	c.Visit("https://gestiona.comunidad.madrid/nomecalles/ListaCalles.icm?munic=000")

	return calles
}

func dniAleatorio() string {
	letras := []string{
		"T",
		"R",
		"W",
		"A",
		"G",
		"M",
		"Y",
		"F",
		"P",
		"D",
		"X",
		"B",
		"N",
		"J",
		"Z",
		"S",
		"Q",
		"V",
		"H",
		"L",
		"C",
		"K",
		"E",
		"T",
	}
	dni := ""
	total := 0
	for i := 0; i < 8; i++ {
		n := rand.Intn(9)
		dni += strconv.Itoa(n)
		total += n
	}

	letra := letras[total%23]
	return dni + letra
}
func GetCarteros() []models.Cartero {

	nombres := getNombres()
	apellidos := getApellidos()

	var carteros []models.Cartero
	for i := 0; i < 100; i++ {
		cartero := models.Cartero{
			Nombre:   nombres[rand.Intn(len(nombres))],
			Apellido: apellidos[rand.Intn(len(apellidos))] + " " + apellidos[rand.Intn(len(apellidos))],
			DNI:      dniAleatorio(),
		}
		carteros = append(carteros, cartero)
	}
	return carteros

}

func GetDirecciones(calles []models.Calle) []models.Direccion {
	var direcciones []models.Direccion
	letras := []byte("ABCDEFGH")
	for i := 0; i < 10000; i++ {
		calle := calles[rand.Intn(len(calles))]
		direcciones = append(direcciones, models.Direccion{
			Numero:    rand.Intn(120) + 1,
			Letra:     string(letras[rand.Intn(7)]),
			Piso:      rand.Intn(7) + 1,
			Portal:    rand.Intn(4) + 1,
			Calle:     calle.Nombre,
			Municipio: calle.Municipio,
		})

	}
	return direcciones
}

func GetRutasPre() []models.RutaPredefinida {
	var rutas []models.RutaPredefinida

	for i := 0; i < 300; i++ {
		rutas = append(rutas, models.RutaPredefinida{ID: i})
	}
	return rutas
}

func GetSegmentos(calles []models.Calle, rutas []models.RutaPredefinida, areaDeEnvio []models.AreaDeEnvio) []models.Segmento {
	var segmentos []models.Segmento
	for _, calle := range calles {
		inicio := 0
		final := rand.Intn(40-8) + 8
		segmentos = append(segmentos, models.Segmento{
			Inicio:        inicio,
			Final:         final,
			Calle:         calle.Nombre,
			Municipio:     calle.Municipio,
			Ruta:          rutas[rand.Intn(len(rutas))].ID,
			IDAreaDeEnvio: areaDeEnvio[rand.Intn(len(areaDeEnvio))].ID,
		})
		for i := 1; i < rand.Intn(10); i++ {
			inicio = final + 1
			final = inicio + rand.Intn(40)
			segmentos = append(segmentos, models.Segmento{
				Inicio:        inicio,
				Final:         final,
				Calle:         calle.Nombre,
				Municipio:     calle.Municipio,
				Ruta:          rutas[rand.Intn(len(rutas))].ID,
				IDAreaDeEnvio: areaDeEnvio[rand.Intn(len(areaDeEnvio))].ID,
			})
		}
	}
	return segmentos
}

func GetUsuarios(direcciones []models.Direccion) []models.Usuario {
	apellidos := getApellidos()
	nombres := getNombres()
	var usuarios []models.Usuario
	for i := 0; i < 1000; i++ {
		nombre := nombres[rand.Intn(len(nombres))]
		apellido := apellidos[rand.Intn(len(apellidos))] + " " + apellidos[rand.Intn(len(apellidos))]
		var dni string
		if rand.Intn(2) == 1 {
			dni = dniAleatorio()
		}
		var correo string
		if rand.Intn(2) == 1 {
			correo = strings.ToLower(strings.Replace(nombre+apellido+"@gmail.com", " ", "", -1))
		}

		var direccion models.Direccion
		// if rand.Intn(2) == 1 {
		direccion = direcciones[rand.Intn(len(direcciones))]
		// }
		usuario := models.Usuario{
			ID:        i,
			Nombre:    nombre,
			Apellidos: apellido,
			DNI:       dni,
			Direccion: direccion,
			Correo:    correo,
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios
}

func getNombres() []string {
	c := colly.NewCollector()
	ninas := colly.NewCollector()

	var nombres []string
	ninas.OnHTML("#__next > div > main > div:nth-child(2) > div.infinite-article.selected.article-nombres-de-nina > div > div.container.article-container > div.row.article-content > div.col-md-7.col-sm-12 > div > div.RichText__StyledWrapper-sc-14kpice-0.YYzDT.font-secondary.rich-text > div:nth-child(8) > ol > li> p", func(e *colly.HTMLElement) {
		nombres = append(nombres, e.Text)
	})
	c.OnHTML("body > div.mainContainer > div.articleContainer > main > article > div.innerArticle__body > div > div.innerArticle__content.newsType > div.newsType__content.news-body-complete > ul > li", func(e *colly.HTMLElement) {
		nombres = append(nombres, e.Text)
	})

	ninas.Visit("https://www.dodot.es/embarazo/nombres-de-bebes/articulo/nombres-de-nina")
	c.Visit("https://www.elconfidencial.com/alma-corazon-vida/2021-09-25/nombres-ninos-bebes-padres-originales-populares_3294675/")

	return nombres
}

func getApellidos() []string {
	c := colly.NewCollector()

	var apellidos []string

	c.OnHTML("#post-23617 > div > table > tbody > tr > td:nth-child(2)", func(e *colly.HTMLElement) {
		apellidos = append(apellidos, e.Text)
	})

	c.Visit("https://www.saberespractico.com/curiosidades/apellidos-mas-comunes-en-espana/")

	apellidos = apellidos[1:]
	return apellidos
}

func GetAutorizaA(usuarios []models.Usuario) []models.AutorizaA {
	var autorizaA []models.AutorizaA
	for i := 0; i < 100; i++ {
		autorizaA = append(autorizaA, models.AutorizaA{
			IDReceptor: usuarios[rand.Intn(len(usuarios)-1)].ID,
			IDEmisor:   usuarios[rand.Intn(len(usuarios)-1)].ID,
		})
	}
	return autorizaA
}

func getProv() []string {
	return []string{
		"MAD",
		"BAR",
		"VAL",
		"SEV",
		"EXT",
		"MAL",
		"ALI",
		"AST",
		"BIL",
		"CAN",
		"GAL",
		"ZAR",
		"AND",
		"BAD",
	}
}

func GetAreasDeEnvio(oficinas []models.Oficina) []models.AreaDeEnvio {
	var areas []models.AreaDeEnvio
	provincias := getProv()

	for _, prov := range provincias {
		for i := 0; i < rand.Intn(20); i++ {
			var numero string
			if i > 9 {
				numero = strconv.Itoa(i)
			} else {
				numero = "0" + strconv.Itoa(i)
			}
			areas = append(areas, models.AreaDeEnvio{
				ID:      "AR-" + prov + "-" + numero,
				Oficina: oficinas[rand.Intn(len(oficinas))].Codigo,
			})
		}
	}
	return areas

}

func GetAsosianA(car []models.Cartero, are []models.AreaDeEnvio) []models.AsosianA {
	var asos []models.AsosianA
	for _, area := range are {
		for i := 0; i < rand.Intn(10); i++ {
			asos = append(asos, models.AsosianA{
				IDAreaDeEnvio: area.ID,
				DNICartero:    car[rand.Intn(len(car))].DNI,
			})
		}
	}
	return asos
}

func GetRecogidas(carteros []models.Cartero, direcciones []models.Direccion) []models.Recogida {
	var recogidas []models.Recogida
	for i := 0; i < 4000; i++ {
		recogidas = append(recogidas, models.Recogida{
			Identificador: i,
			Fecha:         fechaAleatoria(),
			DNICartero:    carteros[rand.Intn(len(carteros))].DNI,
			Direccion:     direcciones[rand.Intn(len(direcciones))],
		})
	}
	return recogidas
}

func GetCentrosDeClasificacion(municipios []models.Municipio) []models.CentroDeClasificacion {
	var centros []models.CentroDeClasificacion
	for _, municipio := range municipios {
		if rand.Intn(5) == 1 {
			for i := 0; i < rand.Intn(4); i++ {
				var numero string

				numero = "0" + strconv.Itoa(i)

				centros = append(centros, models.CentroDeClasificacion{
					Nombre:          "CC-" + strings.ToUpper(municipio.Nombre[:3]) + "-" + numero,
					Codigo:          len(centros),
					NombreMunicipio: municipio.Nombre,
					NEnviosMax:      rand.Intn(600-200) + 200,
				})
			}
		}
	}

	return centros
}

func fechaAleatoria() time.Time {
	min := time.Date(2000, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2023, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func GetOficinas(direc []models.Direccion, centros []models.CentroDeClasificacion) []models.Oficina {
	var oficinas []models.Oficina
	for _, prov := range getProv() {
		for i := 0; i < rand.Intn(30); i++ {
			var numero string
			if i > 9 {
				numero = strconv.Itoa(i)
			} else {
				numero = "0" + strconv.Itoa(i)
			}
			oficinas = append(oficinas, models.Oficina{
				Codigo:    "OF-" + prov + "-" + numero,
				Centro:    centros[rand.Intn((len(centros)))],
				Direccion: direc[rand.Intn((len(direc)))],
			})
		}
	}
	return oficinas
}

func GetCoches(ofi []models.Oficina) []models.Coche {
	var coches []models.Coche
	for i := 0; i < 100; i++ {
		var numMat string
		for i := 0; i < 4; i++ {
			numMat += strconv.Itoa(rand.Intn(10))
		}
		coches = append(coches, models.Coche{
			CodigoOficina: ofi[rand.Intn(len(ofi))].Codigo,
			Capacidad:     rand.Float64() * float64(300),
			Matricula:     numMat + cadenaRandom(3),
		})
	}
	return coches
}

func cadenaRandom(n int) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetTrabajaEn(carteros []models.Cartero, oficinas []models.Oficina) []models.TrabajaEn {
	turnos := []string{"Mañana", "Tarde", "Parcial"}
	var trabaja []models.TrabajaEn
	for _, cartero := range carteros {
		turno := turnos[rand.Intn(3)]
		var horaComienzo int
		var horaFinal int
		switch turno {
		case "Mañana":
			horaComienzo = rand.Intn(9-6) + 6
			horaFinal = horaComienzo + 8
		case "Tarde":
			horaComienzo = rand.Intn(16-13) + 13
			horaFinal = horaComienzo + 8
		case "Parcial":
			horaComienzo = rand.Intn(24)
			horaFinal = (horaComienzo + 4) % 24
		}

		trabaja = append(trabaja, models.TrabajaEn{
			DNICartero:    cartero.DNI,
			CodigoOficina: oficinas[rand.Intn(len(oficinas))].Codigo,
			Horario:       strconv.Itoa(horaComienzo) + ":00-" + strconv.Itoa(horaFinal) + ":00",
			Turno:         turno,
			FechaInicio:   fechaAleatoria(),
		})

		if turno == "Parcial" && horaFinal < 15 {
			horaComienzo = (horaComienzo + 5) % 24
			horaFinal = (horaFinal + 5) % 24
			trabaja = append(trabaja, models.TrabajaEn{
				DNICartero:    cartero.DNI,
				CodigoOficina: oficinas[rand.Intn(len(oficinas))].Codigo,
				Horario:       strconv.Itoa(horaComienzo) + ":00-" + strconv.Itoa(horaFinal) + ":00",
				Turno:         turno,
				FechaInicio:   fechaAleatoria(),
			})
		}
	}
	return trabaja
}

func GetReparto(coches []models.Coche, rutas []models.RutaPredefinida, carteros []models.Cartero, oficinas []models.Oficina) []models.Reparto {
	var repartos []models.Reparto
	for i := 0; i < 3000; i++ {
		repartos = append(repartos, models.Reparto{
			ID:                i,
			IDRutaPredefinida: rutas[rand.Intn(len(rutas))].ID,
			MatriculaCoche:    coches[rand.Intn(len(coches))].Matricula,
			CodigoOficina:     oficinas[rand.Intn(len(oficinas))].Codigo,
			DNICartero:        carteros[rand.Intn(len(carteros))].DNI,
		})
	}
	return repartos
}

func GetCartasCertificadas(repartos []models.Reparto, usuarios []models.Usuario) []models.CartaCertificada {
	var cartas []models.CartaCertificada

	for i := 0; i < 10000; i++ {
		var numero string
		for i := 0; i < 10; i++ {
			numero += strconv.Itoa(rand.Intn(10))
		}

		urgencias := []int{1, 2, 3}

		var usuario1 models.Usuario
		usEncontrado := false
		for !usEncontrado {
			us := usuarios[rand.Intn(len(usuarios))]
			if us.DNI != "" && us.Correo != "" {
				usEncontrado = true
				usuario1 = us
			}
		}

		var usuario2 models.Usuario
		usEncontrado = false

		for !usEncontrado {
			us := usuarios[rand.Intn(len(usuarios))]
			if us.DNI != "" && us.Correo != "" {
				usEncontrado = true
				usuario2 = us
			}
		}
		cartas = append(cartas, models.CartaCertificada{
			Identificador: "CE" + numero,
			Fecha:         fechaAleatoria(),
			IDReparto:     repartos[rand.Intn(len(repartos))].ID,
			IDEmisor:      usuario1.ID,
			IDReceptor:    usuario2.ID,
			Urgencia:      urgencias[rand.Intn(len(urgencias))],
		})
	}
	return cartas
}

func GetCartas(repartos []models.Reparto, usuarios []models.Usuario) []models.Carta {
	var cartas []models.Carta

	for i := 0; i < 10000; i++ {
		var numero string
		for i := 0; i < 10; i++ {
			numero += strconv.Itoa(rand.Intn(10))
		}

		cartas = append(cartas, models.Carta{
			Identificador: "CT" + numero,
			Fecha:         fechaAleatoria(),
			IDReparto:     repartos[rand.Intn(len(repartos))].ID,
			IDEmisor:      usuarios[rand.Intn(len(usuarios))].ID,
			IDReceptor:    usuarios[rand.Intn(len(usuarios))].ID,
			Formato:       "A" + strconv.Itoa(rand.Intn(8)),
		})
	}
	return cartas
}

func GetPaquetes(repartos []models.Reparto, usuarios []models.Usuario) []models.Paquete {
	var paquetes []models.Paquete

	for i := 0; i < 10000; i++ {
		var numero string
		for i := 0; i < 10; i++ {
			numero += strconv.Itoa(rand.Intn(10))
		}

		paquetes = append(paquetes, models.Paquete{
			Identificador: "PQ" + numero,
			Fecha:         fechaAleatoria(),
			IDReparto:     repartos[rand.Intn(len(repartos))].ID,
			IDEmisor:      usuarios[rand.Intn(len(usuarios))].ID,
			IDReceptor:    usuarios[rand.Intn(len(usuarios))].ID,
			Peso:          rand.Float64() * 40,
			Dimensiones:   strconv.Itoa(100) + "x" + strconv.Itoa(100) + "x" + strconv.Itoa(100),
		})
	}
	return paquetes
}

func GetCuentaCon(segmentos []models.Segmento, rutas []models.RutaPredefinida) []models.CuentaCon {
	var cuentaCon []models.CuentaCon
	for _, ruta := range rutas {
		for i := 0; i < rand.Intn(50)+8; i++ {
			cuentaCon = append(cuentaCon, models.CuentaCon{
				Segmento: segmentos[rand.Intn(len(segmentos))],
				Orden:    i,
				Ruta:     ruta.ID,
			})
		}
	}
	return cuentaCon
}
