/**
 * Copyright 2020 disp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */
 package utils

 import (
	"testing"
	_ "encoding/json"
	_"fmt"
	_ "os"
	_"time"
)

func TestZKClient(t *testing.T) {
	var err error
	fmt.Println("TestZKClient...")

	zkclient := new(MyZookeeperClient)
	err = zkclient.Connect("127.0.0.0:2181")
	if err != nil {
		//t.Logf(err.Error())
		//t.FailNow()
		t.Fatalf(err.Error())
	}

	err = zkclient.CreateNode(z, "temp1/tmp2/tmp3")
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	err = zkclient.UpdateNode("temp1/tmp2/temp3", "disp", 100)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	_, err = zkclient.NodeExists("temp1")
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	zkclient.Close()

}
