package main

import (
	"context"
	"fmt"
	"homework/internal/service"
	pkg "homework/internal/service/package"
	"homework/internal/storage/db"
	"homework/internal/util"
	"homework/internal/view"
	"log"
)

func main() {
	repository := db.NewSQLRepository(context.Background(), util.NewConfig())

	packageService := pkg.NewPackageService()
	orderService := service.NewOrderService(repository, packageService)
	validationService := service.NewValidationService(repository, packageService)

	commands := view.NewCLI(orderService, validationService)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bye!")
}
