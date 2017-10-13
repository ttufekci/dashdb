CREATE TABLE `anothertesttable` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `age` varchar(45) DEFAULT NULL,
  `code` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

CREATE TABLE `primkeynamecustom` (
  `idprimkeynamecustom` int(11) NOT NULL,
  `primkeynamecustomcol` varchar(45) DEFAULT NULL,
  `primkeynamecustomcol1` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`idprimkeynamecustom`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `testtbl` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;