CREATE SCHEMA IF NOT EXISTS correos 
DEFAULT CHARACTER SET utf8 
COLLATE utf8_spanish2_ci;


USE correos;


CREATE TABLE cartero (
  DNI VARCHAR(9),
  nombre VARCHAR(250) NOT NULL,
  apellidos VARCHAR(250) NOT NULL,
  PRIMARY KEY (DNI)
);

CREATE TABLE municipio (
  nombre VARCHAR(250) ,
  provincia VARCHAR(250) NOT NULL,
  PRIMARY KEY (nombre)
);

CREATE TABLE calle (
  nombre VARCHAR(250),
  nombreMunicipio VARCHAR(250),
  PRIMARY KEY(nombre, nombreMunicipio),
  CONSTRAINT FOREIGN KEY (nombreMunicipio) REFERENCES municipio (nombre)
	ON UPDATE CASCADE
  ON DELETE CASCADE
);

CREATE TABLE direccion (
  numero INT,
  letra CHAR,
  piso INT,
  portal INT, 
  nombreCalle VARCHAR(250),
  nombreMunicipio VARCHAR(250),
  PRIMARY KEY (numero, letra, piso, portal, nombreCalle, nombreMunicipio),
  CONSTRAINT 
	FOREIGN KEY (nombreCalle, nombreMunicipio) REFERENCES calle (nombre, nombreMunicipio)
	ON UPDATE CASCADE
	ON DELETE CASCADE
);

CREATE TABLE rutaPredefinida (
  ID INT,
  PRIMARY KEY (ID)
);

CREATE TABLE segmento (
  nInicio INT,
  nFinal INT, 
  nombreCalle VARCHAR(250),
  nombreMunicipio VARCHAR(250),
  IDRutaPredefinida INT DEFAULT NULL,
  orden INT DEFAULT NULL,
  PRIMARY KEY (nInicio, nFinal, nombreCalle, nombreMunicipio),
  CONSTRAINT 
	FOREIGN KEY (nombreCalle, nombreMunicipio) REFERENCES calle (nombre, nombreMunicipio)
		ON UPDATE CASCADE
    ON DELETE CASCADE,
	FOREIGN KEY (IDRutaPredefinida) REFERENCES rutaPredefinida (ID)
		ON UPDATE CASCADE
    ON DELETE CASCADE
);

CREATE TABLE usuario (
  ID INT,
  nombre VARCHAR(250) NOT NULL,
  apellidos VARCHAR(250) NOT NULL,
  DNI VARCHAR(9) DEFAULT NULL,
  correo VARCHAR(250) DEFAULT NULL,
  numero INT DEFAULT NULL,
  letra CHAR DEFAULT NULL,
  piso INT DEFAULT NULL,
  portal INT DEFAULT NULL, 
  nombreCalle VARCHAR(250) DEFAULT NULL,
  nombreMunicipio VARCHAR(250) DEFAULT NULL,
  PRIMARY KEY (ID),
  CONSTRAINT 
	FOREIGN KEY (numero, letra, piso, portal, nombreCalle, nombreMunicipio) REFERENCES direccion (numero, letra, piso, portal, nombreCalle, nombreMunicipio)
		ON UPDATE CASCADE
    ON DELETE SET NULL
);

CREATE TABLE autorizaA (
  IDEmisor INT,
  IDReceptor INT,
  PRIMARY KEY (IDEmisor, IDReceptor),
  CONSTRAINT 
	FOREIGN KEY (IDEmisor) REFERENCES usuario (ID)
		ON UPDATE CASCADE
    ON DELETE CASCADE,
  FOREIGN KEY (IDReceptor) REFERENCES usuario(ID)
		ON UPDATE CASCADE
    ON DELETE CASCADE
);

CREATE TABLE areaDeEnvio (
  ID VARCHAR(250),
  IDPadre VARCHAR(250) DEFAULT NULL,
  PRIMARY KEY (ID),
  CONSTRAINT
	FOREIGN KEY(IDPadre) REFERENCES areaDeEnvio(ID)
	ON UPDATE CASCADE

);

