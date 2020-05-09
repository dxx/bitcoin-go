package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func main()  {
	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bucketName := []byte("bucket")

	// 创建
	db.Update(func (tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		// 判断 Bucket 是否存在
		if bucket == nil {
			// 创建 Bucket
			bucket, err = tx.CreateBucket(bucketName)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

	// 存入数据
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		// key 不存在则新增，否则更新
		return bucket.Put([]byte("key1"), []byte("value1"))
	})

	// 获取数据
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		value := bucket.Get([]byte("key1"))
		fmt.Printf("%s", value)
		return nil
	})
}
