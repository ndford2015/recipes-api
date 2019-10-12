package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type recipe struct {
	Id int `json:"id"`
	Name  string `json:"name"`
	Url string `json:"url"`
	Description string `json:"description"`
}

func getRecipeIds(db *sql.DB, names []string) ([]string, error){
	
	rows,err := getColByNames(db, "recipe_id", names)
	
	if err != nil {
        return nil, err
	}
	
	var recipeIds []string
	for rows.Next() {
		var r string
		rows.Scan(&r)
		recipeIds = append(recipeIds, r)
	}

	return recipeIds, nil
}

func getRecipes(db *sql.DB, ids []string) ([]recipe, error){
	var whereSelectors []string
	for _,value := range ids {
		whereSelectors = append(whereSelectors, fmt.Sprintf(`"%s"`, value))
	}
	query := fmt.Sprintf("select * from recipes.recipes where id in (%s)", strings.Join(whereSelectors, ","))
	rows, err := db.Query(query)
	if err != nil {
        return nil, err
	}
	recipes := []recipe{}
	for rows.Next() {
		var r recipe
		if err := rows.Scan(&r.Id, &r.Name, &r.Url, &r.Description); err != nil {
            return nil, err
		}
		recipes = append(recipes, r)
	}
	return recipes, err
}