package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	// Инициализируем GTK.
	gtk.Init(nil)
	fmt.Println("Запуск")
	// Создаём билдер
	b, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Загружаем в билдер окно из файла Glade
	err = b.AddFromFile("mainForm.glade")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Получаем объект главного окна по ID
	obj, err := b.GetObject("dialog1")
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
	actions(b)
	// Отображаем все виджеты в окне
	dlg.ShowAll()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()
}

func actions(bld *gtk.Builder) {
	obj, _ := bld.GetObject("textview1")
	textview1 := obj.(*gtk.TextView)
	obj, _ = bld.GetObject("button1")
	button1 := obj.(*gtk.Button)
	obj, _ = bld.GetObject("button2")
	button2 := obj.(*gtk.Button)

	button1.Connect("clicked", func() {
		text := "ТЕСТ ТЕСТ"
		t_buff, _ := textview1.GetBuffer()
		t_buff.SetText(text)

	})

	button2.Connect("clicked", func() {
		text := "ТЕСТ2 ТЕСТ2"
		t_buff, _ := textview1.GetBuffer()
		start, end := t_buff.GetBounds()
		textold, _ := t_buff.GetText(start, end, true)
		t_buff.SetText(textold + text + "\n")

	})

}
