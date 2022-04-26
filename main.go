package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/alexflint/go-arg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/jakehl/goid"
)

//go:embed static/*
var static embed.FS

func main() {
	var args struct {
		MaxFileSize   int    `default:"1000000" help:"Maximum file size in kb"`
		FileChunkSize int    `default:"20000" help:"File chunk size in kb"`
		Port          int    `default:"8080" help:"Port to listen on"`
		DataPath      string `default:"./data" help:"Path to data directory"`
	}
	arg.MustParse(&args)
	if args.MaxFileSize < 1 {
		log.Fatal("MaxFileSize can't be less than 1 (kb)")
		return
	}
	if args.FileChunkSize < 1 {
		log.Fatal("FileChunkSize can't be less than 1 (kb)")
		return
	}

	if _, err := os.Stat(args.DataPath); os.IsNotExist(err) {
		err := os.Mkdir(args.DataPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	app := fiber.New(fiber.Config{
		AppName:   "Feather File Upload",
		BodyLimit: args.FileChunkSize * 1000,
	})

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(static),
		PathPrefix: "static",
		Index:      "index.html",
	}))

	app.Static("/", "./data", fiber.Static{
		Download: true,
	})

	app.Get("/chunksize", func(c *fiber.Ctx) error {
		c.SendStatus(200)
		return c.SendString(strconv.Itoa(args.FileChunkSize))
	})

	IDToPathID := make(map[string]string)

	app.Post("/upload", func(c *fiber.Ctx) error {
		if c.Query("id") == "" || c.GetReqHeaders()["File-Name"] == "" {
			return c.SendStatus(400)
		}
		fileName := c.GetReqHeaders()["File-Name"]
		if IDToPathID[c.Query("id")] == "" {
			IDToPathID[c.Query("id")] = goid.NewV4UUID().String()

			folder := path.Join(args.DataPath, IDToPathID[c.Query("id")])

			//mkdir
			err := os.Mkdir(folder, os.ModePerm)
			if err != nil {
				panic(err)
			}

			//make file
			_, err = os.Create(path.Join(folder, fileName))
			if err != nil {
				panic(err)
			}
		}
		folder := path.Join(args.DataPath, IDToPathID[c.Query("id")])
		file, err := os.OpenFile(path.Join(folder, fileName), os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
		if _, err := file.Write(c.Body()); err != nil {
			log.Fatal(err)
		}
		log.Println("Appended to " + "/" + IDToPathID[c.Query("id")] + "/" + fileName)
		return c.SendString("/" + IDToPathID[c.Query("id")] + "/" + fileName)
	})

	app.Listen(":" + strconv.Itoa(args.Port))
}
