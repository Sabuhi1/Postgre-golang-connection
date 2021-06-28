package psqlconn

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

type Config map[string]string

func Trim(s string) string {
	s = strings.TrimPrefix(s, `"`)
	s = strings.TrimSuffix(s, `"`)
	return s
}

func ReadConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{}

	if len(filename) == 0 {
		return config, nil
	}

	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		// check if the line has = sign
		// and process the line. Ignore the rest.
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				// assign the config map
				config[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		check(err)
	}
	return config, nil
}

func DBconn(dbfilepath string) {
	config, err := ReadConfig(dbfilepath)
	check(err)
	// Host Trimming
	host := config["host"]
	host = Trim(host)
	//Port Str to Int
	port_str := config["port"]
	number, err := strconv.ParseUint(port_str, 10, 32)
	port := int(number)
	//user trimming
	user := config["user"]
	user = Trim(user)
	//Password Trimming
	password := config["password"]
	password = Trim(password)
	//Dbname Trimming
	dbname := config["dbname"]
	dbname = Trim(dbname)

	//SQL Kismi
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	fmt.Println("Successfully connected!")
}
