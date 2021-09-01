# Http Terraform Backend

REST Service wich stores terraform states and supports locking 
written in GO

## Usage

### startup standalone
```
./stageservice -port <PORT> -bind <ADDRESS> -storeto <DIR>
```

* PORT: Portnumber the http service listens to
* ADDRESS: IP address the http service binds to
* DIR: Directory where states and lock are written to

### startup in docker
```
docker run -d stateservice:0.1.0
```

### use in terraform

Example: 
```terraform
terraform {
  backend "http" {
    address = "http://127.0.0.1:8080/test"
    lock_address = "http://127.0.0.1:8080/test"
    unlock_address = "http://127.0.0.1:8080/test"
  }
}
```
* _127.0.0.1_: is the host the __stageservice__ is runing on
* _8080_: is the port the __stageservice__ listens
* _test_: is an arbitrary terraform project name

The __stageservice__ requires locking, therefore the URLs in the 
terraform backend for ```address```, ```lock_address``` and 
```unlock_address``` must not differ.

## Building
* for building executable run ```make build```
* for creating a docker container ```make docker```

