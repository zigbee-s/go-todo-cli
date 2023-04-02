package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	ID          string    `json:"id"`
	Task        string    `json:"task"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"createdAt"`
	CompletedAt time.Time
}

type dbData struct {
	Data []item `json:"data"`
}

type Todos []item

func Add(task string) error {

	todo := item{
		Task:      task,
		Done:      false,
		CreatedAt: time.Now(),
	}

	jsonBytes, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/todo/", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code and handle any errors
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func List() error {
	url := "http://localhost:8080/todo/"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dbDatas dbData
	json.Unmarshal(body, &dbDatas)
	var t Todos = dbDatas.Data
	t.Print()
	return nil
}

func (t *Todos) Print() {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignRight, Text: "CreatedAt"},
			{Align: simpletable.AlignRight, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell

	for idx, item := range *t {
		idx++
		task := blue(item.Task)
		done := blue("no")
		if item.Done {
			task = green(fmt.Sprintf("\u2705 %s", item.Task))
			done = green("yes")
		}
		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: task},
			{Text: done},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: item.CompletedAt.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("You have %d pending todos", t.CountPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) CountPending() int {
	total := 0
	for _, item := range *t {
		if !item.Done {
			total++
		}
	}

	return total
}

func Complete(index int) error {
	url := "http://localhost:8080/todo/"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dbDatas dbData
	json.Unmarshal(body, &dbDatas)

	for i, todo := range dbDatas.Data {
		i++
		if i == index {
			updatedTodo := item{
				Task:        todo.Task,
				Done:        true,
				CreatedAt:   todo.CreatedAt,
				CompletedAt: time.Now(),
			}

			jsonBytes, err := json.Marshal(updatedTodo)
			if err != nil {
				return err
			}

			req, err := http.NewRequest("PUT", "http://localhost:8080/todo/"+todo.ID, bytes.NewBuffer(jsonBytes))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Check the response status code and handle any errors
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}

			return nil
		}
	}
	return nil
}

func Delete(index int) error {
	url := "http://localhost:8080/todo/"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dbDatas dbData
	json.Unmarshal(body, &dbDatas)

	for i, todo := range dbDatas.Data {
		i++
		if i == index {
			req, err := http.NewRequest("DELETE", "http://localhost:8080/todo/"+todo.ID, nil)
			if err != nil {
				return err
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Check the response status code and handle any errors
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}

			return nil
		}
	}
	return nil
}
