package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"io/ioutil"
	"main/storeage"
	"net/http"
)

type Server struct {
	port int
	bind string
	builder storeage.StorageBuilderInterface
}

type LockInfo struct {
	ID string
	Operation string
	Info string
	Who string
	Version string
	Created string
	Path string
}

func InitServer(port int, bind string, builder storeage.StorageBuilderInterface) *Server {
	s := Server{
		port: port,
		bind: bind,
		builder: builder,
	}
	return &s
}


func (s* Server) home(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info(r.URL)
	w.Write([]byte("terraform state server"))
}

func (s* Server) handler_get(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info(r.URL)
	vars := mux.Vars(r)
	prjName := vars["project"]
	glog.Info(fmt.Sprintf("get project: %s", prjName))
	prj, err := s.builder.Build(prjName)
	if (err != nil) {
		glog.Error(err)
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(prj.Get())
}

func (s* Server) handler_unlock(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info(r.URL)
	vars := mux.Vars(r)
	prjName := vars["project"]
	glog.Info(fmt.Sprintf("unlock project: %s", prjName))
	prj, err := s.builder.Build(prjName)
	if err != nil {
		glog.Error(err)
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	var lockInfo LockInfo
    err = json.NewDecoder(r.Body).Decode(&lockInfo)
	defer r.Body.Close()
	if err != nil {
		glog.Error(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	if !prj.Unlock(lockInfo.ID) {
		err = errors.New("cannot unlock")
		glog.Error(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
}

func (s* Server) handler_lock(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info(r.URL)
	vars := mux.Vars(r)
	prjName := vars["project"]
	glog.Info(fmt.Sprintf("lock project: %s", prjName))
	prj, err := s.builder.Build(prjName)
	if err != nil {
		glog.Error(err)
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	var lockInfo LockInfo
    err = json.NewDecoder(r.Body).Decode(&lockInfo)
	defer r.Body.Close()
	if err != nil {
		glog.Error(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	isLocked, _ := prj.IsLocked()
	if isLocked {
		w.WriteHeader(http.StatusLocked)
	} else {
		if prj.Lock(lockInfo.ID) {
			w.WriteHeader(200)
		} else {
			glog.Error(errors.New("cannot log"))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
}

func (s* Server) handler_post(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info(r.URL)
	vars := mux.Vars(r)
	prjName := vars["project"]
	lockId := r.URL.Query().Get("ID")
	glog.Info(fmt.Sprintf("put project: %s with lockID '%s'", prjName, lockId))
	if lockId == "" {
		err := errors.New("missing query param ID")
		glog.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	prj, err := s.builder.Build(prjName)
	if err != nil {
		glog.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if prj.Put(lockId, body) {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(http.StatusLocked)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (s* Server) handler_delete(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info(r.URL)
	vars := mux.Vars(r)
	prjName := vars["project"]
	glog.Info(fmt.Sprintf("delete project: %s", prjName))
	prj, err := s.builder.Build(prjName)
	if err != nil {
		glog.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	prj.Delete()
	w.Header().Set("Content-Type", "application/json")
}

func (s* Server) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.home).Methods("GET")
	r.HandleFunc("/{project}", s.handler_get).Methods( "GET")
	r.HandleFunc("/{project}", s.handler_lock).Methods( "LOCK")
	r.HandleFunc("/{project}", s.handler_unlock).Methods( "UNLOCK")
	r.HandleFunc("/{project}", s.handler_post).Methods( "POST")
	r.HandleFunc("/{project}", s.handler_delete).Methods( "DELETE")
	http.Handle("/", r)
	glog.Info(fmt.Sprintf("Listen to %s:%d",s.bind,s.port))
	http.ListenAndServe(fmt.Sprintf("%s:%d",s.bind,s.port), r)
}
