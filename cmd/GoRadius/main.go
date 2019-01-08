package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/nudelfabrik/GoRadius"
	"github.com/nudelfabrik/GoRadius/database"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	var userName string
	var vlan int
	var hashMethod string

	file := flag.String("f", "/usr/local/etc/raddb/freeradius.db", "SQLite3 Database file")

	addCommand := flag.NewFlagSet("add", flag.ExitOnError)

	addCommand.StringVar(&userName, "name", "", "Name and Group of new User.")
	addCommand.IntVar(&vlan, "vlan", 0, "VLAN of new User.")
	addCommand.StringVar(&hashMethod, "hash", "nt", "Password Hash. Currently supported: nt")

	flag.Parse()
	db := database.NewDatabase(*file)

	if len(flag.Args()) < 1 {
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Subcommands:\n\n  add - Add New User")
		addCommand.PrintDefaults()
		os.Exit(1)
	}

	switch flag.Args()[0] {
	case "add":
		addCommand.Parse(flag.Args()[1:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if addCommand.Parsed() {
		user := GoRadius.User{}
		if userName == "" {
			userName = requestInput("Username")
		}
		if vlan == 0 {
			var err error
			vlan, err = strconv.Atoi(requestInput("VLAN"))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		user.Name = userName
		user.VLAN = vlan

		fmt.Print("Password:")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Println("Cannot read password:", err)
			os.Exit(1)
		}

		switch hashMethod {
		case "NT", "nt", "Nt", "nT":
			user.PwHash = string(GoRadius.NTHash(string(bytePassword)))
		default:
			fmt.Println("Invalid Password Hash")
			addCommand.PrintDefaults()
			os.Exit(1)

		}
		err = db.AddUser(user)
		if err != nil {
			fmt.Println("cannot add User:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

}

func requestInput(name string) (value string) {
	fmt.Printf("%s: ", name)
	fmt.Scanf("%s\n", &value)
	return
}
