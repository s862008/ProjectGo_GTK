package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
	_ "github.com/nakagami/firebirdsql"
)

type Config struct {
	baseName string
	basePath string
	basePSW  string
	baseUser string
}

func main() {
	// Читаем натройки из файла
	var cfg Config
	getConfig(&cfg)
	fmt.Println(cfg)

	// Инициализируем GTK.
	gtk.Init(nil)
	fmt.Println("Запуск")
	// Создаём билдер
	bld, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Загружаем в билдер окно из файла Glade
	err = bld.AddFromFile("mainForm.glade")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Получаем объект главного окна по ID
	obj, err := bld.GetObject("dialog1")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Преобразуем из объекта именно окно типа gtk.Window
	// и соединяем с сигналом "destroy" чтобы можно было закрыть
	// приложение при закрытии окна
	dlg := obj.(*gtk.Dialog)
	dlg.Connect("destroy", func() {
		gtk.MainQuit()
	})
	actions(bld)
	// Отображаем все виджеты в окне
	dlg.ShowAll()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()
}
func getConfig(cfg *Config) {
	file, err := os.Open("config.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2
	reader.Comment = '#'
	reader.Comma = ';'
	for {
		record, e := reader.Read()
		if e != nil {
			fmt.Println(e) //EOF
			break
		}
		fmt.Println(record)
		switch record[0] {
		case "baseName":
			cfg.baseName = record[1]
		case "basePath":
			cfg.basePath = record[1]
		case "basePSW":
			cfg.basePSW = record[1]
		case "baseUser":
			cfg.baseUser = record[1]
		}
	}
}

func actions(bld *gtk.Builder) {
	obj, _ := bld.GetObject("textview1")
	textview1 := obj.(*gtk.TextView)
	obj, _ = bld.GetObject("button1")
	button1 := obj.(*gtk.Button)
	obj, _ = bld.GetObject("button2")
	button2 := obj.(*gtk.Button)

	button1.Connect("clicked", func() {

		var n int
		conn, _ := sql.Open("firebirdsql", "SYSDBA:masterkey@localhost")
		defer conn.Close()
		conn.QueryRow("SELECT Count(*) FROM rdb$relations").Scan(&n)
		fmt.Println("Relations count=", n)

		t_buff, _ := textview1.GetBuffer()
		t_buff.SetText(strconv.Itoa(n))

	})

	button2.Connect("clicked", func() {
		text := "ТЕСТ2 ТЕСТ2"
		t_buff, _ := textview1.GetBuffer()
		start, end := t_buff.GetBounds()
		textold, _ := t_buff.GetText(start, end, true)
		t_buff.SetText(textold + text + "\n")

	})

}
