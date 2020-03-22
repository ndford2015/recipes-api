package main

import (
    "database/sql"
	"fmt"
	"strings"
)

type ingredient struct {
    Name  string `json:"name"`
}
// TODO: Error handling
func getRelatedIngredients(db *sql.DB, names []string) ([]ingredient, error) {
	rows,err := getColByNames(db, "name", names)
	
	if err != nil {
        return nil, err
	}
	return dbRowsToIngredient(rows), nil
}

/* Get rows of specified column given a list of ingredient names - randomly ordered */
func getColByNames(db *sql.DB, colName string, ingrNames []string) (*sql.Rows, error) {
	var query []string
	var whereSelectors []string
	var selectStr string = fmt.Sprintf("select distinct base.%s from recipes.ingredients_v2 base", colName)
	query = append(query, selectStr)
	for index, value := range ingrNames {
		var innerJoin string = fmt.Sprintf(
			`inner join (select name, recipe_id from recipes.ingredients_v2 where name = "%s") t%d on base.recipe_id = t%d.recipe_id`, value, index, index)
		query = append(query, innerJoin)
		whereSelectors = append(whereSelectors, fmt.Sprintf(`"%s"`, value))
	}

	var whereStr = fmt.Sprintf("where base.name not in (%s)", strings.Join(whereSelectors, ","))
	query = append(query, whereStr)
	query = append(query, "order by rand()");
	return db.Query(strings.Join(query, " "))
}

/* Get ingredients specified with a specific type name. E.g meat, fish, etc. */
func getIngredientsByType(db *sql.DB, typeName string) ([]ingredient, error) {
	var query string = fmt.Sprintf(`select distinct name from recipes.ingredients_v2 where ingr_type = "%s" order by rand()`, typeName)
	rows,err := db.Query(query);
	if err != nil {
        return nil, err
	}
	return dbRowsToIngredient(rows), nil
}

func getTopIngredients(db *sql.DB) ([]ingredient, error) {
	var query string = "select name from recipes.ingredients_v2 group by name having count(*) > 25 order by rand()"
	rows,err := db.Query(query);
	if err != nil {
        return nil, err
	}
	return dbRowsToIngredient(rows), nil
}

func dbRowsToIngredient(rows *sql.Rows) ([]ingredient) {
	ingredients := []ingredient{}
	for rows.Next() {
		var i ingredient
		rows.Scan(&i.Name)
		ingredients = append(ingredients, i)
	}
	return ingredients
}