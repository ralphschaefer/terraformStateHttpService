package storeage

import (
	"errors"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"path"
)

const (
	Lockfile = "lock"
	Datafile = "state.json"
)

type FileStorageBuilder struct {
	Directory string
}

func (sb* FileStorageBuilder) Build(project string) (StorageInterface, error) {
	return InitFileStorage(sb.Directory, project)
}

type FileStorage struct {
	directory string
	project string
	projectDir string
}

func createDir(dir string) error {
	finfo, err := os.Stat(dir)
	if err == nil {
		if finfo.IsDir() {
			return nil
		} else {
			return errors.New("not directory")
		}
	}
	return os.MkdirAll(dir, 0777)
}

func exixts(dir string) bool {
	_, err := os.Stat(path.Join(dir))
	return !os.IsNotExist(err)
}

func InitFileStorage(dir string, project string) (*FileStorage, error) {
	s := FileStorage{
		directory: dir,
		project: project,
		projectDir: path.Join(dir, project),
	}
	err := createDir(s.projectDir)
	if (err != nil) {
		return nil, err
	} else {
		return &s, nil
	}
}

func (s* FileStorage) IsLocked() (bool, *string) {
	filename := path.Join(s.projectDir, Lockfile)
	if exixts(filename) {
		logid, err := ioutil.ReadFile(filename)
		if err != nil {
			glog.Fatal(err)
		}
		logidstr := string(logid)
		return true, &logidstr
	} else {
		return false, nil
	}
}

func (s* FileStorage) Lock(id string) bool {
	isLocked, _ := s.IsLocked()
	if isLocked {
		return false
	}
	f, err := os.Create(path.Join(s.projectDir, Lockfile))
	defer f.Close()
	if err != nil {
		glog.Fatal(err)
	}
	_, err = f.Write([]byte(id))
   	if err != nil {
      	glog.Fatal(err)
   	}
	return true
}

func (s* FileStorage) Unlock(id string) bool{
	isLocked, logId := s.IsLocked()
	if (isLocked) && (*logId == id) {
		file := path.Join(s.projectDir, Lockfile)
		err := os.Remove(file)
		if err != nil {
			glog.Fatal(err)
		}
		return true
	} else {
		return false
	}
}


func (s* FileStorage) Delete() {
	if (exixts(path.Join(s.projectDir, Lockfile))) {
		err := os.Remove(path.Join(s.projectDir, Lockfile))
		if err != nil {
			glog.Fatal(err)
		}
	}
	if (exixts(path.Join(s.projectDir, Datafile))) {
		err := os.Remove(path.Join(s.projectDir, Datafile))
		if err != nil {
			glog.Fatal(err)
		}
	}
	err := os.Remove(s.projectDir)
	if err != nil {
		glog.Fatal(err)
	}
}

func (s* FileStorage) Put(id string, content []byte) bool {
	file, err := os.OpenFile(
		path.Join(s.projectDir, Datafile),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()
	if err != nil {
	  	glog.Fatal(err)
	}
	isLocked, lockId := s.IsLocked()
	if isLocked && (*lockId == id) {
		_, err = file.Write(content)
		if err != nil {
			glog.Fatal(err)
		}
		return true
	} else {
		return false
	}

}

func (s* FileStorage) Get() []byte {
	filename := path.Join(s.projectDir, Datafile)
	if exixts(filename) {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			glog.Fatal(err)
		}
		return data
	} else {
		return make([]byte, 0)
	}
}