CREATE TABLE asocianA (
  IDAreaDeEnvio VARCHAR(250),
  DNICartero VARCHAR(9),
  PRIMARY KEY (IDAreaDeEnvio, DNICartero),
  CONSTRAINT 
	FOREIGN KEY (IDAreaDeEnvio) REFERENCES areaDeEnvio (ID)
		ON UPDATE CASCADE
    ON DELETE CASCADE,
	FOREIGN KEY (DNICartero) REFERENCES cartero (DNI)
		ON UPDATE CASCADE
    ON DELETE CASCADE
);

CREATE TABLE recogida (
  identificador INT,
  fecha DATE DEFAULT NULL,
  DNICartero VARCHAR(9) NOT NULL,
  numero INT DEFAULT NULL,
  letra CHAR DEFAULT NULL,
  piso INT DEFAULT NULL,
  portal INT DEFAULT NULL, 
  nombreCalle VARCHAR(250) DEFAULT NULL,
  nombreMunicipio VARCHAR(250) DEFAULT NULL,
  PRIMARY KEY (identificador),
  CONSTRAINT
	FOREIGN KEY (DNICartero) REFERENCES cartero (DNI)
		ON UPDATE CASCADE,
    FOREIGN KEY (numero, letra, piso, portal, nombreCalle, nombreMunicipio) REFERENCES direccion (numero, letra, piso, portal, nombreCalle, nombreMunicipio)
		ON UPDATE CASCADE
    ON DELETE SET NULL
);

CREATE TABLE centroDeClasificacion (
  codigo INT,
  nombreMunicipio VARCHAR(250) DEFAULT NULL,
  nombre VARCHAR(250) DEFAULT NULL,
  nEnviosMax INT DEFAULT NULL,
  PRIMARY KEY (codigo),
  CONSTRAINT 
	FOREIGN KEY (nombreMunicipio) REFERENCES municipio (nombre)
		ON UPDATE CASCADE
    ON DELETE SET NULL
);

CREATE TABLE oficina (
  codigo VARCHAR(250),
  codigoCentroClasificacion INT NOT NULL,
  numero INT,
  letra CHAR,
  piso INT,
  portal INT,
  nombreCalle VARCHAR(250), 
  nombreMunicipio VARCHAR(250),
  PRIMARY KEY (codigo),
  CONSTRAINT 
	FOREIGN KEY (numero, letra, piso, portal, nombreCalle, nombreMunicipio) REFERENCES direccion(numero , letra, piso, portal, nombreCalle, nombreMunicipio)
		ON UPDATE CASCADE
        ON DELETE SET NULL,
	FOREIGN KEY (codigoCentroClasificacion) REFERENCES centroDeClasificacion (codigo)
		ON UPDATE CASCADE
);

CREATE TABLE coche (
  matricula VARCHAR(7),
  codigoOficina VARCHAR(250) DEFAULT NULL,
  capacidad INT DEFAULT NULL,
  PRIMARY KEY (matricula),
  CONSTRAINT 
  FOREIGN KEY (codigoOficina) REFERENCES oficina (codigo)
	ON UPDATE CASCADE
    ON DELETE SET NULL
);

CREATE TABLE trabajaEn (
  DNICartero VARCHAR(9),
  codigoOficina VARCHAR(250),
  horario VARCHAR (50) DEFAULT NULL,
  fechaInicio DATE NOT NULL,
  turno VARCHAR(250) DEFAULT NULL,
  PRIMARY KEY (DNICartero, codigoOficina),
  CONSTRAINT 
	FOREIGN KEY (DNICartero) REFERENCES cartero (DNI)
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	FOREIGN KEY (codigoOficina) REFERENCES oficina (codigo)
		ON UPDATE CASCADE
    ON DELETE CASCADE
);

