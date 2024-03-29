package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/xmp-er/Redis_Go/validatior"
)

func transactional_cmds(s []string, Map map[string]string, Backup_Map map[string]string, conn net.Conn) string {
	//we are taking only MULTI from this

	fmt.Println("OK") //inital ok

	var transactional_store []string //storing everything here then will execute at the end

	for {
		str, err := read_connection_input(conn)

		if err != nil {
			continue
		}

		//checking if the input is correct
		_, err = validatior.Validate_input(str)
		if err != nil {
			fmt.Println(err)
			continue
		}

		st := strings.Split(str, " ")
		if st[0] == "EXEC" {
			break
		} else if st[0] == "DISCARD" {
			return "OK"
		}

		transactional_store = append(transactional_store, str)
		fmt.Println("QUEUED")
	}

	var res string = ""
	var final_res string = ""
	for i := 0; i < len(transactional_store); i++ {
		cmd := transactional_store[i]
		st := strings.Split(cmd, " ")
		switch st[0] {
		case "SET", "GET", "DEL":
			res = crud(st, Map, Backup_Map)
		case "INCR", "INCRBY":
			res = incr_cmds(st, Map, Backup_Map)
		case "MULTI": //recursive MULTI Function maybe
			res = transactional_cmds(st, Map, Backup_Map, conn)
		}
		final_res += fmt.Sprintf("%d)%s\n", i+1, res)
		fmt.Println(i+1, ")", res)
	}
	return final_res
}
