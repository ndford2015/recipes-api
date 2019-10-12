package main

func main() {
    a := App{} 
    a.Initialize("root", "", "recipes")

    a.Run(":8080")
}