CREATE TABLE reparto (
  IDRutaPredefinida INT NOT NULL,
  ID INT,
  matriculaCoche VARCHAR(7) NOT NULL,
  codigoOficina VARCHAR(250) NOT NULL,
  DNICartero VARCHAR(9) NOT NULL,
  PRIMARY KEY (ID),
  CONSTRAINT 
	FOREIGN KEY (IDRutaPredefinida) REFERENCES rutaPredefinida (ID)
		ON UPDATE CASCADE,
    FOREIGN KEY (matriculaCoche) REFERENCES coche (matricula)
		ON UPDATE CASCADE,
	FOREIGN KEY (codigoOficina) REFERENCES oficina (codigo)
		ON UPDATE CASCADE,
    	FOREIGN KEY (DNICartero) REFERENCES cartero (DNI)
		ON UPDATE CASCADE
);

CREATE TABLE cartaCertificada (
  identificador VARCHAR(12),
  fecha DATETIME,
  IDReparto INT DEFAULT NULL,
  urgencia INT NOT NULL,
  IDEmisor INT NOT NULL,
  IDReceptor INT NOT NULL,
  PRIMARY KEY (identificador),
  CONSTRAINT 
	FOREIGN KEY (IDReparto) REFERENCES reparto (ID)
		ON UPDATE CASCADE
        ON DELETE SET NULL,
	FOREIGN KEY (IDEmisor) REFERENCES usuario (ID)
		ON UPDATE CASCADE,
	FOREIGN KEY (IDReceptor) REFERENCES usuario(ID)
		ON UPDATE CASCADE
);

CREATE TABLE carta (
  identificador VARCHAR(12),
  fecha DATETIME,
  IDReparto INT DEFAULT NULL,
  formato VARCHAR (2) DEFAULT NULL,
  IDEmisor INT NOT NULL,
  IDReceptor INT NOT NULL,
  PRIMARY KEY (identificador),
  CONSTRAINT 
	FOREIGN KEY (IDReparto) REFERENCES reparto (ID)
		ON UPDATE CASCADE
        ON DELETE SET NULL,
	FOREIGN KEY (IDEmisor) REFERENCES usuario (ID)
		ON UPDATE CASCADE,
	FOREIGN KEY (IDReceptor) REFERENCES usuario(ID)
		ON UPDATE CASCADE
);

CREATE TABLE paquete (
  identificador VARCHAR(12),
  fecha DATETIME,
  IDReparto INT DEFAULT NULL,
  comentario VARCHAR (250) DEFAULT NULL,
  peso INT NOT NULL,
  dimensiones VARCHAR(250) NOT NULL,
  IDEmisor INT NOT NULL,
  IDReceptor INT NOT NULL,
  PRIMARY KEY (identificador),
  CONSTRAINT 
	FOREIGN KEY (IDReparto) REFERENCES reparto (ID)
		ON UPDATE CASCADE
        ON DELETE SET NULL,
	FOREIGN KEY (IDEmisor) REFERENCES usuario (ID)
		ON UPDATE CASCADE,
	FOREIGN KEY (IDReceptor) REFERENCES usuario(ID)
		ON UPDATE CASCADE
);

15:27:41	CREATE TABLE reparto (   IDRutaPredefinida VARCHAR(250) NOT NULL,   ID INT,   matriculaCoche VARCHAR(7) NOT NULL,   codigoOficina VARCHAR(250) NOT NULL,   DNICartero VARCHAR(9) NOT NULL,   PRIMARY KEY (ID),   CONSTRAINT   FOREIGN KEY (IDRutaPredefinida) REFERENCES rutaPredefinida (ID)   ON UPDATE CASCADE,     FOREIGN KEY (matriculaCoche) REFERENCES coche (matricula)   ON UPDATE CASCADE,  FOREIGN KEY (codigoOficina) REFERENCES oficina (codigo)   ON UPDATE CASCADE,      FOREIGN KEY (DNICartero) REFERENCES cartero (DNI)   ON UPDATE CASCADE )	Error Code: 3780. Referencing column 'IDRutaPredefinida' and referenced column 'ID' in foreign key constraint 'reparto_ibfk_1' are incompatible.	0.0020 sec

