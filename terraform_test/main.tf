

terraform {
  backend "http" {
    address = "http://127.0.0.1:8080/test"
    lock_address = "http://127.0.0.1:8080/test"
    unlock_address = "http://127.0.0.1:8080/test"
  }
}


provider "local" {
}

provider "random" {
}


resource "random_string" "test" {
  length  = 20
  special = false
  upper   = false
}

resource "random_uuid" "testname" {
}

resource "local_file" "test" {
  filename = join(".", ["testfile", random_uuid.testname.result])
  content = random_string.test.result
}

output "name" {
  value = local_file.test.filename
}