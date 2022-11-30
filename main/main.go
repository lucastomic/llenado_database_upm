package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lucastomic/llenado_database_upm/models"
	search "github.com/lucastomic/llenado_database_upm/searcher"
)

func getDatabase() (db *sql.DB, e error) {
	usuario := "root"
	pass := ""
	host := "tcp(127.0.0.1:3306)"
	dbname := "correos"
	// Debe tener la forma usuario: contraseña@protocolo(host:puerto)/nombreBaseDeDatos
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, dbname))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initSchema(db *sql.DB) error {
	nombreArchivo := "create_schema.sql"
	bytesLeidos, err := ioutil.ReadFile(nombreArchivo)
	if err != nil {
		fmt.Printf("Error leyendo archivo: %v \n", err)
	}

	sqlSentence := string(bytesLeidos)

	// Preparamos para prevenir inyecciones SQL
	sentenciaPreparada, err := db.Prepare(sqlSentence)
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()
	// Ejecutar sentencia, un valor por cada '?'
	_, err = sentenciaPreparada.Exec()
	if err != nil {
		return err
	}
	return nil
}

func insertar(db *sql.DB, insertSentence string, inserciones []models.Insercion) []error {
	var errores []error
	for _, insercion := range inserciones {
		sentenciaPreparada, err := db.Prepare(insertSentence)
		if err != nil {
			errores = append(errores, err)
		}

		_, err = sentenciaPreparada.Exec(insercion.GetProps()...)
		if err != nil {
			// fmt.Println(insercion.GetProps()...)
			errores = append(errores, err)
		}
		defer sentenciaPreparada.Close()
	}
	return errores
}

