package main

import (
    "database/sql"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"encoding/json"
	"net/http"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

type ingredientChoiceResponse struct {
	Ingredients  []ingredient `json:"ingredients"`
	RecipeIds []string `json:"recipeIds"`
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes() 
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/ingredients/related/", a.getIngredientChoiceResponse).Methods("GET")
	a.Router.HandleFunc("/ingredients/top/", a.getTopIngredients).Methods("GET")
	a.Router.HandleFunc("/ingredients/", a.getIngredientsByType).Methods("GET")
	a.Router.HandleFunc("/recipes/", a.getRecipes).Methods("GET")
}

func (a *App) getIngredientsByType(w http.ResponseWriter, r *http.Request) {
	typeFilter := r.URL.Query()["type"]
	if (typeFilter == nil) {
		respondWithError(w, http.StatusInternalServerError, "No ingredient type defined");
		return 
	}
	ingredients, err := getIngredientsByType(a.DB, typeFilter[0])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, ingredients)
}

func (a *App) getTopIngredients(w http.ResponseWriter, r *http.Request) {
	ingredients, err := getTopIngredients(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, ingredients)
}

func (a *App) getRecipes(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["id"]
	recipes, err := getRecipes(a.DB, ids)
	log.Println(ids)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, recipes)
}

func (a *App) getIngredientChoiceResponse(w http.ResponseWriter, r *http.Request) {
	names := r.URL.Query()["name"]

	ingredients, err := getRelatedIngredients(a.DB, names)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	recipeIds, err := getRecipeIds(a.DB, names)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	icr := ingredientChoiceResponse{Ingredients: ingredients, RecipeIds: recipeIds}
	respondWithJSON(w, http.StatusOK, icr)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(a.Router))) 
}
