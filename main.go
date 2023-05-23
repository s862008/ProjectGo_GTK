package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	_ "github.com/nakagami/firebirdsql"
)

type Config struct {
	baseName string
	basePath string
	basePSW  string
	baseUser string
}

var cfg Config

func main() {
	// Читаем натройки из файла

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
	//objlst, _ := bld.GetObject("liststore1")
	//	liststore1 := objlst.(*gtk.ListStore)

	objtree, _ := bld.GetObject("treeview1")
	treeview1 := objtree.(*gtk.TreeView)
	{
		renderer, _ := gtk.CellRendererTextNew()
		column1, _ := gtk.TreeViewColumnNewWithAttribute("ID", renderer, "text", 0)
		column2, _ := gtk.TreeViewColumnNewWithAttribute("TASK", renderer, "text", 1)
		treeview1.AppendColumn(column1)
		treeview1.AppendColumn(column2)
	}

	liststore1, _ := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_BOOLEAN)
	//liststore1.SetValue(liststore1.Append(), 0, "wwwww")
	//liststore1.SetValue(liststore1.Append(), 1, "ww234w")

	// Добавляем данные в список.
	data := [][]interface{}{
		{"Элемент 1", "a", true},
		{"Элемент 2", "q", false},
		{"Элемент 3", "test", true},
	}

	for _, val := range data {
		iter := liststore1.Append()
		liststore1.SetValue(iter, 0, val[0])
		liststore1.SetValue(iter, 1, val[1])
		liststore1.SetValue(iter, 2, val[2])
	}

	// Создаем CellRendererToggle для второго столбца (чекбокса).
	rendererToggle, _ := gtk.CellRendererToggleNew()
	col2, _ := gtk.TreeViewColumnNewWithAttribute("Столбец 3", rendererToggle, "active", 2)

	rendererToggle.Connect("toggled", func(renderer *gtk.CellRendererToggle, path *gtk.TreePath) {
		// Обработчик события "toggled" для CellRendererToggle при отметке/снятии отметки на чекбоксе.
		iter, _ := liststore1.GetIter(path)
		act, _ := liststore1.GetValue(iter, 2)
		fmt.Println(act.IsValue())
		// if act.IsValue() {
		// 	liststore1.SetValue(iter, 1, false)
		// } else {
		// 	liststore1.SetValue(iter, 1, true)
		// }
	})

	treeview1.SetModel(liststore1)
	treeview1.AppendColumn(col2)

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
	button1.SetLabel("Получить данные")
	obj, _ = bld.GetObject("button2")
	button2 := obj.(*gtk.Button)
	button2.SetLabel("Выход")

	button1.Connect("clicked", func() {

		var n int
		conn, _ := sql.Open("firebirdsql", cfg.baseUser+":"+cfg.basePSW+"@"+cfg.basePath)
		defer conn.Close()
		conn.QueryRow("SELECT Count(*) FROM rdb$relations").Scan(&n)
		fmt.Println("Relations count=", n)

		t_buff, _ := textview1.GetBuffer()
		t_buff.SetText(strconv.Itoa(n))

	})

	button2.Connect("clicked", func() {

		gtk.MainQuit()

	})

}
