package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gernest/utron"
	"github.com/gernest/utron/controller"
	"github.com/gernest/utron/router"
)

type indexView struct {
	t *template.Template
}

func newIndexView(file string) (*indexView, error) {
	t, err := template.ParseFiles(file)
	if err != nil {
		return nil, err
	}
	return &indexView{t: t}, nil
}

func (v indexView) Render(out io.Writer, name string, data interface{}) error {
	return v.t.ExecuteTemplate(out, name, data)
}

type Upload struct {
	controller.BaseController
	Routes []string
}

func (u *Upload) Index() {
	u.Ctx.Template = "index.html"
}

//Save saves the uploaded file
func (u *Upload) Save() {
	r := u.Ctx.Request()
	err := r.ParseMultipartForm(15485760)
	if err != nil {
		fmt.Println(err)
		return
	}
	f, h, err := r.FormFile("upload")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	o, err := os.OpenFile(h.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer o.Close()
	io.Copy(o, f)
	fmt.Println("OK")
}

func NewUplad() controller.Controller {
	return &Upload{
		Routes: []string{
			"get;/;Index",
			"post;/upload;Save",
			"get;/delete/{id};Delete",
		},
	}
}

func main() {
	app := utron.NewApp()
	v, err := newIndexView("index.html")
	if err != nil {
		app.Log.Errors(err)
		return
	}
	app.View = v
	app.Router.Options = &router.Options{
		View: v,
	}
	app.Router.Add(NewUplad)
	port := ":8090"
	app.Log.Info("staring server on port", port)
	log.Fatal(http.ListenAndServe(port, app))

}
