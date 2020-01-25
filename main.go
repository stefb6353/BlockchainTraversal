package main

import (
	"fmt"
	"encoding/json"
	"github.com/bmatsuo/lmdb-go/lmdb"
	"strings"
)

// OpenLMDBBDatabase opens it
func OpenLMDBBDatabase() error {

	env, err := lmdb.NewEnv()
	if err != nil {
		panic(err)
	}

	defer env.Close()

	err = env.SetMaxDBs(10)
	if err != nil {
		panic(err)
	}

	err = env.SetMapSize(1 << 40)
	if err != nil {
		panic(err)
	}

	err = env.Open("./blockchain.lmdb", lmdb.NoSubdir, 0777)
	defer env.Close()
	if err != nil {
		env.Close()
		panic(err)

	}

	staleReaders, err := env.ReaderCheck()
	if err != nil {
		panic(err)
	}

	if staleReaders > 0 {
		fmt.Printf("cleared %d stale readers", staleReaders)
	}

	txn, err := env.BeginTxn(nil, 0)
	if err != nil {
		panic(err)
	}
	defer txn.Abort()

	_, err = txn.OpenDBI("index_batch", 0)
	if err != nil {
		panic(err)
	}

	_, err = txn.OpenDBI("index_block_num", 0)
	if err != nil {
		panic(err)
	}

	_, err = txn.OpenDBI("index_transaction", 0)
	if err != nil {
		panic(err)
	}

	dbiMain, err := txn.OpenDBI("main", 0)
	if err != nil {
		panic(err)
	}

	cur, err := txn.OpenCursor(dbiMain)
	if err != nil {
		panic(err)
	}

	blockindex := 1;
	for {
		k, v, err := cur.Get(nil, nil, lmdb.Next)

		if lmdb.IsNotFound(err) {
			return nil
		}

		if err != nil {
			panic(err)
		}
		if(k != nil){

			ind := strings.Index(string(v), "{")
			if(ind == -1){
				fmt.Println("Index error", ind)
			}else{
				fmt.Println("Index: ", blockindex)
				substring := string(v)[ind:len(string(v))]

				var f interface{}
				err := json.Unmarshal([]byte(substring), &f)
				if err != nil {
					fmt.Println(err, "err", )
				}else{
					m := f.(map[string]interface{})
					for k1, v1 := range m {
						switch vv := v1.(type) {
						case string:
							fmt.Println(k1, "is string", vv)
						case float64:
							fmt.Println(k1, "is float64", vv)
						case []interface{}:
							fmt.Println(k1, "is an array:")
							for i, u := range vv {
								fmt.Println(i, u)
							}
						default:
							fmt.Println(k1, "is of a type I don't know how to handle")
						}
					}
				}	
			}	
		}
		blockindex++
	}

}

func main() {

	err := OpenLMDBBDatabase()
	if err != nil {
		panic(err)
	}

}