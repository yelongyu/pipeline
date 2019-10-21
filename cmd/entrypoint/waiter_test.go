/*
Copyright 2019 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestRealWaiterWaitMissingFile(t *testing.T) {
	// Create a temp file and then immediately delete it to get
	// a legitimate tmp path and ensure the file doesnt exist
	// prior to testing Wait().
	tmp, err := ioutil.TempFile("", "real_waiter_test_file")
	if err != nil {
		t.Errorf("error creating temp file: %v", err)
	}
	os.Remove(tmp.Name())
	rw := realWaiter{}
	doneCh := make(chan struct{})
	go func() {
		err := rw.Wait(tmp.Name(), false)
		if err != nil {
			t.Errorf("error waiting on tmp file %q", tmp.Name())
		}
		close(doneCh)
	}()
	select {
	case <-doneCh:
		t.Errorf("did not expect Wait() to have detected a file at path %q", tmp.Name())
	case <-time.After(2 * waitPollingInterval):
		// Success
	}
}

func TestRealWaiterWaitWithFile(t *testing.T) {
	tmp, err := ioutil.TempFile("", "real_waiter_test_file")
	if err != nil {
		t.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	rw := realWaiter{}
	doneCh := make(chan struct{})
	go func() {
		err := rw.Wait(tmp.Name(), false)
		if err != nil {
			t.Errorf("error waiting on tmp file %q", tmp.Name())
		}
		close(doneCh)
	}()
	select {
	case <-doneCh:
		// Success
	case <-time.After(2 * waitPollingInterval):
		t.Errorf("expected Wait() to have detected the file's existence by now")
	}
}

func TestRealWaiterWaitMissingContent(t *testing.T) {
	tmp, err := ioutil.TempFile("", "real_waiter_test_file")
	if err != nil {
		t.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	rw := realWaiter{}
	doneCh := make(chan struct{})
	go func() {
		err := rw.Wait(tmp.Name(), true)
		if err != nil {
			t.Errorf("error waiting on tmp file %q", tmp.Name())
		}
		close(doneCh)
	}()
	select {
	case <-doneCh:
		t.Errorf("no data was written to tmp file, did not expect Wait() to have detected a non-zero file size and returned")
	case <-time.After(2 * waitPollingInterval):
		// Success
	}
}

func TestRealWaiterWaitWithContent(t *testing.T) {
	tmp, err := ioutil.TempFile("", "real_waiter_test_file")
	if err != nil {
		t.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	rw := realWaiter{}
	doneCh := make(chan struct{})
	go func() {
		err := rw.Wait(tmp.Name(), true)
		if err != nil {
			t.Errorf("error waiting on tmp file %q", tmp.Name())
		}
		close(doneCh)
	}()
	if err := ioutil.WriteFile(tmp.Name(), []byte("😺"), 0700); err != nil {
		t.Errorf("error writing content to temp file: %v", err)
	}
	select {
	case <-doneCh:
		// Success
	case <-time.After(2 * waitPollingInterval):
		t.Errorf("expected Wait() to have detected a non-zero file size by now")
	}
}
