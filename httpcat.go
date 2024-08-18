package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const apiHost string = "https://http.cat"

func main() {
	app := &cli.App{
		Name:  "http-cat",
		Usage: "Returns a HTTP 🐱 kitty!",
		Action: func(cCtx *cli.Context) error {
			var http_status_code string = cCtx.Args().Get(0)

			if !isStringNumber(http_status_code) {
				fmt.Println(
					"You must provide a HTTP status code! " +
						"E.g. http-cat 404, http-cat 500. Run 'http-cat --help' for more help.",
				)
				return nil
			}

			slog.Debug("Making request to HTTP Cat API...")

			response, err := http.Get(apiHost + "/" + http_status_code)

			if err != nil {
				log.Fatalln("GET request to HTTP Cat API failed:", err)
			}

			var temp_directory_path string = getTempDir()

			if _, err := os.Stat(temp_directory_path); os.IsNotExist(err) {
				_ = os.Mkdir(temp_directory_path, 0777) // everyone can read and write
			}

			var path_to_image string = filepath.Join(getTempDir(), "kitty.jpeg")

			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatalln("Failed to read image from API:", err)
			}

			err = os.WriteFile(path_to_image, body, 0777)

			if err != nil {
				log.Fatalln("Failed to write the 🐈 kitty in the temp directory:", err)
			}

			output, err := exec.Command("chafa", path_to_image, "--scale", "0.7").Output()

			if err != nil {
				log.Fatalln(
					"Failed to execute chafa to display 🐈 kitty! "+
						"Make sure you have it installed and in path. Error:", err,
				)
			}

			fmt.Print(string(output))

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func isStringNumber(string string) bool {
	if _, err := strconv.Atoi(string); err == nil {
		return true
	}

	return false
}

// Currently only supports Linux, fuck you Windows, no kitties for you for the time being.
func getTempDir() string {
	var temp_directory string = os.Getenv("TMPDIR")

	if temp_directory == "" {
		temp_directory = "/tmp"
	}

	return filepath.Join(temp_directory, "http-cat")
}
