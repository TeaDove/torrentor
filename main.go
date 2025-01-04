package main

import (
	"fmt"
	"github.com/anacrolix/torrent"
)

func main() {
	magnetLink := "magnet:?xt=urn:btih:1AE80FD51FC9591C3369EC1BFA0EDBD3E6CDF019&tr=http%3A%2F%2Fbt.t-ru.org%2Fann%3Fmagnet&dn=%D0%AD%D1%80%D0%B8%D1%85%20%D0%9C%D0%B0%D1%80%D0%B8%D1%8F%20%D0%A0%D0%B5%D0%BC%D0%B0%D1%80%D0%BA%20-%20%D0%A1%D0%BE%D0%B1%D1%80%D0%B0%D0%BD%D0%B8%D0%B5%20%D1%81%D0%BE%D1%87%D0%B8%D0%BD%D0%B5%D0%BD%D0%B8%D0%B9%20%D0%B2%2016%20%D1%82%D0%BE%D0%BC%D0%B0%D1%85%20%5B2011%2C%20EPUB%2C%20RUS%5D" // Замените на вашу магнит-ссылку

	//Создаем новый клиент
	config := torrent.NewDefaultClientConfig()
	config.DataDir = "./data/torrent/"

	client, err := torrent.NewClient(config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Загружаем торрент из магнитной ссылки
	t, err := client.AddMagnet(magnetLink)
	if err != nil {
		panic(err)
	}

	// Ждем, пока торрент будет загружен
	<-t.GotInfo()
	fmt.Printf("Torrent info loaded: %s\n", t.Info().Name)

	t.DownloadAll()
	<-t.Complete().On()

	fmt.Println("Download completed successfully!")
}
