package main



import (

	"fmt"

	"encoding/json"

	"github.com/bmatsuo/lmdb-go/lmdb"

	"log"

	"net/http"

	"encoding/base64"

	"strconv"

	"strings"

)



func blockchaindump(w http.ResponseWriter, r *http.Request) {

    if r.URL.Path != "/blockchain" {

        http.Error(w, "404 not found.", http.StatusNotFound)

        return

    }

 

    switch r.Method {

    case "GET":     

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

				return

			}



			if err != nil {

				panic(err)

			}

			if(k != nil){



				ind := strings.Index(string(v), "{")

				if(ind == -1){

					fmt.Println("Index error", ind)

				}else{

					fmt.Fprintf(w, "Block: %d\n", blockindex)

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

								if k1 == "payload" {

									decoded, errd := base64.StdEncoding.DecodeString(vv)

									if errd == nil {

										fmt.Fprintf(w, "%s is string %s\n", k1, string(decoded))

									}

								}else {

									fmt.Fprintf(w, "%s is string %s\n", k1, vv)

								}

							case float64:

								fmt.Fprintf(w, "%s is float64 %f\n", k1, vv)

							case []interface{}:

								fmt.Fprintf(w, "%s is an array\n", k1)

								for i, u := range vv {

									fmt.Fprintf(w, "%d %s\n",i, u)

								}

							default:

								fmt.Fprintf(w, "%s is of a type I don't know how to handle\n", k1)

							}

						}

					}	

					fmt.Fprintf(w, "\n")

				}	

			}

			blockindex++

		}

    default:

        fmt.Fprintf(w, "Sorry, only GET methods are supported.")

    }

}



func blockchaindumpblock(w http.ResponseWriter, r *http.Request) {

    if r.URL.Path != "/blockchain/dump/block" {

        http.Error(w, "404 not found.", http.StatusNotFound)

        return

    }

 

    switch r.Method {

    case "POST":

        // Call ParseForm() to parse the raw query and update r.PostForm and r.Form.

        if err := r.ParseForm(); err != nil {

            fmt.Fprintf(w, "ParseForm() err: %v", err)

            return

        }

        searchblock := r.FormValue("block")

		fmt.Println(searchblock)

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

				return

			}



			if err != nil {

				panic(err)

			}

			if(k != nil){



				ind := strings.Index(string(v), "{")

				if(ind == -1){

					fmt.Println("Index error", ind)

				}else{

					substring := string(v)[ind:len(string(v))]



					var f interface{}

					err := json.Unmarshal([]byte(substring), &f)

					if err != nil {

						fmt.Println(err, "err", )

					}else{

						if(searchblock == strconv.Itoa(blockindex)){

							m := f.(map[string]interface{})

							for k1, v1 := range m {

								switch vv := v1.(type) {

								case string:

									if k1 == "payload" {

										decoded, errd := base64.StdEncoding.DecodeString(vv)

										if errd == nil {

											fmt.Fprintf(w, "%s is string %s\n", k1, string(decoded))

										}

									}else {

										fmt.Fprintf(w, "%s is string %s\n", k1, vv)

									}

								case float64:

									fmt.Fprintf(w, "%s is float64 %f\n", k1, vv)

								case []interface{}:

									fmt.Fprintf(w, "%s is an array\n", k1)

									for i, u := range vv {

										fmt.Fprintf(w, "%d %s\n",i, u)

									}

								default:

									fmt.Fprintf(w, "%s is of a type I don't know how to handle\n", k1)

								}

							}

						}

					}	

				}	

			}

			blockindex++

		}

    default:

        fmt.Fprintf(w, "Sorry, only POST methods are supported.")

    }

}



func blockchainsearch(w http.ResponseWriter, r *http.Request) {

    if r.URL.Path != "/blockchain/search" {

        http.Error(w, "404 not found.", http.StatusNotFound)

        return

    }

 

    switch r.Method {

    case "POST":

        // Call ParseForm() to parse the raw query and update r.PostForm and r.Form.

        if err := r.ParseForm(); err != nil {

            fmt.Fprintf(w, "ParseForm() err: %v", err)

            return

        }

        searchkey := r.FormValue("key")

		searchvalue := r.FormValue("value")

		fmt.Println(searchkey, searchvalue)

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

				return

			}



			if err != nil {

				panic(err)

			}

			if(k != nil){



				ind := strings.Index(string(v), "{")

				if(ind == -1){

					fmt.Println("Index error", ind)

				}else{

					substring := string(v)[ind:len(string(v))]



					var f interface{}

					err := json.Unmarshal([]byte(substring), &f)

					if err != nil {

						fmt.Println(err, "err", )

					}else{

						m := f.(map[string]interface{})

						for k1, v1 := range m {

							if( (k1 == searchkey) && (v1 == searchvalue)){

								fmt.Fprintf(w, "Found in block %d\n", blockindex)

							}

						}

					}	

				}	

			}

			blockindex++

		}

    default:

        fmt.Fprintf(w, "Sorry, only POST methods are supported.")

    }

}

 

func main() {

	http.HandleFunc("/blockchain", blockchaindump)

	http.HandleFunc("/blockchain/dump/block", blockchaindumpblock)

	http.HandleFunc("/blockchain/search", blockchainsearch)

 

    fmt.Printf("Starting Blockchain traversal HTTP Server...\n")

    if err := http.ListenAndServe(":8080", nil); err != nil {

        log.Fatal(err)

	}

}
