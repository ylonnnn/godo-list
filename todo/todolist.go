package todo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Entry struct {
	Id        int       `json:"id"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"createdAt"`
}

type EntryList struct {
	Entries []Entry `json:"entries"`
}

func (list *EntryList) CreateEntry(desc string) {
	entryLen := len(list.Entries)
	latestId := 0

	// If there are currently no
	if entryLen > 0 {
		latestId = list.Entries[entryLen-1].Id
	}

	list.Entries = append(list.Entries, Entry{
		Id:        latestId + 1,
		Desc:      desc,
		CreatedAt: time.Now(),
	})
}

func (list *EntryList) DeleteEntry(id int) {
	filtered := list.Entries[:0]

	for _, entry := range list.Entries {
		if entry.Id == id {
			continue
		}

		filtered = append(filtered, entry)
	}

	list.Entries = filtered
}

func (list *EntryList) Display() {
	for i := 0; i < len(list.Entries); i++ {
		entry := list.Entries[i]
		fmt.Printf("%d.) [%d] %s - %s\n", i+1, entry.Id, entry.Desc, entry.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

func ManageList() {
	list := LoadList()
	quit := false

	fmt.Println("===== TO-DO LIST =====")

	for {
		if quit {
			break
		}

		options := []string{"Add", "Delete", "List", "Quit"}
		fmt.Println("\nOPTIONS")

		for i, option := range options {
			fmt.Printf("%d.) %s\n", i+1, option)
		}

		fmt.Println()

		var option int
		fmt.Println("Enter an Option: ")

		_, err := fmt.Scan(&option)
		if err != nil {
			log.Fatal("[Error]", err)
		}

		handlers := []func(){
			func() {
				reader := bufio.NewReader(os.Stdin)
				fmt.Println("Enter a task to add: ")

				task, err := reader.ReadString('\n')
				if err != nil {
					log.Fatal("[Error]", err)
				}

				list.CreateEntry(strings.TrimRight(task, "\r\n"))
			},

			func() {
				var id int

				fmt.Println("Enter the ID of the task to delete: ")
				fmt.Scan(&id)

				list.DeleteEntry(id)
			},
			func() { list.Display() },
			func() { quit = true },
		}

		if option < 1 || option > len(options) {
			fmt.Println("Invalid Option:", option)
			continue
		}

		handlers[option-1]()
	}

	SaveEntries(list)
}

func LoadList() *EntryList {
	data, err := os.ReadFile("entries.json")
	if err != nil {
		mode := os.FileMode(0644)

		os.WriteFile("entries.json", []byte("{}"), mode)
		return LoadList()
	}

	list := EntryList{Entries: make([]Entry, 0, 16)}
	json.Unmarshal(data, &list)

	return &list
}

func SaveEntries(entries *EntryList) {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("entries.json", data, 0644)
}
