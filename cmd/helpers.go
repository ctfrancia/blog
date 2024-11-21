package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/ctfrancia/blog/internal/model"
	"github.com/ctfrancia/blog/pkg/post"
)

// The serverError helper writes an error message and stack trace to the errorLog
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Used for later when making a bad request
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// not found error handler
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) renderIndex(w http.ResponseWriter, _ *http.Request, idxData model.IndexData) {
	ts, ok := app.templateCache["index.gohtml"]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", "index.gohtml"))
		return
	}
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, idxData)
	if err != nil {
		fmt.Println(err)
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) render(w http.ResponseWriter, _ *http.Request, name string, td post.PostData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
	}

	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *post.PostData, _ *http.Request) *post.PostData {
	if td == nil {
		td = &post.PostData{}
	}

	return td
}
