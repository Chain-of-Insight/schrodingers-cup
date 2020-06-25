package handlers

import (
	"net/http"
	"io/ioutil"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RuleObject struct {
	Type string `json:"kind"`
	Name string `json:"filename"`
	Code string `json:"code"`
}

type RuleList struct {
	Immutable []RuleObject
	Mutable []RuleObject
}

// @description Rules
// @success 200 {object} RuleList "Current list of rules and their contents"
// @router /rules/list [get]
// @produce json
func ListRules(c echo.Context) error {
	
	// Count immutable files
	immutablesCount, err := fileCount("nomsu/rules/immutable")
	if err != nil {
		return err
	}

	// Count mutable files
	mutablesCount, err := fileCount("nomsu/rules/mutable")
	if err != nil {
		return err
	}

	list := new(RuleList)

	// Parse immutables
	for i := 0; i < immutablesCount; i++ {
		// Read rule
		byte_contents, err := ioutil.ReadFile("nomsu/rules/immutable/rule" + strconv.Itoa(i) + ".nom")
		if err != nil {
			return err
		}
		contents := string(byte_contents)
		// Create data structure
		rule := RuleObject{}
		rule.Type = "immutable"
		rule.Name = "rule" + strconv.Itoa(i) + ".nom"
		rule.Code = contents
		// Append to parent struct
		list.Immutable = append(list.Immutable, rule)
	}

	// Parse mutables
	for i := 0; i < mutablesCount; i++ {
		// Read rule
		byte_contents, err := ioutil.ReadFile("nomsu/rules/mutable/rule" + strconv.Itoa(i) + ".nom")
		if err != nil {
			return err
		}
		contents := string(byte_contents)
		// Create data structure
		rule := RuleObject{}
		rule.Type = "mutable"
		rule.Name = "rule" + strconv.Itoa(i) + ".nom"
		rule.Code = contents
		// Append to parent struct
		list.Mutable = append(list.Mutable, rule)
	}

	// Return model / resp.
	r := &RuleList{
		Immutable: list.Immutable,
		Mutable: list.Mutable,
	}

	return c.JSON(http.StatusOK, r)
}

func fileCount(path string) (int, error){
    i := 0
    files, err := ioutil.ReadDir(path)
    if err != nil {
        return 0, err
    }
    for _, file := range files {
        if !file.IsDir() { 
            i++
        }
    }
    return i, nil
}
