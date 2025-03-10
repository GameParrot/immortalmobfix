package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/df-mc/goleveldb/leveldb"
	"github.com/df-mc/goleveldb/leveldb/opt"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

func main() {
	var worldPath string
	if len(os.Args) > 1 {
		worldPath = os.Args[1]
	} else {
		fmt.Print("Enter world path: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		worldPath = strings.Trim(strings.TrimSuffix(scanner.Text(), " "), "\"'")
	}
	db, err := leveldb.OpenFile(path.Join(worldPath, "db"), &opt.Options{Compression: opt.FlateCompression})
	if err != nil {
		panic(err)
	}
	iter := db.NewIterator(nil, &opt.ReadOptions{})
	for iter.Next() {
		if strings.HasPrefix(string(iter.Key()), "actorprefix") {
			var actorData = map[string]any{}
			if err := nbt.UnmarshalEncoding(iter.Value(), &actorData, nbt.LittleEndian); err == nil {
				if dead, ok := actorData["Dead"]; ok {
					if dead == byte(1) {
						actorData["Dead"] = byte(0)
						nbtData, err := nbt.MarshalEncoding(&actorData, nbt.LittleEndian)
						if err == nil {
							if id, ok := actorData["identifier"]; ok {
								fmt.Println("Fixed dead", id)
							}
							db.Put(iter.Key(), nbtData, &opt.WriteOptions{})
						}
					}
				}
			}
		}
	}
	iter.Release()
	db.Close()
}
