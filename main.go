
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"vigo360.es/new/internal"
	"vigo360.es/new/internal/database"
)

var (
	version string
)

func main() {
	if err := checkEnv(); err != nil {
		fmt.Printf("<3>error validando entorno: %s\n", err.Error())
		os.Exit(1)
	}

	if err := run(); err != nil {
		fmt.Printf("<3>%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Printf("<6>iniciando vigo360 versión %s\n", version)
	var PORT string = ":" + os.Getenv("PORT")

	var db = database.GetDB()
	var container = internal.NewMysqlContainer(db)

	var s = internal.NewServer(container)

	fmt.Printf("<6>iniciando servidor web en %s\n", PORT)
	http.Handle("/", s.Router)

	var err = http.ListenAndServe(PORT, nil)
	return err
}

func checkEnv() error {
	if val, is := os.LookupEnv("PORT"); !is || val == "" {
		return fmt.Errorf("es necesario especificar PORT")
	} else {
		i, e := strconv.Atoi(val)
		if e != nil {
			return fmt.Errorf("PORT tiene que ser un número")
		}
		if i < 0 || i > 65535 {
			return fmt.Errorf("PORT debe ser un puerto TCP válido")
		}
	}

	if val, is := os.LookupEnv("UPLOAD_PATH"); !is || val == "" {
		return fmt.Errorf("es necesario especificar UPLOAD_PATH")
	} else {
		info, err := os.Stat(val)
		if err != nil {
			return fmt.Errorf("error comprobando validez de UPLOAD_PATH: %s", err.Error())
		}
		if !info.IsDir() {
			return fmt.Errorf("UPLOAD_PATH tiene que ser un directorio: %s", err.Error())
		}
		err = os.WriteFile(val+"/.test", []byte{0x00}, os.ModePerm)
		if err != nil {
			return fmt.Errorf("no se puede escribir a UPLOAD_PATH: %s", err.Error())
		}
		os.Remove(val + "/.test")
	}

	if val, is := os.LookupEnv("DOMAIN"); !is || val == "" {
		return fmt.Errorf("es necesario especificar DOMAIN")
	}

	return nil
}
