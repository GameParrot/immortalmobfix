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
	var scanner *bufio.Scanner
	var worldPath string
	if len(os.Args) > 1 {
		worldPath = os.Args[1]
	} else {
		fmt.Print("Enter world path: ")
		scanner = bufio.NewScanner(os.Stdin)
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
						mobType, _ := getMapValueAsString(actorData, "identifier")
						if mobType == "minecraft:ender_dragon" {
							continue
						}
						if actorData["Persistent"] == byte(1) {
							actorData["Dead"] = byte(0)
							nbtData, err := nbt.MarshalEncoding(&actorData, nbt.LittleEndian)
							if err == nil {
								name, _ := getMapValueAsString(actorData, "Name")
								printStr := "Fixed dead " + mobType
								if name != "" {
									printStr += " (" + name + ")"
								}
								fmt.Println(printStr)
								db.Put(iter.Key(), nbtData, &opt.WriteOptions{})
							}
						}
					}
				}
			}
		}
	}
	iter.Release()
	db.Close()
	if scanner == nil {
		fmt.Println("Fixed immortal mobs in " + worldPath)
	} else {
		fmt.Println("Fixed immortal mobs in " + worldPath + ". You may now close this window.")
		scanner.Scan() // wait for return press and exit
	}
}

func getMapValueAsString(m map[string]any, key string) (string, bool) {
	if val, ok := m[key]; ok {
		if valStr, ok := val.(string); ok {
			return valStr, true
		}
	}
	return "", false
}
