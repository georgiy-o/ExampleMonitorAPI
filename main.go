package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var DB *sql.DB

func main() {
	if len(os.Args) > 1 {
		switch command := strings.ToLower(os.Args[1]); {
		case command == "--help":
			printHelp()
		case command == "--createdb":
			if _, err := os.Stat("./monitors.txt"); os.IsNotExist(err) {
				fmt.Println("ERROR! File \"monitors.txt\" does not exist!")
				return
			}
			if _, err := os.Stat("./products.db"); err == nil {
				err = os.Remove("./products.db")
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			CreateDB()
			AddMonitorsFromFile("./monitors.txt")
			fmt.Println("OK. File products.db is created!")
			return

		case command == "--start":
			http.HandleFunc("/category/monitors", GetMonitors)
			http.HandleFunc("/category/monitor/", GetStatForMonitor)
			http.HandleFunc("/category/monitor_click/", AddClickForMonitor)
			fmt.Println("The server is running!")
			fmt.Println("Looking forward to requests...")
			if err := http.ListenAndServe(":8030", nil); err != nil {
				log.Fatal("Failed to start server!", err)
			}
		default:
			printHelp()
		}

	}
}

func CreateDB() {
	OpenDB()
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS monitors (id integer, name varchar(255) not null, count integer)`)
	if err != nil {
		log.Fatal(err)
	}
	DB.Close()
}

func OpenDB() {
	db, err := sql.Open("sqlite3", "products.db")
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}

func AddMonitorsFromFile(filename string) {
	var file *os.File
	var err error
	if file, err = os.Open(filename); err != nil {
		log.Fatal("ERROR. Failed to open file: ", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	OpenDB()
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), ",")
		id := arr[0]
		monitorName := arr[1]

		_, err = DB.Exec("INSERT INTO monitors (id, name, count) VALUES ($1, $2, 0)", id, monitorName)
	}
}

func AddClickForMonitor(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	setCORSHeaders(writer)
	if err != nil {
		fmt.Fprintf(writer, "{%s}", err)
	} else {
		monitorId := strings.TrimPrefix(request.URL.Path, "/category/monitor_click/")
		OpenDB()
		countValue := 0
		rows, _ := DB.Query("select count from monitors where id=" + monitorId)
		for rows.Next() {
			rows.Scan(&countValue)
		}
		countValue++
		_, _ = DB.Exec("update monitors set count=" + strconv.Itoa(countValue) + " where id=" + monitorId)
	}
}

func GetStatForMonitor(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	setCORSHeaders(writer)
	if err != nil {
		fmt.Fprintf(writer, "{%s}", err)
	} else {
		countValue := 0
		monitorId := strings.TrimPrefix(request.URL.Path, "/category/monitor/")
		OpenDB()
		rows, _ := DB.Query("select count from monitors where id=$1", monitorId)
		for rows.Next() {
			rows.Scan(&countValue)
		}
		strOut := "{ \"id\": \"" + monitorId + "\", \"count\": \"" + strconv.Itoa(countValue) + "\"}"
		fmt.Fprintf(writer, strOut)
	}
}
func GetMonitors(writer http.ResponseWriter, request *http.Request) {
	OpenDB()
	monitors := GetFromDBNameModel("monitors")
	err := request.ParseForm()
	setCORSHeaders(writer)
	if err != nil {
		fmt.Fprintf(writer, "{%s}", err)
	} else {
		strOut := "{ \"monitors\": ["
		for i := 0; i < len(monitors)-1; i++ {
			id := monitors[i][0]
			name := monitors[i][1]
			strOut += fmt.Sprintf("[%v, %v], ", id, name)
		}
		strOut += fmt.Sprintf("[%v, %v]] } ", monitors[len(monitors)-1][0], monitors[len(monitors)-1][1])
		fmt.Fprintf(writer, strOut)
	}
}

func GetFromDBNameModel(tblName string) [][]interface{} {
	var arr [][]interface{}
	var id int
	var name string

	rows, _ := DB.Query("SELECT id, name FROM " + tblName)
	for rows.Next() {
		rows.Scan(&id, &name)
		arr = append(arr, []interface{}{id, name})
	}
	return arr
}

func printHelp() {
	fmt.Println()
	fmt.Println("Help:                        ./counter --help")
	fmt.Println("Create products database:    ./counter --createdb")
	fmt.Println("Start server:                ./counter --start")
}

func setCORSHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
}