func main() {

	db, err := getDatabase()

	if err != nil {
		fmt.Printf("Error obteniendo base de datos: %v \n", err)
		return
	}

	// Terminar conexión al terminar función
	defer db.Close()

	var wg sync.WaitGroup

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error conectando: %v \n", err)
		return
	}

	// Listo, aquí ya podemos usar a db
	fmt.Printf("Conectado correctamente \n")

	var errores []error
	// Iniciamos el Schema si es que no existe
	// err = initSchema(db)
	// if err != nil {
	// 	fmt.Printf("Error creando la base de datos: %v \n", err)
	// 	return
	// }

	// Insertar municipios

	var municipiosSinProcesar []models.Municipio
	wg.Add(1)
	go func() {
		defer wg.Done()

		var municipios []models.Insercion
		municipiosSinProcesar = search.GetMunicipios()
		for _, val := range municipiosSinProcesar {
			municipios = append(municipios, val)
		}
		errores = insertar(db, "INSERT INTO municipio(nombre,provincia) VALUES (?,?)", municipios)

		if len(errores) != 0 {
			for _, e := range errores {
				fmt.Println("Error al insertar municipio: ", e)
			}
		}
	}()

	// Insertar Cartero
	var carterosSinProcesar []models.Cartero

	wg.Add(1)
	go func() {
		defer wg.Done()

		var carteros []models.Insercion
		carterosSinProcesar = search.GetCarteros()
		for _, val := range carterosSinProcesar {
			carteros = append(carteros, val)
		}
		errores = insertar(db, "INSERT INTO cartero(nombre,apellidos, DNI) VALUES (?,?,?)", carteros)
		if len(errores) != 0 {
			for _, e := range errores {
				fmt.Println("Error al insertar cartero: ", e)
			}
		}
	}()

	// Insertar rutas predefinidas
	var rutasPreSinProcesar []models.RutaPredefinida
	wg.Add(1)
	go func() {
		defer wg.Done()
		var rutas []models.Insercion
		rutasPreSinProcesar = search.GetRutasPre()
		for _, val := range rutasPreSinProcesar {
			rutas = append(rutas, val)
		}
		errores = insertar(db, "INSERT INTO rutaPredefinida(ID) VALUES (?)", rutas)
		if len(errores) != 0 {
			for _, e := range errores {
				fmt.Println("Error al insertar ruta: ", e)
			}
		}

	}()

	wg.Wait()

	// Insertar calle
	var calles []models.Insercion
	callesSinProcesar := search.GetCalles(municipiosSinProcesar)
	for _, val := range callesSinProcesar {
		calles = append(calles, val)
	}
	errores = insertar(db, "INSERT INTO calle(nombre,nombreMunicipio) VALUES (?,?)", calles)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar calle: ", e)
		}
	}

	// Insertar direcciones

	var direcciones []models.Insercion
	direccionesSinProcesar := search.GetDirecciones(callesSinProcesar)
	for _, val := range direccionesSinProcesar {
		direcciones = append(direcciones, val)
	}
	errores = insertar(db, "INSERT INTO direccion(numero,letra, piso, portal, nombreCalle, nombreMunicipio) VALUES (?,?,?,?,?,?)", direcciones)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar dirección: ", e)
		}
	}

	// Insertar usuarios
	var usuarios []models.Insercion
	usuariosSinProcesar := search.GetUsuarios(direccionesSinProcesar)
	for _, val := range usuariosSinProcesar {
		usuarios = append(usuarios, val)
	}
	errores = insertar(db, "INSERT INTO usuario(ID, nombre, apellidos, DNI, correo, numero, letra,piso, portal, nombreCalle, nombreMunicipio) VALUES (?,?,?,?,?,?,?,?,?,?,?)", usuarios)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar usuario: ", e)
		}
	}

	// Insertar autoriza a
	var autoriza []models.Insercion
	for _, val := range search.GetAutorizaA(usuariosSinProcesar) {
		autoriza = append(autoriza, val)
	}
	errores = insertar(db, "INSERT INTO autorizaA(IDEmisor,IDReceptor) VALUES (?,?)", autoriza)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar autorizacion: ", e)
		}
	}

	// Inserta areas de envio
	var areasDeEnvio []models.Insercion
	areasSinProcesar := search.GetAreasDeEnvio()
	for _, val := range areasSinProcesar {
		areasDeEnvio = append(areasDeEnvio, val)
	}
	errores = insertar(db, "INSERT INTO areaDeEnvio(ID, codigoOficina) VALUES (?,?)", areasDeEnvio)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar Area de envio: ", e)
		}
	}

	// Insertar segmentos
	var segmentos []models.Insercion
	for _, val := range search.GetSegmentos(callesSinProcesar, rutasPreSinProcesar, areasSinProcesar) {
		segmentos = append(segmentos, val)
	}
	errores = insertar(db, "INSERT INTO segmento(nInicio, nFinal, nombreCalle, nombreMunicipio,IDAreaDeEnvio) VALUES (?,?,?,?,?)", segmentos)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar segmento: ", e)
		}
	}

	for _, area := range areasSinProcesar {
		if rand.Intn(2) == 1 {
			areaId := strings.Split(area.ID, "-")
			var numero string
			if i := rand.Intn(4); i > 9 {
				numero = strconv.Itoa(i)
				if numero == areaId[2] {
					n, _ := strconv.Atoi(numero)
					numero = strconv.Itoa(n + 1)
				}
			} else {
				numero = "0" + strconv.Itoa(i)
				if numero == areaId[2] {
					n, _ := strconv.Atoi(numero)
					numero = strconv.Itoa(n + 1)
				}
			}
			sqlSent := "UPDATE areaDeEnvio SET IDPadre = 'AR-" + areaId[1] + "-" + numero + "' WHERE ID = '" + area.ID + "' ;"
			sentenciaPreparada, err := db.Prepare(sqlSent)
			if err != nil {
				fmt.Println("Error al actulizar area de envío: ", err)
			}

			_, err = sentenciaPreparada.Exec()
			if err != nil {
				fmt.Println("Error al actulizar area de envío: ", err)
			}

			defer sentenciaPreparada.Close()
		}
	}

	// Insertar asosiasiones
	var asosiasiones []models.Insercion
	for _, val := range search.GetAsosianA(carterosSinProcesar, areasSinProcesar) {
		asosiasiones = append(asosiasiones, val)
	}
	errores = insertar(db, "INSERT INTO asocianA(IDAreaDeEnvio,DNICartero) VALUES (?,?)", asosiasiones)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar asosiación: ", e)
		}
	}

	// Insertar recogidas
	var recogidas []models.Insercion
	for _, val := range search.GetRecogidas(carterosSinProcesar, direccionesSinProcesar) {
		recogidas = append(recogidas, val)
	}
	errores = insertar(db, "INSERT INTO recogida(identificador,fecha, DNICartero, numero,letra,piso,portal,nombreCalle,nombreMunicipio) VALUES (?,?,?,?,?,?,?,?,?)", recogidas)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar recogida: ", e)
		}
	}

	// Insertar centro de clasificación
	var centros []models.Insercion
	centrosSinProcesar := search.GetCentrosDeClasificacion(municipiosSinProcesar)
	for _, val := range centrosSinProcesar {
		centros = append(centros, val)
	}
	errores = insertar(db, "INSERT INTO centroDeClasificacion(codigo,nombreMunicipio,nombre,nEnviosMax) VALUES (?,?,?,?)", centros)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar centro de clasificaicón: ", e)
		}
	}

	// Insertar oficinas
	var oficinas []models.Insercion
	oficinasSP := search.GetOficinas(direccionesSinProcesar, centrosSinProcesar)
	for _, val := range oficinasSP {
		oficinas = append(oficinas, val)
	}
	errores = insertar(db, "INSERT INTO oficina(codigo,codigoCentroClasificacion,numero,letra,piso,portal,nombreCalle,nombreMunicipio) VALUES (?,?,?,?,?,?,?,?)", oficinas)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar oficina: ", e)
		}
	}

	// Insertar coches
	var coches []models.Insercion
	cochesSP := search.GetCoches(oficinasSP)
	for _, val := range cochesSP {
		coches = append(coches, val)
	}
	errores = insertar(db, "INSERT INTO coche(matricula,codigoOficina,capacidad) VALUES (?,?,?)", coches)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar coche: ", e)
		}
	}

	// Insertar trabajaEn
	var trabajaEn []models.Insercion
	for _, val := range search.GetTrabajaEn(carterosSinProcesar, oficinasSP) {
		trabajaEn = append(trabajaEn, val)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()

		errores = insertar(db, "INSERT INTO trabajaEn(DNICartero,codigoOficina,horario,fechaInicio,turno) VALUES (?,?,?,?,?)", trabajaEn)
		if len(errores) != 0 {
			for _, e := range errores {
				fmt.Println("Error al insertar trabajaEn: ", e)
			}
		}
	}()

	// Insertar repartos
	var repartos []models.Insercion
	repartosSP := search.GetReparto(cochesSP, rutasPreSinProcesar, carterosSinProcesar, oficinasSP)
	for _, val := range repartosSP {
		repartos = append(repartos, val)
	}
	errores = insertar(db, "INSERT INTO reparto(IDRutaPredefinida,ID,matriculaCoche,codigoOficina,DNICartero) VALUES (?,?,?,?,?)", repartos)
	if len(errores) != 0 {
		for _, e := range errores {
			fmt.Println("Error al insertar  reparto: ", e)
		}
	}

	// Insertar cartasCer
	wg.Add(1)

	go func() {
		defer wg.Done()

		var cartasCer []models.Insercion
		for _, val := range search.GetCartasCertificadas(repartosSP, usuariosSinProcesar) {
			cartasCer = append(cartasCer, val)
		}
		errores = insertar(db, "INSERT INTO cartaCertificada(Identificador,fecha,IDReparto,urgencia,IDEmisor,IDReceptor) VALUES (?,?,?,?,?,?)", cartasCer)
		if len(errores) != 0 {
			for _, e := range errores {
				fmt.Println("Error al insertar carta certificada: ", e)
			}
		}
	}()

	// Insertar cartas
	wg.Add(1)

	go func() {
		defer wg.Done()

		var cartas []models.Insercion
		for _, val := range search.GetCartas(repartosSP, usuariosSinProcesar) {
			cartas = append(cartas, val)
		}
		errores = insertar(db, "INSERT INTO carta(Identificador,fecha,IDReparto,formato,IDEmisor,IDReceptor) VALUES (?,?,?,?,?,?)", cartas)
		if len(errores) != 0 {
			for _, e := range errores {
				fmt.Println("Error al insertar  carta : ", e)
			}
		}
	}()

	// Insertar paquete

	wg.Add(1)
	go func() {
		defer wg.Done()

		var paquete []models.Insercion
		for _, val := range search.GetPaquetes(repartosSP, usuariosSinProcesar) {
			paquete = append(paquete, val)
		}
		errores = insertar(db, "INSERT INTO paquete(Identificador,fecha,IDReparto,comentario, peso, dimensiones,IDEmisor,IDReceptor) VALUES (?,?,?,?,?,?,?,?)", paquete)
		if len(errores) != 0 {
			for _, e := range errores {
				fmt.Println("Error al insertar  carta : ", e)
			}
		}
	}()

	wg.Wait()
}
