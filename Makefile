# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

VERSION :=
TAG := $(shell git describe --abbrev=0 --tags HEAD 2>/dev/null)
COMMIT := $(shell git rev-parse HEAD)
ifeq ($(TAG),)
    VERSION := latest
else
    ifeq ($(COMMIT), $(shell git rev-list -n1 $(TAG)))
        VERSION := $(TAG)
    else
        VERSION := $(TAG)-$(COMMIT)
    endif
endif

container: build s3-container quick-container
.PHONY: container

clean:
	rm -f s3fs-container
.PHONY: clean

s3-container:
        docker build -t s3fs:latest mount/ 
.PHONY: s3-container

quick-container:
	docker build -t childsb/s3fs-provisioner:latest . 
.PHONY: quick-container

all build:
	go build 
.PHONY: all build
