package main

import (
	"apify/actor/example/store"
	"log/slog"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
)

var (
	xlog = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

type (
	KVStoreExporter map[string]any
	AutherNText     struct {
		Author string `json:"author"`
		Text   string `json:"text"`
	}
)

func (kv *KVStoreExporter) Export(parsedData chan any) error {
	if len(*kv) == 0 {
		*kv = make(KVStoreExporter)
	}
	data := []AutherNText{}
	for d := range parsedData {
		e := d.(map[string]interface{})
		data = append(data, AutherNText{Author: e["author"].(string), Text: e["text"].(string)})
	}
	(*kv)["data"] = data

	return nil
}

func main() {
	// Get Token and defualt KV store id
	xlog.Info("Example actor written in Go.")

	// Get default KV store
	kv := store.KVStoreDefault()

	// Get input
	input, err := kv.Get("INPUT")
	if err != nil {
		xlog.Error("failed to get input from kv store", "error", err)
		return
	}
	xlog.Info("input from kv store", "input", input)
	url := input["url"].(string)
	if strings.EqualFold(url, "") {
		xlog.Error("no url in input", "url", url)
	}
	// Scrape data
	kvExporter := KVStoreExporter{}
	scrape([]string{url}, &kvExporter)

	xlog.Info("saving scrapped data to kv store", "data", kvExporter)
	if err := kv.Put("data", kvExporter); err != nil {
		xlog.Error("error while add value to store", "data", kvExporter, "error", err)
	}
	xlog.Info("actor is done")
}

func scrape(urls []string, exporter export.Exporter) {
	g := geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: urls,
		ParseFunc: quotesParse,
		Exporters: []export.Exporter{&export.JSON{}, exporter},
	})
	g.Start()
	slog.Info("exportor content", "data", exporter)
}

func quotesParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("div.quote").Each(func(i int, s *goquery.Selection) {
		g.Exports <- map[string]interface{}{
			"text":   s.Find("span.text").Text(),
			"author": s.Find("small.author").Text(),
		}
	})
	if href, ok := r.HTMLDoc.Find("li.next > a").Attr("href"); ok {
		g.Get(r.JoinURL(href), quotesParse)
	}
}
