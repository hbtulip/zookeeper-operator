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
	"fmt"
	"strconv"
	"strings"
	"time"

	zkopv1 "hmxq.top/zookeeper-operator/api/v1"
	"github.com/samuel/go-zookeeper/zk"
)
//基于https://github.com/samuel/go-zookeepe ，封装实现MyZookeeperClient 

type ZookeeperClient interface {
	Connect(string) error
	CreateNode(*v1beta1.ZookeeperCluster, string) error
	NodeExists(string) (int32, error)
	UpdateNode(string, string, int32) error
	Close()
}

type MyZookeeperClient  struct {
	conn *zk.Conn
}

//连接zk服务
func (client *MyZookeeperClient ) Connect(zkUri string) (err error) {
	host := []string{zkUri}
	conn, _, err := zk.Connect(host, time.Second*5)
	if err != nil {
		return fmt.Errorf("Failed to connect to zookeeper: %s, Reason: %v", zkUri, err)
	}
	client.conn = conn
	return nil
}

//创建zk持久节点
func (client *MyZookeeperClient ) CreateNode(zoo *zkopv1.ZookeeperCluster, zNodePath string) (err error) {
	paths := strings.Split(zNodePath, "/")
	pathLength := len(paths)
	var parentPath string
	for i := 1; i < pathLength-1; i++ {
		parentPath += "/" + paths[i]
		if _, err := client.conn.Create(parentPath, nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
			return fmt.Errorf("Error creating parent zkNode: %s: %v", parentPath, err)
		}
	}
	data := "CLUSTER_SIZE=" + strconv.Itoa(int(zoo.Spec.Replicas))
	childNode := parentPath + "/" + paths[pathLength-1]
	if _, err := client.conn.Create(childNode, []byte(data), 0, zk.WorldACL(zk.PermAll)); err != nil {
		return fmt.Errorf("Error creating sub zkNode: %s: %v", childNode, err)
	}
	return nil
}

//更新zk节点信息
func (client *MyZookeeperClient ) UpdateNode(path string, data string, version int32) (err error) {
	if _, err := client.conn.Set(path, []byte(data), version); err != nil {
		return fmt.Errorf("Error updating zkNode: %v", err)
	}
	return nil
}

//判断zk节点是否存在
func (client *MyZookeeperClient ) NodeExists(zNodePath string) (version int32, err error) {
	exists, zNodeStat, err := client.conn.Exists(zNodePath)
	if err != nil || !exists {
		return -1, fmt.Errorf("Znode exists check failed for path %s: %v", zNodePath, err)
	}
	return zNodeStat.Version, err
}

//关闭zk连接
func (client *MyZookeeperClient ) Close() {
	client.conn.Close()
}
